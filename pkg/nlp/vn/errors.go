package vn

import "strings"

// // UniqueConstraintError raised when the db complains about violation of unique constraint value
// type UniqueConstraintError struct {
// 	msg string
// 	err error
// }

// func (e *UniqueConstraintError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // ServerNotRespondingError raised when the NLP server did not respond
// type ServerNotRespondingError struct {
// 	msg string
// 	err error
// }

// func (e *ServerNotRespondingError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // OversizedTextError raised when the server complains that the text is too long to be annotated
// type OversizedTextError struct {
// 	msg string
// 	err error
// }

// func (e *OversizedTextError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // JSONParsingError raised when there is a problem with un-marshalling the JSON obj to a struct type
// type JSONParsingError struct {
// 	msg string
// 	err error
// }

// func (e *JSONParsingError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // ServerResponseError raised when the server returns an error
// type ServerResponseError struct {
// 	msg string
// 	err error
// }

// func (e *ServerResponseError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // NLPClientError raised when the client is not working
// type NLPClientError struct {
// 	msg string
// 	err error
// }

// func (e *NLPClientError) Error() string {
// 	return e.msg + e.err.Error()
// }

// // NilAnnotationError raised when the annotation returned as nil, meaning that the text might be too large
// type NilAnnotationError struct {
// 	msg string
// }

// func (e *NilAnnotationError) Error() string {
// 	return e.msg
// }

// ----

// ErrorType is the types of error monitored for Vn NLP
type ErrorType string

// List of error types
const (
	ErrorTypeServerNotResponding       ErrorType = "server_not_responding"
	ErrorTypeServerError               ErrorType = "server_error"
	ErrorTypeRequestFailed             ErrorType = "request_failed"
	ErrorTypeNilAnnotation             ErrorType = "nil_annotation"
	ErrorTypeUniqueConstraintViolation ErrorType = "unique_constraint_violation"
)

// Error is the custom error for VN NLP
type Error struct {
	Type ErrorType
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	return strings.Join([]string{string(e.Type), e.Msg, e.Err.Error()}, ": ")
}
