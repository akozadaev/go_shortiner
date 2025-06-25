package repository

import (
	"context"
	"errors"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/test/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReportRepository_PrepareReportData_Success(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()

	expectedLinks := []model.Link{
		{Source: "https://example.com", Shortened: "ex123"},
		{Source: "https://example.org", Shortened: "go456"},
	}

	mockRepo.On("PrepareReportData", ctx).Return(expectedLinks, nil)

	links, err := mockRepo.PrepareReportData(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedLinks, links)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_PrepareReportData_Error(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()
	expectedErr := errors.New("db failure")

	mockRepo.On("PrepareReportData", ctx).Return(nil, expectedErr)

	links, err := mockRepo.PrepareReportData(ctx)
	assert.Error(t, err)
	assert.Nil(t, links)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_SaveReportData_Success(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()

	report := &model.PreparedReport{
		Timestamp:    time.Now(),
		Source:       "https://example.com",
		Shortened:    "ex123",
		UserEmail:    "john@example.com",
		UserFullName: "John Doe",
	}

	mockRepo.On("SaveReportData", ctx, report).Return(nil)

	err := mockRepo.SaveReportData(ctx, report)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_SaveReportData_Error(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()

	report := &model.PreparedReport{
		Timestamp:    time.Now(),
		Source:       "https://fail.com",
		Shortened:    "fail123",
		UserEmail:    "fail@example.com",
		UserFullName: "Failure User",
	}
	expectedErr := errors.New("insert failed")

	mockRepo.On("SaveReportData", ctx, report).Return(expectedErr)

	err := mockRepo.SaveReportData(ctx, report)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_GetReportData_Success(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()
	startDate := time.Now().Add(-24 * time.Hour)

	expectedReport := &[]model.PreparedReport{
		{
			ID:           1,
			CreatedAt:    time.Now().Add(-23 * time.Hour),
			Timestamp:    time.Now().Add(-23 * time.Hour),
			Source:       "https://example.com",
			Shortened:    "ex123",
			UserEmail:    "john@example.com",
			UserFullName: "John Doe",
		},
		{
			ID:           2,
			CreatedAt:    time.Now().Add(-22 * time.Hour),
			Timestamp:    time.Now().Add(-22 * time.Hour),
			Source:       "https://another.com",
			Shortened:    "an456",
			UserEmail:    "jane@example.com",
			UserFullName: "Jane Smith",
		},
	}

	mockRepo.On("GetReportData", ctx, startDate).Return(expectedReport, nil)

	result, err := mockRepo.GetReportData(ctx, startDate)
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, result)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_GetReportData_Error(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()
	startDate := time.Now()

	expectedErr := errors.New("query error")

	mockRepo.On("GetReportData", ctx, startDate).Return((*[]model.PreparedReport)(nil), expectedErr)

	result, err := mockRepo.GetReportData(ctx, startDate)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_CreateReport_Success(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()
	startDate := time.Now()

	expected := &[]model.PreparedReport{
		{
			ID:           10,
			Timestamp:    startDate,
			Source:       "https://gen1.com",
			Shortened:    "g1",
			UserEmail:    "a@example.com",
			UserFullName: "User A",
		},
		{
			ID:           11,
			Timestamp:    startDate,
			Source:       "https://gen2.com",
			Shortened:    "g2",
			UserEmail:    "b@example.com",
			UserFullName: "User B",
		},
	}

	mockRepo.On("CreateReport", ctx, startDate).Return(expected, nil)

	result, err := mockRepo.CreateReport(ctx, startDate)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestReportRepository_CreateReport_Error(t *testing.T) {
	mockRepo := mocks.NewReportRepository(t)
	ctx := context.Background()
	startDate := time.Now()

	expectedErr := errors.New("generation error")

	mockRepo.On("CreateReport", ctx, startDate).Return((*[]model.PreparedReport)(nil), expectedErr)

	result, err := mockRepo.CreateReport(ctx, startDate)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
