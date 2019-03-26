package rest_backend

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/makeroo/taxi_scout/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)


func TestAccounts(t *testing.T) {
	logger := zap.NewExample().Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)

	mockCookieManager := mocks.NewMockCookieManager(mockCtrl)

	server := &RestServer{
		Dao:    mockDatastore,
		Logger: logger,
		Configuration: Configuration{
			SecureCookies: mockCookieManager,
			HttpsCookies: false,
		},
	}

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(nil).
			Do(func(name string, value string, dst *int32) {
				*dst = 9
				}),
		mockDatastore.EXPECT().QueryAccounts(int32(1), int32(9)).
			Return([]*storage.Account{
				&storage.Account{1, "name", "email", "addr"},
			}, nil),
		)

	request := httptest.NewRequest("GET", "/accounts?group=1", nil)

	recorder := httptest.NewRecorder()

	http.SetCookie(recorder, &http.Cookie{Name: "_ts_u", Value: "expected"})
	for _, v := range recorder.HeaderMap["Set-Cookie"] {
		request.Header.Add("Cookie", v)
	}

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("expected 200, found: %v", recorder.Code)
	}

	var resData []storage.Account

	decoder := json.NewDecoder(recorder.Body)
	if err := decoder.Decode(&resData); err != nil {
		t.Errorf("json response decode failed: %v", err)
	}

	if len(resData) != 1 {
		t.Errorf("expecting 1 accounts, found: %d", len(resData))
	}

	acc := resData[0]
	expectedAcc := storage.Account{1, "name", "email", "addr"}

	if acc != expectedAcc {
		t.Errorf("wrong account: expecting %v, found %v", acc, expectedAcc)
	}
}
