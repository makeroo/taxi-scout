package rest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	tserrors "github.com/makeroo/taxi_scout/errors"
	"github.com/makeroo/taxi_scout/mocks"

	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)

func serverSetup(t *testing.T) (mockCtrl *gomock.Controller, mockDatastore *mocks.MockDatastore, mockCookieManager *mocks.MockCookieManager, server *Server) {
	logger := zap.NewExample().Sugar()

	mockCtrl = gomock.NewController(t)

	mockDatastore = mocks.NewMockDatastore(mockCtrl)

	mockCookieManager = mocks.NewMockCookieManager(mockCtrl)

	server = &Server{
		Dao:    mockDatastore,
		Logger: logger,
		Configuration: Configuration{
			SecureCookies: mockCookieManager,
			HTTPSCookies:  false,
		},
	}

	return
}

func cookieSetup(request *http.Request, recorder *httptest.ResponseRecorder, name string, value string) {
	http.SetCookie(recorder, &http.Cookie{Name: name, Value: value})
	for _, v := range recorder.HeaderMap["Set-Cookie"] {
		request.Header.Add("Cookie", v)
	}
}

func testResponse(t *testing.T, response *httptest.ResponseRecorder, expectedCode int, expectedValue interface{}, decoderFunc func(decoder *json.Decoder) (interface{}, error)) {
	if code := response.Code; code != expectedCode {
		t.Errorf("http query failed: expected=%v received=%v", expectedCode, code)
		return
	}

	cType := response.HeaderMap["Content-Type"]
	if len(cType) != 1 || cType[0] != "application/json" {
		t.Errorf("content-type mismatch: expected=application/json, received=%s", cType)
		return
	}

	// reflect approach does NOT work: interf type is interface{} instead of an actual type
	// eg. []storage.Account
	//	expectedType := reflect.TypeOf(expectedValue)
	//	received := reflect.New(expectedType)
	//	elem := received.Elem()
	//	interf := elem.Interface()

	decoder := json.NewDecoder(response.Body)

	receivedData, err := decoderFunc(decoder)
	if err != nil {
		t.Errorf("json response decode failed: %v", err)
		return
	}

	if !reflect.DeepEqual(receivedData, expectedValue) {
		t.Errorf("mismatch http response: expected=%v received=%v", expectedValue, receivedData)

		return
	}
}

func TestAccountsGet_WithCookie(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(nil).
			Do(func(name string, value string, dst *int32) {
				*dst = 9
			}),
		mockDatastore.EXPECT().QueryAccounts(int32(1), int32(9)).
			Return([]*storage.Account{
				{ID: 1, Name: "name", Email: "email", Address: "addr"},
			}, nil),
	)

	request := httptest.NewRequest("GET", "/accounts?group=1", nil)

	recorder := httptest.NewRecorder()

	cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 200, []storage.Account{
		{ID: 1, Name: "name", Email: "email", Address: "addr"},
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v []storage.Account
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsGet_NoCookie(t *testing.T) {
	mockCtrl, _, _, server := serverSetup(t)
	defer mockCtrl.Finish()

	request := httptest.NewRequest("GET", "/accounts?group=1", nil)

	recorder := httptest.NewRecorder()

	//cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 401, map[string]string{
		"error": "not_authorized",
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]string
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsGet_IllegalCookie(t *testing.T) {
	mockCtrl, _, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(errors.New("decode error")),
	)

	request := httptest.NewRequest("GET", "/accounts?group=1", nil)

	recorder := httptest.NewRecorder()

	cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 400, map[string]string{
		"error": "bad_request",
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]string
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_WithoutCookie_NewAccount(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockDatastore.EXPECT().QueryInvitationToken("xxx", storage.NoRequestingUser).Return(&storage.Account{
			ID: 1, Name: "name", Email: "email", Address: "addr",
		}, false, nil),
		mockCookieManager.EXPECT().Encode("_ts_u", int32(1)).Return("cookie1", nil),
		/*		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
					Return(nil).
					Do(func(name string, value string, dst *int32) {
						*dst = 9
					}),
				mockDatastore.EXPECT().QueryAccounts(int32(1), int32(9)).
					Return([]*storage.Account{
						{1, "name", "email", "addr"},
					}, nil),*/
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	//cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 200, map[string]interface{}{
		"id":            float64(1), // FIXME: this is float32 on 32bit os
		"authenticated": true,
		"new_account":   true,
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]interface{}
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_WithoutCookie_NewGroup(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockDatastore.EXPECT().QueryInvitationToken("xxx", storage.NoRequestingUser).Return(&storage.Account{
			ID: 1, Name: "name", Email: "email", Address: "addr",
		}, true, nil),
		mockCookieManager.EXPECT().Encode("_ts_u", int32(1)).Return("cookie1", nil),
		/*		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
					Return(nil).
					Do(func(name string, value string, dst *int32) {
						*dst = 9
					}),
				mockDatastore.EXPECT().QueryAccounts(int32(1), int32(9)).
					Return([]*storage.Account{
						{1, "name", "email", "addr"},
					}, nil),*/
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	//cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 200, map[string]interface{}{
		"id":            float64(1), // FIXME: this is float32 on 32bit os
		"authenticated": true,
		"new_account":   false,
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]interface{}
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_WithCookie_AccountMismatch(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(nil).
			Do(func(name string, value string, dst *int32) {
				*dst = 9
			}),
		mockDatastore.EXPECT().QueryInvitationToken("xxx", int32(9)).Return(nil, false, tserrors.StolenToken),
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 403, map[string]string{
		"error": "stolen_token",
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]string
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_ValidTokenWithCookie_NewGroup(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(nil).
			Do(func(name string, value string, dst *int32) {
				*dst = 9
			}),
		mockDatastore.EXPECT().QueryInvitationToken("xxx", int32(9)).Return(&storage.Account{
			ID: 1, Name: "name", Email: "email", Address: "addr",
		}, true, nil),
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 200, map[string]interface{}{
		"id":            float64(1), // FIXME: this is float32 on 32bit os
		"authenticated": false,
		"new_account":   false,
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]interface{}
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_InvalidTokenWithoutCookie(t *testing.T) {
	mockCtrl, mockDatastore, _, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockDatastore.EXPECT().QueryInvitationToken("xxx", storage.NoRequestingUser).Return(nil, false, sql.ErrNoRows),
		//mockCookieManager.EXPECT().Encode("_ts_u", int32(1)).Return("cookie1", nil),
		/*		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
					Return(nil).
					Do(func(name string, value string, dst *int32) {
						*dst = 9
					}),
				mockDatastore.EXPECT().QueryAccounts(int32(1), int32(9)).
					Return([]*storage.Account{
						{1, "name", "email", "addr"},
					}, nil),*/
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	//cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 404, map[string]interface{}{
		"error": "not_found",
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]interface{}
		err := decoder.Decode(&v)

		return v, err
	})
}

func TestAccountsPost_InvalidTokenWithCookie(t *testing.T) {
	mockCtrl, mockDatastore, mockCookieManager, server := serverSetup(t)
	defer mockCtrl.Finish()

	gomock.InOrder(
		mockCookieManager.EXPECT().Decode("_ts_u", "expected", gomock.Any()).
			Return(nil).
			Do(func(name string, value string, dst *int32) {
				*dst = 9
			}),
		mockDatastore.EXPECT().QueryInvitationToken("xxx", int32(9)).Return(nil, false, sql.ErrNoRows),
		//mockCookieManager.EXPECT().Encode("_ts_u", int32(1)).Return("cookie1", nil),
	)

	body, _ := json.Marshal(AccountsRequest{
		Token: "xxx",
	})

	request := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))

	recorder := httptest.NewRecorder()

	cookieSetup(request, recorder, "_ts_u", "expected")

	handler := http.HandlerFunc(server.Accounts)
	handler.ServeHTTP(recorder, request)

	testResponse(t, recorder, 200, map[string]interface{}{
		"id":            float64(9),
		"new_account":   false,
		"authenticated": false,
	}, func(decoder *json.Decoder) (interface{}, error) {
		var v map[string]interface{}
		err := decoder.Decode(&v)

		return v, err
	})
}
