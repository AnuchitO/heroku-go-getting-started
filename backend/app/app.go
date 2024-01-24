package app

import "time"

type status string

const (
	Success status = "success"
	Fail    status = "error"
)

const (
	storeErrorStatus = 450
)

type Response struct {
	Status  status `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorField struct {
	Field string `json:"field"`
	Value any    `json:"value"`
	Tag   string `json:"tag"`
}

type Error Response

func (err *Error) Error() string {
	return err.Message
}

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
