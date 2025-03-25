package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
	user_model "go_shurtiner/internal/app/model"
	repository "go_shurtiner/internal/app/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type BasicAuth struct {
	repo repository.UserRepository
}

func NewBasicAuth(repo repository.UserRepository) Authentication {
	return &BasicAuth{
		repo: repo,
	}
}

func (b *BasicAuth) HasAuthHeader(r *http.Request) bool {
	header := r.Header.Get("Authorization")
	if header == "" {
		return false
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "Basic" {
		return false
	}

	return true
}

func (b *BasicAuth) Authenticate(r *http.Request) (*user_model.User, error) {
	if !b.HasAuthHeader(r) {
		return nil, errors.New("header_not_found")
	}

	header := r.Header.Get("Authorization")
	fmt.Println("header: ")
	fmt.Println(header)
	token, err := b.createTokenFromHeader(header)
	if err != nil {
		return nil, err
	}
	requestedUser, err := b.parseToken(token)
	if err != nil {
		return nil, err
	}
	user, err := b.repo.GetUser(r.Context(), requestedUser.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestedUser.Password))
	if err != nil {
		return nil, err
	}

	return user, err
}

func (b *BasicAuth) createTokenFromHeader(header string) (string, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", errors.New("invalid_token")
	}

	token := parts[1]

	return token, nil
}

func (b *BasicAuth) parseToken(token string) (*user_model.AuthUserRequest, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid_token")
	}

	return &user_model.AuthUserRequest{
		Email:    parts[0],
		Password: parts[1],
	}, nil
}

func (b *BasicAuth) UnauthorizedResponse(err error) UnauthorizedResponse {
	return UnauthorizedResponse{
		Errors: map[string]string{
			"basic_token": err.Error(),
		},
	}
}
