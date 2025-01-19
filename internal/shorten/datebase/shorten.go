package datebase

import (
	"context"
	"go_shurtiner/internal/database/database"
	"go_shurtiner/internal/shorten/model"
	"go_shurtiner/pkg/logging"
	"gorm.io/gorm"
)

type ShortenRepository interface {
	SaveLink(ctx context.Context, link *model.Link) error
	FindLink(ctx context.Context, shortened string) (model.Link, error)
}

func NewShortenRepository(db *gorm.DB) ShortenRepository {
	return &shortenRepository{db: db}
}

type shortenRepository struct {
	db *gorm.DB
}

func (s shortenRepository) SaveLink(ctx context.Context, links *model.Link) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, s.db)

	if err := db.WithContext(ctx).Create(links).Error; err != nil {
		logger.Errorw("failed to save link", "err", err)
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (s shortenRepository) FindLink(ctx context.Context, shortened string) (model.Link, error) {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, s.db)
	var err error
	var links model.Link
	if err = db.WithContext(ctx).First(&links, "shortened = ?", shortened).Error; err != nil {
		logger.Errorw("failed to get link", "err", err)
	}
	links.Shortened = links.Shortened
	return links, err
}
