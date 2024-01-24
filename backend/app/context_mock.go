package app

import "net/http"

type MockAppContext struct {
	Context
	Params       map[string]string
	Email        string
	ResponseData interface{}
	ResponseCode int
	// Other necessary fields and methods
}

func (m *MockAppContext) GetString(key string) string {
	return m.Email
}

func (m *MockAppContext) Param(key string) string {
	return m.Params[key]
}

func (m *MockAppContext) OK(v interface{}) {
	m.ResponseData = v
	m.ResponseCode = http.StatusOK
}

func (m *MockAppContext) BadRequest(err error) {
	m.ResponseData = err
	m.ResponseCode = http.StatusBadRequest
}

func (m *MockAppContext) StoreError(err error) {
	m.ResponseData = err
	m.ResponseCode = 450
}
