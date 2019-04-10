package storage

import (
	"database/sql"
	"database/sql/driver"
	"github.com/google/uuid"
	"github.com/makeroo/taxi_scout/ts_errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"go.uber.org/zap"
)

func TestSqlDatastore_QueryAccounts(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	//mock.ExpectBegin()
	mock.ExpectPrepare("SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?")
	mock.ExpectQuery("SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?").
		WithArgs(3, 7, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(1))
	mock.ExpectPrepare("SELECT a.id, a.name, a.email, a.address FROM account a JOIN account_grant g ON g.account_id = a.id WHERE g.group_id = ?")
	mock.ExpectQuery("SELECT a.id, a.name, a.email, a.address FROM account a JOIN account_grant g ON g.account_id = a.id WHERE g.group_id = ?").
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "address"}).
			AddRow(45, "pippo", "pippo@dom", "addr"))
	//mock.ExpectCommit()

	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	accounts, err := dao.QueryAccounts(7, 3)
	if err != nil {
		t.Fatalf("QueryAccounts failed: error=%s", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expecting 1 accounts, found: %d", len(accounts))
	}

	expectedAccount := Account{45, "pippo", "pippo@dom", "addr"}
	if *accounts[0] != expectedAccount {
		t.Errorf("account mismatch: expected %v, found %v", expectedAccount, accounts[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

type anyUUID struct{}

func (_ anyUUID) Match(v driver.Value) bool {
	if s, ok := v.(string); ok {
		_, err := uuid.Parse(s)

		return err == nil
	}

	return false
}

type now struct{}

func (_ now) Match(v driver.Value) bool {
	if t, ok := v.(time.Time); ok {
		d := time.Now().Sub(t)
		dn := d.Nanoseconds()

		return dn < 10000000
	}

	return false
}

func TestSqlDatastore_CreateInvitationForExistingMemberOk(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`)
	mock.ExpectPrepare(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`)
	mock.ExpectExec(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`).
		WithArgs(anyUUID{}, now{}, "mail@h").
		WillReturnResult(sqlmock.NewResult(109, 1))

	mock.ExpectCommit()

	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	invitation, err := dao.CreateInvitationForExistingMember("mail@h")
	if err != nil {
		t.Fatalf("CreateInvitationForExistingMember failed: error=%s", err)
	}

	if invitation.ScoutGroup != nil {
		t.Errorf("invitation mismatch: expected group %v, found %v", nil, invitation.ScoutGroup)
	}
	if invitation.Email != "mail@h" {
		t.Errorf("invitation mismatch: exected email %v, found %v", "mail@h", invitation.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSqlDatastore_CreateInvitationForExistingMemberFail(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`)
	mock.ExpectPrepare(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`)
	mock.ExpectExec(`
INSERT INTO invitation (token, email, created_on, group_id)
     SELECT ?, a.email, ?, g.group_id
       FROM account a
       JOIN account_grant g ON g.account_id=a.id
      WHERE a.email = ?
      LIMIT 1`).
		WithArgs(anyUUID{}, now{}, "mail@h").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectRollback()

	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	_, err = dao.CreateInvitationForExistingMember("mail@h")
	if err != ts_errors.Forbidden {
		t.Fatalf("CreateInvitationForExistingMember failed: expected forbidden, received=%s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSqlDatastore_QueryInvitationToken_NoToken(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`)
	mock.ExpectPrepare(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`)
	mock.ExpectQuery(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`).
		WithArgs("xxx").
		WillReturnRows(sqlmock.NewRows([]string{"email", "created_on", "group_id", "id", "name", "address"}))
	mock.ExpectRollback()

	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	_, _, err = dao.QueryInvitationToken("xxx", NoRequestingUser)

	if err != sql.ErrNoRows {
		t.Errorf("QueryInvitationToken did not failed with no rows: error=%v", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSqlDatastore_QueryInvitationToken_ValidToken(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`)
	mock.ExpectPrepare(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`)
	mock.ExpectQuery(`
   SELECT i.email, i.created_on, i.group_id,
  	      a.id, a.name, a.address
     FROM invitation i
LEFT JOIN account a ON i.email = a.email
    WHERE i.token = ?`).
		WithArgs("xxx").
		WillReturnRows(sqlmock.NewRows([]string{"email", "created_on", "group_id", "id", "name", "address"}).
			AddRow("mail@h", time.Now(), 32, sql.NullInt64{}, sql.NullString{}, sql.NullString{}))
	mock.ExpectPrepare("INSERT INTO account (email) SELECT email FROM invitation WHERE token = ?")
	mock.ExpectPrepare("INSERT INTO account (email) SELECT email FROM invitation WHERE token = ?")
	mock.ExpectExec("INSERT INTO account (email) SELECT email FROM invitation WHERE token = ?").
		WithArgs("xxx").
		WillReturnResult(sqlmock.NewResult(23, 1))
	mock.ExpectPrepare("SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?")
	mock.ExpectPrepare("SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?")
	mock.ExpectQuery("SELECT count(*) FROM account_grant WHERE account_id = ? AND group_id = ? AND permission_id = ?").
		WithArgs(23, 32, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(0))
	mock.ExpectPrepare("INSERT INTO account_grant (permission_id, account_id, group_id) VALUES (?, ?, ?)")
	mock.ExpectPrepare("INSERT INTO account_grant (permission_id, account_id, group_id) VALUES (?, ?, ?)")
	mock.ExpectExec("INSERT INTO account_grant (permission_id, account_id, group_id) VALUES (?, ?, ?)").
		WithArgs(1, 23, 32).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare("DELETE FROM invitation WHERE token = ?")
	mock.ExpectPrepare("DELETE FROM invitation WHERE token = ?")
	mock.ExpectExec("DELETE FROM invitation WHERE token = ?").
		WithArgs("xxx").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()


	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	account, found, err := dao.QueryInvitationToken("xxx", NoRequestingUser)

	if err != nil {
		t.Errorf("QueryInvitationToken failed: error=%v", err)
		return
	}

	if found {
		t.Errorf("QueryInvitationToken did not create account")
		return
	}

	expectedAccount := Account{23, "", "mail@h", ""}
	if *account != expectedAccount {
		t.Errorf("expected account: %v, found %v", expectedAccount, account)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSqlDatastore_QueryAccount(t *testing.T) {
	// TODO
}

func TestSqlDatastore_AuthenticateAccount(t *testing.T) {
	// TODO
}

func TestSqlDatastore_UpdateAccountPassword(t *testing.T) {
	// TODO
}

func TestSqlDatastore_AccountGroups(t *testing.T) {
	// TODO
}

func TestSqlDatastore_AccountScouts(t *testing.T) {
	// TODO
}

func TestSqlDatastore_AccountUpdate(t *testing.T) {
	// TODO
}

func TestSqlDatastore_InsertOrUpdateScout(t *testing.T) {
	// TODO
}

func TestSqlDatastore_RemoveScout(t *testing.T) {
	// TODO
}
