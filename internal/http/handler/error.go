package handler

import "fmt"

type ErrorCode string

const (
	InvalidQueryValue   = ErrorCode("InvalidQueryValue")
	InvalidUriValue     = ErrorCode("InvalidUriValue")
	InvalidBodyValue    = ErrorCode("InvalidBodyValue")
	NotFoundEntity      = ErrorCode("NotFoundEntity")
	DuplicateEntry      = ErrorCode("DuplicateEntry")
	InternalServerError = ErrorCode("InternalServerError")
)

type ErrorResponse struct {
	Errors  string      `json:"errors"`
	Details interface{} `json:"details"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse{Errors:%s, Details:%v}", e.Errors, e.Details)
}
