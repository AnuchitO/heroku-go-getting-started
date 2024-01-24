package authen_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/authen"
	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"gitdev.devops.krungthai.com/aster/ariskill/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() (*app.Router, func()) {
	zapLogger, graceful := logger.NewZap()
	gin.SetMode(gin.TestMode)
	r := app.NewRouter(zapLogger)
	return r, graceful
}

func anyToJsonMarshalToBytesBuffer(body any) *bytes.Buffer {
	data, _ := json.Marshal(body)
	return bytes.NewBuffer(data)
}

func readerCloserToString(rc io.ReadCloser) string {
	data, _ := io.ReadAll(rc)
	return string(data)
}

func TestAuthenHandlerExchangeForTokens(t *testing.T) {
	type testCase struct {
		name       string            // Name of the subtest
		body       any               // The body of the request
		googleOidc config.GoogleOidc // The environment google oidc
		tpRet1     *http.Response    // The response from transport
		tpRet2     error             // The error from transport
		exptRes    any               // The expected response
	}

	cases := []testCase{{
		name: "request ExchangeType is auth_code, return doesn't have access_token",
		body: authen.AuthenRequest{
			ExchangeType: "auth_code",
			Value:        "123",
		},
		googleOidc: config.GoogleOidc{
			ClientId:     "cid",
			ClientSecret: "cs",
			RedirectUri:  "url",
			IsDevMode:    false,
		},
		tpRet1: &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(anyToJsonMarshalToBytesBuffer(map[string]any{
				"field1":       "value",
				"access_token": "token",
			})),
			Header: make(http.Header),
		},
		tpRet2: nil,
		exptRes: map[string]any{
			"field1": "value",
		},
	}, {
		name: "request ExchangeType is refresh_token, return doesn't have access_token",
		body: authen.AuthenRequest{
			ExchangeType: "refresh_token",
			Value:        "123",
		},
		googleOidc: config.GoogleOidc{
			ClientId:     "cid",
			ClientSecret: "cs",
			RedirectUri:  "url",
			IsDevMode:    false,
		},
		tpRet1: &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(anyToJsonMarshalToBytesBuffer(map[string]any{
				"field1":       "value",
				"access_token": "token",
			})),
			Header: make(http.Header),
		},
		tpRet2: nil,
		exptRes: map[string]any{
			"field1": "value",
		},
	}, {
		name: "request ExchangeType is invalid",
		body: authen.AuthenRequest{
			ExchangeType: "invalid",
			Value:        "123",
		},
		googleOidc: config.GoogleOidc{
			ClientId:     "cid",
			ClientSecret: "cs",
			RedirectUri:  "url",
			IsDevMode:    false,
		},
		tpRet1:  &http.Response{},
		tpRet2:  nil,
		exptRes: "Wrong exchange type!",
	}, {
		name:       "error shouldbindjson",
		body:       "",
		googleOidc: config.GoogleOidc{},
		tpRet1:     &http.Response{},
		tpRet2:     nil,
		exptRes: map[string]any{
			"message": "json: cannot unmarshal string into Go value of type authen.AuthenRequest",
			"status":  string(app.Fail),
		},
	}, {
		name: "error http.client",
		body: authen.AuthenRequest{
			ExchangeType: "auth_code",
			Value:        "123",
		},
		googleOidc: config.GoogleOidc{
			ClientId:     "cid",
			ClientSecret: "cs",
			RedirectUri:  "url",
			IsDevMode:    false,
		},
		tpRet1: &http.Response{},
		tpRet2: errors.New("error"),
		exptRes: map[string]any{
			"message": "Post \"https://oauth2.googleapis.com/token\": error",
			"status":  string(app.Fail),
		},
	}, {
		name: "error json.unmarshal into tokenresultmap",
		body: authen.AuthenRequest{
			ExchangeType: "auth_code",
			Value:        "123",
		},
		googleOidc: config.GoogleOidc{
			ClientId:     "cid",
			ClientSecret: "cs",
			RedirectUri:  "url",
			IsDevMode:    false,
		},
		tpRet1: &http.Response{
			Body: io.NopCloser(anyToJsonMarshalToBytesBuffer("")),
		},
		tpRet2: nil,
		exptRes: map[string]any{
			"message": "json: cannot unmarshal string into Go value of type map[string]interface {}",
			"status":  string(app.Fail),
		},
	}}

	mockFnName := []string{"RoundTrip", "/auth/token", "/auth/token"}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Arrange

			transport := app.NewMockHttpTransport()
			transport.On(mockFnName[0], mock.Anything).Return(c.tpRet1, c.tpRet2)

			client := app.NewMockHttpClient(transport)

			h := authen.NewAuthenHandler(client, c.googleOidc)

			r, g := setupRouter()
			defer g()

			r.POST(mockFnName[1], h.ExchangeForTokens)

			jsonReqBody, _ := json.Marshal(c.body)

			req := httptest.NewRequest("POST", mockFnName[2], bytes.NewBuffer(jsonReqBody))
			w := httptest.NewRecorder()

			// Act
			r.ServeHTTP(w, req)

			var jsonRes any

			err := json.NewDecoder(w.Body).Decode(&jsonRes)

			log.Printf("%+v\n", jsonRes)
			log.Printf("%+v\n", w.Code)

			assert.NoError(t, err)
			if w.Code == http.StatusOK {
				assert.Equal(t, w.Code, http.StatusOK)
				transport.AssertNumberOfCalls(t, mockFnName[0], 1)
				assert.Equal(t, "https://oauth2.googleapis.com/token", transport.Calls[0].Arguments[0].(*http.Request).URL.String())
				switch c.body.(authen.AuthenRequest).ExchangeType {
				case "auth_code":
					assert.Equal(t, "client_id="+c.googleOidc.ClientId+"&client_secret="+c.googleOidc.ClientSecret+"&code="+c.body.(authen.AuthenRequest).Value+"&grant_type=authorization_code&redirect_uri="+c.googleOidc.RedirectUri, readerCloserToString(transport.Calls[0].Arguments[0].(*http.Request).Body))
				case "refresh_token":
					assert.Equal(t, "client_id=cid&client_secret=cs&grant_type=refresh_token&refresh_token=123", readerCloserToString(transport.Calls[0].Arguments[0].(*http.Request).Body))
				}
			} else if w.Code == http.StatusBadRequest {
				assert.Equal(t, w.Code, http.StatusBadRequest)
				transport.AssertNumberOfCalls(t, mockFnName[0], 0)
			}
			assert.Equal(t, c.exptRes, jsonRes)
		})
	}
}
