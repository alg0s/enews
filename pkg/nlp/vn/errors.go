package vn

import "strings"

// ErrorType is the types of error monitored for Vn NLP
type ErrorType string

// List of error types
const (
	ErrorTypeServerNotResponding       ErrorType = "server_not_responding"
	ErrorTypeServerError               ErrorType = "server_error"
	ErrorTypeRequestFailed             ErrorType = "request_failed"
	ErrorTypeTextTooLong               ErrorType = "text_too_long"
	ErrorTypeUniqueConstraintViolation ErrorType = "unique_constraint_violation"
	ErrorTypeResetServer               ErrorType = "reset_server"
)

// Error is the custom error for VN NLP
type Error struct {
	Type ErrorType
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return strings.Join([]string{string(e.Type), e.Msg}, ": ")
	}
	return strings.Join([]string{string(e.Type), e.Msg, e.Err.Error()}, ": ")
}
