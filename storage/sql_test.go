package storage

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"go.uber.org/zap"
)

func TestSqlDatastore(t *testing.T) {
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
