package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const goroutineTimeout = 5 * time.Second

func HandleRequest(c *gin.Context, f func(c *gin.Context) *Response) {
	ctx := c.Request.Context()
	if _, ok := ctx.Deadline(); !ok {
		handleRequestReal(c, f(c))
		return
	}

	doneChan := make(chan *Response, 1)

	go func() {
		defer close(doneChan)
		select {
		case <-ctx.Done():
			return
		case doneChan <- f(c):
		}
	}()

	select {
	case <-ctx.Done():
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	case res := <-doneChan:
		if res != nil {
			handleRequestReal(c, res)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	case <-time.After(goroutineTimeout):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request processing timeout"})
	}
}

func handleRequestReal(c *gin.Context, res *Response) {
	if res.Err == nil {
		statusCode := res.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		if res.Data != nil {
			c.JSON(statusCode, res.Data)
		} else {
			c.Status(statusCode)
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
