package rest_backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)

type MockDatastore struct {
}

func (db *MockDatastore) QueryAccounts() ([]*storage.Account, error) {
	return []*storage.Account{
		&storage.Account{1, "name", "email"},
	}, nil
}

func (db *MockDatastore) QueryAccount(id int32) (*storage.Account, error) {
	return nil, nil
}

func (db *MockDatastore) InsertAccount(*storage.AccountWithCredentials) (int32, error) {
	return 0, nil
}

func (db *MockDatastore) AuthenticateAccount(email string, pwd string) (int32, error) {
	return 0, nil
}

func (db *MockDatastore) UpdateAccountPassword(id int32, oldPwd string, newPwd string) error {
	return nil
}

func TestAccounts(t *testing.T) {
	logger := zap.NewExample().Sugar()

	mockDatastore := &MockDatastore{}

	server := &RestServer{
		Dao:    mockDatastore,
		Logger: logger,
	}

	request := httptest.NewRequest("GET", "/accounts", nil)

	response := httptest.NewRecorder()

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(response, request)

	var resData []storage.Account

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&resData); err != nil {
		t.Errorf("json response decode failed: %v", err)
	}

	if len(resData) != 1 {
		t.Errorf("expecting 1 accounts, found: %d", len(resData))
	}

	acc := resData[0]
	expectedAcc := storage.Account{1, "name", "email"}

	if acc != expectedAcc {
		t.Errorf("wrong account: expecting %v, found %v", acc, expectedAcc)
	}
}
