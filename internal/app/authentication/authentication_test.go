package authentication

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	user_model "go_shurtiner/internal/app/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUser(ctx context.Context, email string) (*user_model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user_model.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *user_model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type TestableAuth struct {
	*BasicAuth
	Token string
	User  *user_model.User
}

func TestGetUser(t *testing.T) {
	user := user_model.User{
		Name:  "Alexey",
		Email: "akozadaev@inbox.ru",
	}

	ctx := context.WithValue(context.Background(), User, &user)

	authUser, err := GetUser(ctx)

	require.NoError(t, err)
	require.NotNil(t, authUser)
	assert.Equal(t, &user, authUser)
}

func newTestableAuth(mockRepo *MockUserRepo, token string, user *user_model.User) *TestableAuth {
	ba := &BasicAuth{repo: mockRepo}

	return &TestableAuth{
		BasicAuth: ba,
		Token:     token,
		User:      user,
	}
}

// override methods

func (t *TestableAuth) HasAuthHeader(r *http.Request) bool {
	return true
}

func (t *TestableAuth) createTokenFromHeader(header string) (string, error) {
	if header == "Bearer invalid" {
		return "", errors.New("invalid_token")
	}
	return t.Token, nil
}

func (t *TestableAuth) parseToken(token string) (*user_model.User, error) {
	if token == "bad_token" {
		return nil, errors.New("bad_token")
	}
	return t.User, nil
}

// --- Actual Tests ---

func TestAuthenticate_NoAuthHeader(t *testing.T) {
	mockRepo := new(MockUserRepo)
	auth := &BasicAuth{repo: mockRepo}

	req, _ := http.NewRequest("GET", "/", nil)

	user, err := auth.Authenticate(req)
	assert.Nil(t, user)
	assert.EqualError(t, err, "header_not_found")
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepo)
	user := &user_model.User{Email: "user@example.com", Password: "password123"}

	tauth := newTestableAuth(mockRepo, "", user)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid")

	// подменим метод
	//tauth.BasicAuth.HasAuthHeader = tauth.HasAuthHeader
	//tauth.BasicAuth.createTokenFromHeader = tauth.createTokenFromHeader
	//tauth.BasicAuth.parseToken = tauth.parseToken

	res, err := tauth.Authenticate(req)
	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid_token")
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
	storedUser := &user_model.User{Email: "user@example.com", Password: string(hashed)}
	inputUser := &user_model.User{Email: "user@example.com", Password: "wrongpass"}

	mockRepo.On("GetUser", mock.Anything, inputUser.Email).Return(storedUser, nil)

	tauth := newTestableAuth(mockRepo, "valid_token", inputUser)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid")

	// подменим методы
	//tauth.BasicAuth.HasAuthHeader = tauth.HasAuthHeader
	//tauth.BasicAuth.createTokenFromHeader = tauth.createTokenFromHeader
	//tauth.BasicAuth.parseToken = tauth.parseToken

	user, err := tauth.Authenticate(req)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashedPassword")
}

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	password := "securepass"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	storedUser := &user_model.User{Email: "user@example.com", Password: string(hashed)}
	inputUser := &user_model.User{Email: "user@example.com", Password: password}

	mockRepo.On("GetUser", mock.Anything, inputUser.Email).Return(storedUser, nil)

	tauth := newTestableAuth(mockRepo, "valid_token", inputUser)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid")

	//tauth.BasicAuth.HasAuthHeader = tauth.HasAuthHeader
	//tauth.BasicAuth.createTokenFromHeader = tauth.createTokenFromHeader
	//tauth.BasicAuth.parseToken = tauth.parseToken

	user, err := tauth.Authenticate(req)
	assert.NoError(t, err)
	assert.Equal(t, storedUser, user)
}
