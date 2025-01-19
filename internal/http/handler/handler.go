package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleRequest(c *gin.Context, f func(c *gin.Context) *Response) {
	ctx := c.Request.Context()
	if _, ok := ctx.Deadline(); !ok {
		handleRequestReal(c, f(c))
		return
	}
	doneChan := make(chan *Response)
	defer close(doneChan)

	go func() {
		doneChan <- f(c)
	}()
	select {
	case <-ctx.Done():
		break
	case res := <-doneChan:
		handleRequestReal(c, res)
	}
}

func handleRequestReal(c *gin.Context, res *Response) {
	if res.Err == nil {
		statusCode := res.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		if res.Data != nil {
			c.JSON(res.StatusCode, res.Data)
		} else {
			c.Status(res.StatusCode)
		}
		return
	}

	var err *ErrorResponse
	ok := errors.As(res.Err, &err)
	if !ok {
		res.StatusCode = http.StatusInternalServerError
		err = &ErrorResponse{Errors: "An error has occurred, please try again later", Details: res.Err}
	}
	c.AbortWithStatusJSON(res.StatusCode, err)
}
