package storage

import (
	"database/sql"
	"github.com/makeroo/taxi_scout/ts_errors"
	"math"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)


type SqlDatastore struct {
	*sql.DB
	Logger *zap.SugaredLogger
	InvitationDuration time.Duration

	preparedStatements map[string]*sql.Stmt
	preparedStatementsLock *sync.RWMutex
}

func NewSqlDatastore(driver string, dataSourceName string, logger *zap.SugaredLogger) (*SqlDatastore, error) {
	db, err := sql.Open(driver, dataSourceName)

	if err != nil {
		return nil, err
	}

	return NewSqlDatastorec(driver, db, logger)
}

func NewSqlDatastorec(driver string, dataSourceConnection *sql.DB, logger *zap.SugaredLogger) (*SqlDatastore, error) {
	if err := dataSourceConnection.Ping(); err != nil {
		return nil, err
	}

	// FIXME: trasferirlo in un file di settings
	return &SqlDatastore{
		dataSourceConnection,
		logger,
		time.Duration(10) * time.Hour,
		map[string]*sql.Stmt{},
		&sync.RWMutex{},
	}, nil
}

func (db *SqlDatastore) Close() {
	for key, stmt := range db.preparedStatements {
		//fmt.Println("Key:", key, "Value:", value)
		err := stmt.Close()

		if err != nil {
			db.Logger.Warnf("error closing statement: query=%v, error=%v", key, err)
		}
	}

	db.Close()
}

func (db *SqlDatastore) rollback(tx *sql.Tx) {
	err := tx.Rollback()

	if err != nil {
		db.Logger.Errorf("rollback failed: error=%v", err)
	}
}

func (db *SqlDatastore) commit(tx *sql.Tx) error {
	err := tx.Commit()

	if err != nil {
		db.Logger.Errorf("commit failed: error=%v", err)

		db.rollback(tx)

		return err
	}

	return nil
}

func (db *SqlDatastore) execStmt (stmtName string, tx *sql.Tx, val ...interface{}) error {
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

func (db *SqlDatastore) CheckPermission (userId int32, groupId int32, permId int32) error {
	return db.checkPermission(userId, groupId, permId, nil)
}

func (db *SqlDatastore) checkPermission (userId int32, groupId int32, permId int32, tx *sql.Tx) error {
	stmt, err := db.stmt("check_permission", tx)

	if err != nil {
		return err
	}

	row := stmt.QueryRow(userId, groupId, permId)

	var count int64

	err = row.Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return ts_errors.Forbidden
	}

	return nil
}

func (db *SqlDatastore) QueryInvitationToken (token string) (*Account, bool, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, false, err
	}

	stmt, err := db.stmt("fetch_invitation", tx)

	if err != nil {
		db.rollback(tx)

		return nil, false, err
	}

	row := stmt.QueryRow(token)

	var invitationEmail string
	var invitationCreatedOn time.Time
	var scoutGroupId int32
	var accountId sql.NullInt64
	var accountName sql.NullString
	var accountAddress sql.NullString

	err = row.Scan(
		&invitationEmail, &invitationCreatedOn,
		&scoutGroupId,
		&accountId, &accountName, &accountAddress,
		)

	if err != nil {
		db.rollback(tx)

		return nil, false, err
	}

	account := new(Account)
	var found bool

	if accountId.Valid {
		account.Id = int32(accountId.Int64)
		account.Name = accountName.String
		account.Email = invitationEmail
		account.Address = accountAddress.String

		found = true

	} else {
		invitationExpires := invitationCreatedOn.Add(db.InvitationDuration)

		if invitationExpires.Before(time.Now()) {
			err = db.execStmt("delete_invitation", tx, token)

			if err != nil {
				return nil, false, err
			}

			err = db.commit(tx)

			if err == nil {
				err = ts_errors.Expired
			}

			return nil, false, err
		}

		stmt, err = db.stmt("create_account_from_invitation", tx)

		if err != nil {
			db.rollback(tx)
			return nil, false, err
		}

		res, err := stmt.Exec(token)

		if err != nil {
			db.rollback(tx)
			return nil, false, err
		}

		id, err := res.LastInsertId()

		if err != nil {
			db.rollback(tx)

			return nil, false, err
		}

		if id > int64(math.MaxInt32) {
			db.rollback(tx)

			return nil, false, IdOverflow
		}

		account.Id = int32(id)
		account.Email = invitationEmail

		found = false
	}

	err = db.checkPermission(account.Id, scoutGroupId, PermissionMember, tx)

	if err == ts_errors.Forbidden {
		err = db.execStmt("grant", tx, PermissionMember, account.Id, scoutGroupId)
	} else if err != nil {
		db.rollback(tx)
	}

	if err != nil {
		return nil, false, err
	}

	err = db.execStmt("delete_invitation", tx, token)

	if err != nil {
		db.Logger.Errorf("delete invitation failed: error=%v", err)
	} else {
		err = db.commit(tx)
	}

	return account, found, err
}

func (db *SqlDatastore) QueryAccounts(group int32, userId int32) ([]*Account, error) {
	err := db.checkPermission(userId, group, PermissionMember, nil)

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
		err := rows.Scan(&account.Id, &account.Name, &account.Email, &account.Address)

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

func (db *SqlDatastore) QueryAccount(id int32) (*Account, error) {
	stmt, err := db.stmt("query_account", nil)

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	account := new(Account)
	err = row.Scan(&account.Id, &account.Name, &account.Email, &account.Address)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *SqlDatastore) AccountUpdate(account *Account) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	stmt, err := db.stmt("update_account", tx)

	if err != nil {
		db.rollback(tx)

		return err
	}

	_, err = stmt.Exec(account.Name, account.Address, account.Id)

	if err != nil {
		db.rollback(tx)

		return err
	}

	return db.commit(tx)
}

/*
func (db *SqlDatastore) InsertAccount(account *AccountWithCredentials) (int32, error) {
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
		return 0, IdOverflow
	}

	return int32(id), err
}
*/
func (db *SqlDatastore) AuthenticateAccount(email string, pwd string) (int32, error) {
	stmt, err := db.stmt("account_credentials", nil)

	if err != nil {
		return 0, err
	}

	row := stmt.QueryRow(email)

	if err != nil {
		return 0, err
	}

	var accId int32
	var hashedPassword string
	err = row.Scan(&accId, &hashedPassword)

	if err != nil {
		return 0, err
	}

	db.Logger.Debugf("hashed pwd: %s %d --%s--", hashedPassword, accId, pwd)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pwd))

	return accId, err
}

func (db *SqlDatastore) UpdateAccountPassword(id int32, oldPwd string, newPwd string) error {
	return nil // TODO
}

func (db *SqlDatastore) AccountGroups(accountId int32) ([]*ScoutGroup, error) {
	stmt, err := db.stmt("account_groups", nil)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(accountId, PermissionMember)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := make([]*ScoutGroup, 0)

	for rows.Next() {
		group := new(ScoutGroup)
		err := rows.Scan(&group.Id, &group.Name)

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

func (db *SqlDatastore) AccountScouts(accountId int32) ([]*Scout, error) {
	stmt, err := db.stmt("account_scouts", nil)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(accountId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	scouts := make([]*Scout, 0)

	for rows.Next() {
		scout := new(Scout)
		err := rows.Scan(&scout.Id, &scout.Name)

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
func (db *SqlDatastore) checkPermissionStmt(stmtName string, tx *sql.Tx, args... interface{}) error {
	stmt, err := db.stmt(stmtName, tx)

	row := stmt.QueryRow(args...)

	var count int64

	err = row.Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return ts_errors.Forbidden
	}

	return nil
}*/

func (db *SqlDatastore) UpdateScout(scout Scout, tutorId int32) (int32, error) {
	tx, err := db.Begin()

	if err != nil {
		return 0, err
	}

	var r int32 = 0

	if scout.Id < 0 {
		// insert a new scout
		err = db.checkPermission(tutorId, scout.GroupId, PermissionMember, tx)

		if err != nil {
			db.rollback(tx)
			return 0, err
		}


		stmt, err := db.stmt("insert_scout", tx)

		if err != nil {
			db.rollback(tx)
			return 0, err
		}

		res, err := stmt.Exec(scout.Name, scout.GroupId)

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

			return 0, IdOverflow
		}

		r = int32(lid)
	} else {
		stmt, err := db.stmt("check_if_tutor", tx)

		row := stmt.QueryRow(scout.Id, tutorId)

		var count int64

		err = row.Scan(&count)

		if err != nil {
			db.rollback(tx)

			return 0, err
		}

		if count == 0 {
			db.rollback(tx)

			return 0, ts_errors.Forbidden
		}

		stmt, err = db.stmt("update_scout", tx)

		_, err = stmt.Exec(scout.Name, scout.Id, scout.GroupId)

		r = scout.Id
	}

	if err != nil {
		db.rollback(tx)

		return 0, err
	}

	return r, db.commit(tx)
}

func (db *SqlDatastore) stmt(query string, tx *sql.Tx) (*sql.Stmt, error) {
	db.preparedStatementsLock.RLock()
	stmt, found := db.preparedStatements[query]
	db.preparedStatementsLock.RUnlock()

	var err error

	if !found {
		db.preparedStatementsLock.Lock()
		stmt, found = db.preparedStatements[query]
		defer db.preparedStatementsLock.Unlock()

		if !found {
			sqlQuery, found := SqlQueries[query]

			if !found {
				return nil, UnknownQuery // TODO: error parameter
			}

			stmt, err = db.Prepare(sqlQuery)

			if err != nil {
				return nil, err
			}

			db.preparedStatements[sqlQuery] = stmt
		}
	}

	if tx != nil {
		stmt = tx.Stmt(stmt)
	}

	return stmt, nil
}
