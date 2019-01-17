package storage

import (
	"database/sql"
	"errors"
	"math"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)


type SqlDatastore struct {
	*sql.DB
	Logger *zap.SugaredLogger
	InvitationDuration time.Duration

	preparedStatements map[string]*sql.Stmt
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
	return &SqlDatastore{dataSourceConnection, logger, time.Duration(10) * time.Hour, map[string]*sql.Stmt{}}, nil
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

func (db *SqlDatastore) QueryInvitationToken (token string) (*Invitation, *Account, error) {
	stmt, err := db.stmt("fetch_invitation")

	if err != nil {
		return nil, nil, err
	}

	row := stmt.QueryRow(token)

	var invitationEmail string
	var invitationCreatedOn time.Time
	var scoutGroupId int32
	var scoutGroupName string
	var accountId sql.NullInt64
	var accountName sql.NullString
	var accountAddress sql.NullString
	var accountVerified sql.NullBool

	err = row.Scan(
		&invitationEmail, &invitationCreatedOn,
		&scoutGroupId, &scoutGroupName,
		&accountId, &accountName, &accountAddress, &accountVerified,
		)

	if err != nil {
		return nil, nil, err
	}

	if accountId.Valid {
		account := NewAccount()
		account.Id = int32(accountId.Int64)
		account.Name = accountName.String
		account.Email = invitationEmail
		account.Address = accountAddress.String
		account.VerifiedEmail = accountVerified.Bool

		return nil, account, nil

	} else {
		invitation := NewInvitation()

		invitation.Token = token
		invitation.Email = invitationEmail
		invitation.Expires = invitationCreatedOn.Add(db.InvitationDuration)

		scoutGroup := NewScoutGroup()

		scoutGroup.Id = scoutGroupId
		scoutGroup.Name = scoutGroupName

		invitation.ScoutGroup = scoutGroup

		return invitation, nil, nil
	}
}

func (db *SqlDatastore) QueryAccounts() ([]*Account, error) {
	stmt, err := db.stmt("query_accounts")

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := make([]*Account, 0)

	for rows.Next() {
		account := new(Account)
		err := rows.Scan(&account.Id, &account.Name, &account.Email)

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
	stmt, err := db.stmt("query_account")

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	account := new(Account)
	err = row.Scan(&account.Id, &account.Name, &account.Email)

	if err != nil {
		return nil, err
	}

	return account, nil
}

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
	res, err := stmt.Exec(account.Name, account.Email, hashedPassword)

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
		return 0, errors.New(IdOverflow)
	}

	return int32(id), err
}

func (db *SqlDatastore) AuthenticateAccount(email string, pwd string) (int32, error) {
	stmt, err := db.stmt("account_credentials")

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

func (db *SqlDatastore) stmt(query string) (*sql.Stmt, error) {
	stmt, found := db.preparedStatements[query]
	var err error

	if !found {
		sqlQuery, found := SqlQueries[query]

		if !found {
			return nil, errors.New(UnknownQuery) // TODO: error parameter
		}

		stmt, err = db.Prepare(sqlQuery)

		if err != nil {
			return nil, err
		}

		db.preparedStatements[sqlQuery] = stmt
	}

	return stmt, nil
}
