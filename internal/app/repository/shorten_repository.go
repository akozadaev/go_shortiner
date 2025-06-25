package repository

import (
	"context"
	"go_shurtiner/internal/adapter"
	"go_shurtiner/internal/app/model"
	database2 "go_shurtiner/internal/database"
	"go_shurtiner/pkg/logging"
	"gorm.io/gorm"
)

//go:generate mockery --name=ShortenRepository --output=../../../test/mocks --outpkg=mocks --filename=mock_shorten_repository.go
type ShortenRepository interface {
	SaveLink(ctx context.Context, link *model.Link) error
	FindLink(ctx context.Context, shortened string) (model.Link, error)
	FetchLinks(ctx context.Context) ([]model.Link, error)
}

func NewShortenRepository(db *gorm.DB) ShortenRepository {
	return &shortenRepository{db: db}
}

type shortenRepository struct {
	db *gorm.DB
}

func (s shortenRepository) SaveLink(ctx context.Context, links *model.Link) error {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, s.db)

	if err := db.WithContext(ctx).Create(links).Error; err != nil {
		logger.Errorw("failed to save link", "err", err)
		if database2.IsKeyConflictErr(err) {
			return database2.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (s shortenRepository) FindLink(ctx context.Context, shortened string) (model.Link, error) {
	logger := logging.FromContext(ctx)
	db := database2.FromContext(ctx, s.db)
	var err error
	var links model.Link
	if err = db.WithContext(ctx).First(&links, "shortened = ?", shortened).Error; err != nil {
		logger.Errorw("failed to get link", "err", err)
	}

	return links, err
}

func (s shortenRepository) FetchLinks(ctx context.Context) ([]model.Link, error) {
	db := database2.FromContext(ctx, s.db)
	links := make([]model.Link, 0)

	var err error
	if pagination, ok := ctx.Value(adapter.Pagination).(*adapter.PaginationAdapter); ok {
		err = db.Find(&links).Limit(int(pagination.GetLimit())).Error
	} else {
		err = db.Find(&links).Error
	}

	return links, err
}
