package app

import (
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockClock struct {
	mock.Mock
}

func NewMockClock() *mockClock {
	return &mockClock{}
}

func (m *mockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *mockClock) After(d time.Duration) <-chan time.Time {
	args := m.Called(d)
	return args.Get(0).(<-chan time.Time)
}

type mockHttpRoundTrip struct {
	mock.Mock
}

// Interface: http.RoundTripper
func (m *mockHttpRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func NewMockHttpTransport() *mockHttpRoundTrip {
	return &mockHttpRoundTrip{}
}

func NewMockHttpClient(tr http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: tr,
	}
}
