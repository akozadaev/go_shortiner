package auth

import (
	"context"
	"errors"
	"go_shurtiner/internal/app/authentication"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
	user := &model.User{
		Email: "john@example.com",
		Name:  "John",
	}

	ctx := context.WithValue(context.Background(), authentication.User, user)

	getUser, err := authentication.GetUser(ctx)

	assert.NoError(t, err)
	assert.Equal(t, user, getUser)
}

func TestGetUser_NoUserInContext(t *testing.T) {
	ctx := context.Background() // без user

	gotUser, err := authentication.GetUser(ctx)

	assert.Nil(t, gotUser)
	assert.Error(t, err)
	assert.Equal(t, "context has no user", err.Error())
}

func TestAuthentication_Authenticate_Success(t *testing.T) {
	mockAuth := mocks.NewAuthentication(t)

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	expectedUser := &model.User{
		Email: "akozadaev@inbox.ru",
		Name:  "Алексей",
	}

	mockAuth.
		On("Authenticate", req).
		Return(expectedUser, nil)

	user, err := mockAuth.Authenticate(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockAuth.AssertExpectations(t)
}

func TestAuthentication_Authenticate_Failure(t *testing.T) {
	mockAuth := mocks.NewAuthentication(t)

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	mockErr := errors.New("invalid token")

	mockAuth.
		On("Authenticate", req).
		Return(nil, mockErr)

	user, err := mockAuth.Authenticate(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, mockErr, err)
	mockAuth.AssertExpectations(t)
}

func TestAuthentication_UnauthorizedResponse(t *testing.T) {
	mockAuth := mocks.NewAuthentication(t)

	err := errors.New("unauthorized")
	expected := authentication.UnauthorizedResponse{
		Errors: map[string]string{
			"auth": "unauthorized",
		},
	}

	mockAuth.
		On("UnauthorizedResponse", err).
		Return(expected)

	result := mockAuth.UnauthorizedResponse(err)

	assert.Equal(t, expected, result)
	mockAuth.AssertExpectations(t)
}

func TestAuthentication_HasAuthHeader(t *testing.T) {
	mockAuth := mocks.NewAuthentication(t)

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("Authorization", "Bearer abc.def.ghi")

	mockAuth.
		On("HasAuthHeader", req).
		Return(true)

	ok := mockAuth.HasAuthHeader(req)

	assert.True(t, ok)
	mockAuth.AssertExpectations(t)
}
