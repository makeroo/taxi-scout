package storage

import (
	"database/sql"
	"errors"
	"math"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SqlDatastore struct {
	*sql.DB
	Logger *zap.SugaredLogger

	queryAccounts      *sql.Stmt
	queryAccount       *sql.Stmt
	insertAccount      *sql.Stmt
	accountCredentials *sql.Stmt
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

	return &SqlDatastore{dataSourceConnection, logger, nil, nil, nil, nil}, nil
}

func (db *SqlDatastore) Close() {
	if db.queryAccounts != nil {
		db.queryAccounts.Close()
	}
	if db.queryAccount != nil {
		db.queryAccount.Close()
	}
	if db.insertAccount != nil {
		db.insertAccount.Close()
	}
	if db.accountCredentials != nil {
		db.accountCredentials.Close()
	}
	db.Close()
}

func (db *SqlDatastore) QueryAccounts() ([]*Account, error) {
	stmt, err := db.stmtQueryAccounts()

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
	stmt, err := db.stmtQueryAccount()

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	if err != nil {
		return nil, err
	}

	account := new(Account)
	err = row.Scan(&account.Id, &account.Name, &account.Email)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (db *SqlDatastore) InsertAccount(account *AccountWithCredentials) (int32, error) {
	stmt, err := db.stmtInsertAccount()

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
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()

	tx.Commit()

	if id > int64(math.MaxInt32) {
		return 0, errors.New(IdOverflow)
	}

	return int32(id), err
}

func (db *SqlDatastore) AuthenticateAccount(email string, pwd string) (int32, error) {
	stmt, err := db.stmtAccountCredentials()

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

func (db *SqlDatastore) stmtQueryAccounts() (*sql.Stmt, error) {
	if db.queryAccounts == nil {
		stmt, err := db.Prepare("SELECT id, name, email FROM account")

		if err != nil {
			return nil, err
		}

		db.queryAccounts = stmt
	}

	return db.queryAccounts, nil
}

func (db *SqlDatastore) stmtQueryAccount() (*sql.Stmt, error) {
	if db.queryAccount == nil {
		stmt, err := db.Prepare("SELECT id, name, email FROM account WHERE id=?")

		if err != nil {
			return nil, err
		}

		db.queryAccount = stmt
	}

	return db.queryAccount, nil
}

func (db *SqlDatastore) stmtInsertAccount() (*sql.Stmt, error) {
	if db.insertAccount == nil {
		stmt, err := db.Prepare("INSERT INTO account (name, email, pwd) VALUES ( ?, ?, ? )")

		if err != nil {
			return nil, err
		}

		db.insertAccount = stmt
	}

	return db.insertAccount, nil
}

func (db *SqlDatastore) stmtAccountCredentials() (*sql.Stmt, error) {
	if db.accountCredentials == nil {
		stmt, err := db.Prepare("SELECT id, pwd FROM account WHERE email = ?")

		if err != nil {
			return nil, err
		}

		db.accountCredentials = stmt
	}

	return db.accountCredentials, nil
}
