package datebase

import (
	"context"
	"go_shurtiner/internal/shorten/model"
	"gorm.io/gorm"
)

type ShortenRepository interface {
	SaveLink(ctx context.Context, link *model.Link) error
	FindLink(ctx context.Context, source string) model.Link
}

func NewShortenRepository(db *gorm.DB) ShortenRepository {
	return &shortenRepository{db: db}
}

type shortenRepository struct {
	db *gorm.DB
}

func (s shortenRepository) SaveLink(ctx context.Context, links *model.Link) error {
	//TODO implement me
	panic("implement me")
}

func (s shortenRepository) FindLink(ctx context.Context, source string) model.Link {
	//TODO implement me
	panic("implement me")
}
