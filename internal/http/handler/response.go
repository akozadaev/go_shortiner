package handler

type Response struct {
	StatusCode int
	Data       interface{}
	Err        error
}

func NewSuccessResponse(statusCode int, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, error string, details interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Err: &ErrorResponse{
			Errors:  error,
			Details: details,
		},
	}
}

func NewInternalErrorResponse(err error) *Response {
	return &Response{
		Err: err,
	}
}
