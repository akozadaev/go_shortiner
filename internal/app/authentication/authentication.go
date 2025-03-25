package authentication

import (
	"context"
	"errors"
	user_model "go_shurtiner/internal/app/model"
	"net/http"
)

const User = "user"

type Authentication interface {
	Authenticate(r *http.Request) (*user_model.User, error)
	UnauthorizedResponse(err error) UnauthorizedResponse
	HasAuthHeader(r *http.Request) bool
}

type UnauthorizedResponse struct {
	Errors map[string]string `json:"errors"`
}

func GetUser(ctx context.Context) (*user_model.User, error) {
	currentUser, ok := ctx.Value(User).(*user_model.User)
	if !ok {
		return nil, errors.New("context has no user")
	}

	return currentUser, nil
}
