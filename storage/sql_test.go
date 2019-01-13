package storage

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"go.uber.org/zap"
)

func TestSqlDatastore(t *testing.T) {
	logger := zap.NewExample().Sugar()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock failed: error=%s", err)
	}
	defer db.Close()

	//mock.ExpectBegin()
	mock.ExpectPrepare("SELECT id, name, email FROM account")
	mock.ExpectQuery("SELECT id, name, email FROM account").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(45, "pippo", "pippo@dom"))
	//mock.ExpectCommit()

	dao, err := NewSqlDatastorec("mysql", db, logger)
	if err != nil {
		t.Fatalf("failed to create SqlDatastore: error=%s", err)
	}

	accounts, err := dao.QueryAccounts()
	if err != nil {
		t.Fatalf("QueryAccounts failed: error=%s", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expecting 1 accounts, found: %d", len(accounts))
	}

	expectedAccount := NewAccount()
	expectedAccount.Id = 45
	expectedAccount.Name = "pippo"
	expectedAccount.Email = "pippo@dom"
	if *accounts[0] != *expectedAccount {
		t.Errorf("account mismatch: expected %v, found %v", expectedAccount, accounts[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
