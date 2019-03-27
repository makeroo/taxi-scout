package rest_backend

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/makeroo/taxi_scout/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/makeroo/taxi_scout/storage"
	"go.uber.org/zap"
)


func serverSetup(t *testing.T) (mockCtrl *gomock.Controller, mockDatastore *mocks.MockDatastore, mockCookieManager *mocks.MockCookieManager, server *RestServer) {
	logger := zap.NewExample().Sugar()

	mockCtrl = gomock.NewController(t)

	mockDatastore = mocks.NewMockDatastore(mockCtrl)

	mockCookieManager = mocks.NewMockCookieManager(mockCtrl)

	server = &RestServer{
		Dao:    mockDatastore,
		Logger: logger,
		Configuration: Configuration{
			SecureCookies: mockCookieManager,
			HttpsCookies: false,
		},
	}

	return
}

func testSuccessfulResponse (t *testing.T, response *httptest.ResponseRecorder, expectedValue interface{}, decoderFunc func (decoder *json.Decoder) (interface{}, error)) {
	if code := response.Code; code != 200 {
		t.Errorf("http query failed: expected 200 received %v", code)
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

func TestAccounts(t *testing.T) {
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
				{1, "name", "email", "addr"},
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

	testSuccessfulResponse(t, recorder, []storage.Account{
		{1, "name", "email", "addr"},
	}, func (decoder *json.Decoder) (interface{}, error) {
		var v []storage.Account
		err := decoder.Decode(&v)

		return v, err
	})
}
