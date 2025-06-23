package repository

import (
	"context"
	"go_shurtiner/internal/adapter"
	"go_shurtiner/internal/app/model"
	database2 "go_shurtiner/internal/database"
	"go_shurtiner/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepository --output=../../../test/mocks --outpkg=mocks --filename=mock_user_repository.go
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserForApi(ctx context.Context, email string) (*model.UserApi, error)
	GetUserForApiById(ctx context.Context, id string) (*model.UserApi, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
	HashPassword(password string) (string, error)
	FetchUsers(ctx context.Context) ([]model.UserApi, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r userRepository) CreateUser(ctx context.Context, user *model.User) error {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, r.db)

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Errorw("failed to save user", "err", err)
		if database2.IsKeyConflictErr(err) {
			return database2.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (r userRepository) GetUserForApi(ctx context.Context, email string) (*model.UserApi, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, r.db)
	var err error
	var user model.User
	if err = db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		logger.Errorw("failed to get user for API", "err", err)
	}
	result := model.UserApi{
		user.Model,
		user.Name,
		user.LastName,
		user.MiddleName,
		user.Email}
	return &result, err
}

func (r userRepository) GetUserForApiById(ctx context.Context, id string) (*model.UserApi, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, r.db)
	var err error
	var user model.User
	if err = db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		logger.Errorw("failed to get user for API by id", "err", err)
	}
	result := model.UserApi{
		user.Model,
		user.Name,
		user.LastName,
		user.MiddleName,
		user.Email}
	return &result, err
}

func (r userRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, r.db)
	var err error
	var user model.User
	if err = db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		logger.Errorw("failed to get user", "err", err)
	}

	return &user, err
}

func (r userRepository) FetchUsers(ctx context.Context) ([]model.UserApi, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("get all users")

	db := database2.FromContext(ctx, r.db)
	users := make([]model.User, 0)
	result := make([]model.UserApi, 0)
	var err error
	if pagination, ok := ctx.Value(adapter.Pagination).(*adapter.PaginationAdapter); ok {
		err = db.Find(&users).Limit(int(pagination.GetLimit())).Error
	} else {
		err = db.Find(&users).Error
	}

	for _, user := range users {
		result = append(result, model.UserApi{
			user.Model,
			user.Name,
			user.LastName,
			user.MiddleName,
			user.Email})
	}

	return result, err
}

func (r userRepository) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r userRepository) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
