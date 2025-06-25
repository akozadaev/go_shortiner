package repository

import (
	"context"
	"errors"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenRepository_SaveLink_Success(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	link := &model.Link{
		Source:    "https://example.com",
		Shortened: "abc123",
	}

	mockRepo.On("SaveLink", ctx, link).Return(nil)

	err := mockRepo.SaveLink(ctx, link)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestShortenRepository_SaveLink_Failure(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	link := &model.Link{
		Source:    "https://example.org",
		Shortened: "xyz456",
	}
	expectedErr := errors.New("db error")

	mockRepo.On("SaveLink", ctx, link).Return(expectedErr)

	err := mockRepo.SaveLink(ctx, link)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestShortenRepository_FindLink_Success(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	shortened := "abc123"
	expectedLink := model.Link{
		Source:    "https://example.net",
		Shortened: shortened,
	}

	mockRepo.On("FindLink", ctx, shortened).Return(expectedLink, nil)

	link, err := mockRepo.FindLink(ctx, shortened)
	assert.NoError(t, err)
	assert.Equal(t, expectedLink, link)
	mockRepo.AssertExpectations(t)
}

func TestShortenRepository_FindLink_NotFound(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	shortened := "notfound"
	expectedErr := errors.New("record not found")

	mockRepo.On("FindLink", ctx, shortened).Return(model.Link{}, expectedErr)

	link, err := mockRepo.FindLink(ctx, shortened)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, link.ID)
	mockRepo.AssertExpectations(t)
}

func TestShortenRepository_FetchLinks_Success(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	links := []model.Link{
		{Source: "https://1.com", Shortened: "a1"},
		{Source: "https://2.com", Shortened: "a2"},
	}

	mockRepo.On("FetchLinks", ctx).Return(links, nil)

	result, err := mockRepo.FetchLinks(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, links, result)
	mockRepo.AssertExpectations(t)
}

func TestShortenRepository_FetchLinks_Error(t *testing.T) {
	mockRepo := mocks.NewShortenRepository(t)

	ctx := context.Background()
	expectedErr := errors.New("fetch failed")

	mockRepo.On("FetchLinks", ctx).Return(nil, expectedErr)

	result, err := mockRepo.FetchLinks(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
