package storage

import (
	"database/sql"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	tserrors "github.com/makeroo/taxi_scout/errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SQLDatastore is a Datastore implementation based on SQL backend. Currently only MySQL has been tested.
type SQLDatastore struct {
	*sql.DB
	Logger             *zap.SugaredLogger
	InvitationDuration time.Duration

	preparedStatements     map[string]*sql.Stmt
	preparedStatementsLock *sync.RWMutex
}

// NewSQLDatastore builds a SQLDatastore instance connected do specified datasource.
// Currently driver is ignored and only MySQL is supported.
func NewSQLDatastore(driver string, dataSourceName string, logger *zap.SugaredLogger) (*SQLDatastore, error) {
	db, err := sql.Open(driver, dataSourceName)

	if err != nil {
		return nil, err
	}

	return NewSQLDatastoreWithConnection(driver, db, logger)
}

func NewSQLDatastoreWithConnection(driver string, dataSourceConnection *sql.DB, logger *zap.SugaredLogger) (*SQLDatastore, error) {
	if err := dataSourceConnection.Ping(); err != nil {
		return nil, err
	}

	// TODO: use a configuration struct
	return &SQLDatastore{
		dataSourceConnection,
		logger,
		time.Duration(10) * time.Hour,
		map[string]*sql.Stmt{},
		&sync.RWMutex{},
	}, nil
}

func (db *SQLDatastore) Close() {
	for key, stmt := range db.preparedStatements {
		//fmt.Println("Key:", key, "Value:", value)
		err := stmt.Close()

		if err != nil {
			db.Logger.Warnf("error closing statement: query=%v, error=%v", key, err)
		}
	}

	db.Close()
}

func (db *SQLDatastore) rollback(tx *sql.Tx) {
	err := tx.Rollback()

	if err != nil {
		db.Logger.Errorf("rollback failed: error=%v", err)
	}
}

func (db *SQLDatastore) commit(tx *sql.Tx) error {
	err := tx.Commit()

	if err != nil {
		db.Logger.Errorf("commit failed: error=%v", err)

		db.rollback(tx)

		return err
	}

	return nil
}

func (db *SQLDatastore) execStmtAndRollbackOnFail(stmtName string, tx *sql.Tx, val ...interface{}) error {
	stmt, err := db.stmt(stmtName, tx)

	if err != nil {
		db.rollback(tx)
		return err
	}

	_, err = stmt.Exec(val...)

	if err != nil {
		db.rollback(tx)

		return err
	}

	return nil
}

func (db *SQLDatastore) CheckPermission(userID int32, groupID int32, permID int32) error {
	return db.checkPermission(userID, groupID, permID, nil)
}

func (db *SQLDatastore) checkPermission(userID int32, groupID int32, permID int32, tx *sql.Tx) error {
	stmt, err := db.stmt("check_permission", tx)

	if err != nil {
		return err
	}

	row := stmt.QueryRow(userID, groupID, permID)

	var count int64

	err = row.Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return tserrors.Forbidden
	}

	return nil
}

func (db *SQLDatastore) deleteInvitationAndReturnError(tx *sql.Tx, token string, errorToBeReturned error) (err error) {
	err = db.execStmtAndRollbackOnFail("delete_invitation", tx, token)

	if err != nil {
		return
	}

	err = db.commit(tx)

	if err == nil {
		err = errorToBeReturned
	}

	return
}

func (db *SQLDatastore) QueryInvitationToken(token string, requestingUser int32) (invitedAccount *Account, accountAlreadyExisted bool, scoutGroupID int32, joinedGroup bool, err error) {
	var tx *sql.Tx

	tx, err = db.Begin()

	if err != nil {
		return
	}

	var stmt *sql.Stmt

	stmt, err = db.stmt("fetch_invitation", tx)

	if err != nil {
		db.rollback(tx)

		return
	}

	row := stmt.QueryRow(token)

	var invitationEmail string
	var invitationCreatedOn time.Time
	var accountID sql.NullInt64
	var accountName sql.NullString
	var accountAddress sql.NullString

	err = row.Scan(
		&invitationEmail, &invitationCreatedOn,
		&scoutGroupID,
		&accountID, &accountName, &accountAddress,
	)

	if err != nil {
		db.rollback(tx)

		return
	}

	invitedAccount = new(Account)

	if accountID.Valid {
		invitedAccount.ID = int32(accountID.Int64)
		invitedAccount.Name = accountName.String
		invitedAccount.Email = invitationEmail
		invitedAccount.Address = accountAddress.String

		accountAlreadyExisted = true

		if requestingUser != NoRequestingUser && invitedAccount.ID != requestingUser {
			err = db.deleteInvitationAndReturnError(tx, token, tserrors.StolenToken)
			return
		}

	} else {
		if requestingUser != NoRequestingUser {

			// token email differs from requesting user's email
			// otherwise accountID would be valid, being matched by fetch_invitation query
			// so token has been "stolen", that is used by a user that is not the one
			// the invitation was sent

			err = db.deleteInvitationAndReturnError(tx, token, tserrors.StolenToken)
			return
		}

		invitationExpires := invitationCreatedOn.Add(db.InvitationDuration)

		if invitationExpires.Before(time.Now()) {
			err = db.deleteInvitationAndReturnError(tx, token, tserrors.Expired)
			return
		}

		stmt, err = db.stmt("create_account_from_invitation", tx)

		if err != nil {
			db.rollback(tx)
			return
		}

		var res sql.Result
		res, err = stmt.Exec(token)

		if err != nil {
			db.rollback(tx)
			return
		}

		var id int64
		id, err = res.LastInsertId()

		if err != nil {
			db.rollback(tx)

			return
		}

		if id > int64(math.MaxInt32) {
			db.rollback(tx)

			err = ErrIDOverflow
			return
		}

		invitedAccount.ID = int32(id)
		invitedAccount.Email = invitationEmail

		accountAlreadyExisted = false
	}

	err = db.checkPermission(invitedAccount.ID, scoutGroupID, ScoutGroupMember, tx)

	if err == tserrors.Forbidden {
		err = db.execStmtAndRollbackOnFail("grant", tx, ScoutGroupMember, invitedAccount.ID, scoutGroupID)
		joinedGroup = true
	} else if err != nil {
		db.rollback(tx)
	}

	if err != nil {
		return
	}

	err = db.execStmtAndRollbackOnFail("delete_invitation", tx, token)

	if err != nil {
		db.Logger.Errorf("delete invitation failed: error=%v", err)
	} else {
		err = db.commit(tx)
	}

	return
}

func (db *SQLDatastore) CreateInvitationForExistingMember(email string) (*Invitation, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	stmt, err := db.stmt("create_invitation_for_existing_member", tx)

	if err != nil {
		db.rollback(tx)

		return nil, err
	}

	tokenUUID, err := uuid.NewRandom()

	if err != nil {
		db.rollback(tx)

		return nil, err
	}

	token := tokenUUID.String()

	createdOn := time.Now()

	res, err := stmt.Exec(token, createdOn, email)

	if err != nil {
		db.rollback(tx)

		return nil, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		db.rollback(tx)

		return nil, err
	}

	if rowsAffected == 0 {
		db.rollback(tx)

		return nil, tserrors.Forbidden
	}

	return &Invitation{
		Token:   token,
		Email:   email,
		Expires: createdOn.Add(db.InvitationDuration),
	}, db.commit(tx)
}

func (db *SQLDatastore) QueryAccounts(group int32, userID int32) ([]*Account, error) {
	err := db.checkPermission(userID, group, ScoutGroupMember, nil)

	if err != nil {
		return nil, err
	}

	stmt, err := db.stmt("query_accounts", nil)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(group)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := make([]*Account, 0)

	for rows.Next() {
		account := new(Account)
		err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.Address)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (db *SQLDatastore) QueryAccount(id int32) (*Account, error) {
	stmt, err := db.stmt("query_account", nil)

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	account := new(Account)
	err = row.Scan(&account.ID, &account.Name, &account.Email, &account.Address)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *SQLDatastore) AccountUpdate(account *Account) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	stmt, err := db.stmt("update_account", tx)

	if err != nil {
		db.rollback(tx)

		return err
	}

	_, err = stmt.Exec(account.Name, account.Address, account.ID)

	if err != nil {
		db.rollback(tx)

		return err
	}

	return db.commit(tx)
}

/*
func (db *SQLDatastore) InsertAccount(account *AccountWithCredentials) (int32, error) {
	stmt, err := db.stmt("insert_account")

	if err != nil {
		return 0, err
	}

	// start password
	password := []byte(account.Pwd)
	encryptedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	hashedPassword := string(encryptedPassword)

	if err != nil {
		return 0, err
	}
	// end password

	tx, err := db.Begin()

	if err != nil {
		return 0, err
	}

	stmt = tx.Stmt(stmt)
	res, err := stmt.Exec(account.Name, account.Email, hashedPassword, account.Address)

	if err != nil {
		rbErr := tx.Rollback()

		if rbErr != nil {
			db.Logger.Errorf("rollback failed: error=%v, while-processing-error=%v", rbErr, err)
		}

		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}

	if id > int64(math.MaxInt32) {
		return 0, ErrIDOverflow
	}

	return int32(id), err
}
*/
func (db *SQLDatastore) AuthenticateAccount(email string, pwd string) (int32, error) {
	stmt, err := db.stmt("account_credentials", nil)

	if err != nil {
		return 0, err
	}

	row := stmt.QueryRow(email)

	if err != nil {
		return 0, err
	}

	var accID int32
	var hashedPassword string
	err = row.Scan(&accID, &hashedPassword)

	if err != nil {
		return 0, err
	}

	db.Logger.Debugf("hashed pwd: %s %d --%s--", hashedPassword, accID, pwd)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pwd))

	return accID, err
}

func (db *SQLDatastore) UpdateAccountPassword(id int32, oldPwd string, newPwd string) error {
	return nil // TODO
}

func (db *SQLDatastore) AccountGroups(accountID int32) ([]*ScoutGroup, error) {
	stmt, err := db.stmt("account_groups", nil)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(accountID, ScoutGroupMember)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := make([]*ScoutGroup, 0)

	for rows.Next() {
		group := new(ScoutGroup)
		err := rows.Scan(&group.ID, &group.Name)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (db *SQLDatastore) AccountScouts(accountID int32) ([]*Scout, error) {
	stmt, err := db.stmt("account_scouts", nil)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	scouts := make([]*Scout, 0)

	for rows.Next() {
		scout := new(Scout)
		err := rows.Scan(&scout.ID, &scout.Name, &scout.GroupID)

		if err != nil {
			return nil, err
		}

		scouts = append(scouts, scout)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return scouts, nil
}

/*
func (db *SQLDatastore) checkPermissionStmt(stmtName string, tx *sql.Tx, args... interface{}) error {
	stmt, err := db.stmt(stmtName, tx)

	row := stmt.QueryRow(args...)

	var count int64

	err = row.Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return tserrors.Forbidden
	}

	return nil
}*/

func (db *SQLDatastore) InsertOrUpdateScout(scout Scout, tutorID int32) (int32, error) {
	tx, err := db.Begin()

	if err != nil {
		return 0, err
	}

	var r int32

	// we check for group both in insert and update
	// because we support group change in update
	err = db.checkPermission(tutorID, scout.GroupID, ScoutGroupMember, tx)

	if err != nil {
		db.rollback(tx)
		return 0, err
	}

	if scout.ID < 0 {
		// insert a new scout

		stmt, err := db.stmt("insert_scout", tx)

		if err != nil {
			db.rollback(tx)
			return 0, err
		}

		res, err := stmt.Exec(scout.Name, scout.GroupID)

		if err != nil {
			db.rollback(tx)
			return 0, err
		}

		lid, err := res.LastInsertId()

		if err != nil {
			db.rollback(tx)
			return 0, err
		}

		if lid > int64(math.MaxInt32) {
			db.rollback(tx)

			return 0, ErrIDOverflow
		}

		r = int32(lid)

		stmt, err = db.stmt("add_scout", tx)

		if err != nil {
			db.rollback(tx)

			return 0, err
		}

		_, err = stmt.Exec(tutorID, r)

		if err != nil {
			db.rollback(tx)

			return 0, err
		}
	} else {
		stmt, err := db.stmt("check_if_tutor", tx)

		row := stmt.QueryRow(scout.ID, tutorID)

		var count int64

		err = row.Scan(&count)

		if err != nil {
			db.rollback(tx)

			return 0, err
		}

		if count == 0 {
			db.rollback(tx)

			return 0, tserrors.Forbidden
		}

		stmt, err = db.stmt("update_scout", tx)

		_, err = stmt.Exec(scout.Name, scout.GroupID, scout.ID)

		if err != nil {
			db.rollback(tx)

			return 0, err
		}

		r = scout.ID
	}

	return r, db.commit(tx)
}

func (db *SQLDatastore) RemoveScout(scoutID int32, tutorID int32) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	stmt, err := db.stmt("remove_scout_tutor", tx)

	if err != nil {
		db.rollback(tx)
		return err
	}

	res, err := stmt.Exec(scoutID, tutorID)

	if err != nil {
		db.rollback(tx)
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		db.rollback(tx)
		return err
	}

	if count == 0 {
		return nil
	}

	stmt, err = db.stmt("count_tutors", tx)

	if err != nil {
		db.rollback(tx)
		return err
	}

	row := stmt.QueryRow(scoutID)

	err = row.Scan(&count)

	if err != nil {
		db.rollback(tx)
		return err
	}

	if count == 0 {
		stmt, err = db.stmt("remove_scout", tx)

		if err != nil {
			db.rollback(tx)
			return err
		}

		_, err = stmt.Exec(scoutID)

		if err != nil {
			db.rollback(tx)
			return err
		}
	}

	return db.commit(tx)
}

func (db *SQLDatastore) stmt(query string, tx *sql.Tx) (*sql.Stmt, error) {
	db.preparedStatementsLock.RLock()
	stmt, found := db.preparedStatements[query]
	db.preparedStatementsLock.RUnlock()

	var err error

	if !found {
		db.preparedStatementsLock.Lock()
		stmt, found = db.preparedStatements[query]
		defer db.preparedStatementsLock.Unlock()

		if !found {
			sqlQuery, found := SQLQueries[query]

			if !found {
				return nil, ErrUnknownQuery // TODO: error parameter
			}

			stmt, err = db.Prepare(sqlQuery)

			if err != nil {
				return nil, err
			}

			db.preparedStatements[query] = stmt
		}
	}

	if tx != nil {
		stmt = tx.Stmt(stmt)
	}

	return stmt, nil
}
