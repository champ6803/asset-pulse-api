package mocks

import (
	"asset-pulse-api/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockDatabaseRepository is a mock implementation of DatabaseRepository
type MockDatabaseRepository struct {
	mock.Mock
}

func (m *MockDatabaseRepository) GetUsers(ctx context.Context, companyCode, status *string, limit, offset int) (*[]entities.User, error) {
	args := m.Called(ctx, companyCode, status, limit, offset)
	return args.Get(0).(*[]entities.User), args.Error(1)
}

func (m *MockDatabaseRepository) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockDatabaseRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockDatabaseRepository) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDatabaseRepository) CountUsers(ctx context.Context, companyCode, status *string) (int64, error) {
	args := m.Called(ctx, companyCode, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDatabaseRepository) GetSeatOptimizationOpportunities(ctx context.Context, companyCode, departmentCode, appName string) (*[]entities.OptimizationOpportunity, error) {
	args := m.Called(ctx, companyCode, departmentCode, appName)
	return args.Get(0).(*[]entities.OptimizationOpportunity), args.Error(1)
}

func (m *MockDatabaseRepository) CreatePurchaseTemplate(ctx context.Context, template *entities.PurchaseTemplate) (*entities.PurchaseTemplate, error) {
	args := m.Called(ctx, template)
	return args.Get(0).(*entities.PurchaseTemplate), args.Error(1)
}

func (m *MockDatabaseRepository) GetPurchaseTemplates(ctx context.Context, companyCode, departmentCode string) (*[]entities.PurchaseTemplate, error) {
	args := m.Called(ctx, companyCode, departmentCode)
	return args.Get(0).(*[]entities.PurchaseTemplate), args.Error(1)
}

func (m *MockDatabaseRepository) GetCrossSubMatch(ctx context.Context, companyCode, appName string, threshold float64) (*[]entities.SimilarApp, error) {
	args := m.Called(ctx, companyCode, appName, threshold)
	return args.Get(0).(*[]entities.SimilarApp), args.Error(1)
}
