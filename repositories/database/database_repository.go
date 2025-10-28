package repositories

import (
	"asset-pulse-api/entities"
	"context"

	"gorm.io/gorm"
)

type databaseRepository struct {
	db *gorm.DB
}

//go:generate mockery --name=DatabaseRepository
type DatabaseRepository interface {
	GetUsers(ctx context.Context, companyCode *string, status *string, limit int, offset int) (*[]entities.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)
	CountUsers(ctx context.Context, companyCode *string, status *string) (int64, error)

	// Seat Optimization
	GetSeatOptimizationOpportunities(ctx context.Context, companyCode, departmentCode, appName string) (*[]entities.OptimizationOpportunity, error)

	// Purchase Templates
	CreatePurchaseTemplate(ctx context.Context, template *entities.PurchaseTemplate) (*entities.PurchaseTemplate, error)
	GetPurchaseTemplates(ctx context.Context, companyCode, departmentCode string) (*[]entities.PurchaseTemplate, error)

	// Cross-Subsidiary Match
	GetCrossSubMatch(ctx context.Context, companyCode, appName string, threshold float64) (*[]entities.SimilarApp, error)

	// User Licenses
	GetUserLicenseAssignments(ctx context.Context, userID int64) ([]UserLicenseAssignment, error)
}

// UserLicenseAssignment represents a user's license with aggregated data from multiple tables
type UserLicenseAssignment struct {
	ID               int64    `json:"id"`
	AppID            *int64   `json:"app_id"`
	AppName          *string  `json:"app_name"`
	AppAlias         *string  `json:"app_alias"`
	Category         *string  `json:"category"`
	LicenseTier      *string  `json:"license_tier"`
	AssignedAt       string   `json:"assigned_at"`
	ExpireDate       *string  `json:"expire_date"`
	EffectiveDate    *string  `json:"effective_date"`
	Cost             *float64 `json:"cost"`
	Currency         *string  `json:"currency"`
	UsageCount30Days int      `json:"usage_count_30d"`
	LastUsedAt       *string  `json:"last_used_at"`
}

func (d *databaseRepository) GetUsers(ctx context.Context, companyCode *string, status *string, limit int, offset int) (*[]entities.User, error) {
	users := []entities.User{}
	query := d.db.WithContext(ctx)

	if companyCode != nil && *companyCode != "" {
		query = query.Where("company_code = ?", *companyCode)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return &users, nil
}

func (d *databaseRepository) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	user := entities.User{}
	result := d.db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (d *databaseRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user := entities.User{}
	result := d.db.WithContext(ctx).Where("username = ? OR email = ?", username, username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (d *databaseRepository) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	var roles []string
	result := d.db.WithContext(ctx).
		Table("user_roles ur").
		Select("r.name").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Pluck("r.name", &roles)

	if result.Error != nil {
		return nil, result.Error
	}

	return roles, nil
}

func (d *databaseRepository) CountUsers(ctx context.Context, companyCode *string, status *string) (int64, error) {
	var count int64
	query := d.db.WithContext(ctx).Model(&entities.User{})

	if companyCode != nil && *companyCode != "" {
		query = query.Where("company_code = ?", *companyCode)
	}

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (d *databaseRepository) GetSeatOptimizationOpportunities(ctx context.Context, companyCode, departmentCode, appName string) (*[]entities.OptimizationOpportunity, error) {
	opportunities := []entities.OptimizationOpportunity{}
	query := d.db.WithContext(ctx)

	if companyCode != "" {
		query = query.Where("company_code = ?", companyCode)
	}
	if departmentCode != "" {
		query = query.Where("department = ?", departmentCode)
	}
	if appName != "" {
		query = query.Where("app_name = ?", appName)
	}

	result := query.Order("potential_savings DESC").Find(&opportunities)
	if result.Error != nil {
		return nil, result.Error
	}

	return &opportunities, nil
}

func (d *databaseRepository) CreatePurchaseTemplate(ctx context.Context, template *entities.PurchaseTemplate) (*entities.PurchaseTemplate, error) {
	result := d.db.WithContext(ctx).Create(template)
	if result.Error != nil {
		return nil, result.Error
	}

	return template, nil
}

func (d *databaseRepository) GetPurchaseTemplates(ctx context.Context, companyCode, departmentCode string) (*[]entities.PurchaseTemplate, error) {
	templates := []entities.PurchaseTemplate{}
	query := d.db.WithContext(ctx)

	if companyCode != "" {
		query = query.Where("company_code = ?", companyCode)
	}
	if departmentCode != "" {
		query = query.Where("department = ?", departmentCode)
	}

	result := query.Order("created_at DESC").Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	return &templates, nil
}

func (d *databaseRepository) GetCrossSubMatch(ctx context.Context, companyCode, appName string, threshold float64) (*[]entities.SimilarApp, error) {
	similarApps := []entities.SimilarApp{}
	query := d.db.WithContext(ctx)

	if companyCode != "" {
		query = query.Where("company_code != ?", companyCode)
	}
	if appName != "" {
		query = query.Where("app_name ILIKE ?", "%"+appName+"%")
	}
	if threshold > 0 {
		query = query.Where("similarity_score >= ?", threshold)
	}

	result := query.Order("similarity_score DESC").Find(&similarApps)
	if result.Error != nil {
		return nil, result.Error
	}

	return &similarApps, nil
}

func (d *databaseRepository) GetUserLicenseAssignments(ctx context.Context, userID int64) ([]UserLicenseAssignment, error) {
	var licenses []UserLicenseAssignment

	query := `
		SELECT 
			la.id,
			la.app_id,
			la.license_tier,
			la.assigned_at,
			a.name as app_name,
			a.alias as app_alias,
			a.category,
			li.expire_date,
			li.effective_date,
			pb.list_price as cost,
			pb.currency,
			COALESCE(
				(SELECT COUNT(*) 
				 FROM usage_events ue 
				 WHERE ue.app_id = la.app_id 
				   AND ue.user_id = la.user_id 
				   AND ue.event_at > NOW() - INTERVAL '30 days'), 0) as usage_count_30d,
			(SELECT MAX(event_at)::text
			 FROM usage_events ue 
			 WHERE ue.app_id = la.app_id 
			   AND ue.user_id = la.user_id) as last_used_at
		FROM license_assignments la
		JOIN apps a ON a.id = la.app_id
		LEFT JOIN license_inventories li ON li.id = la.license_id
		LEFT JOIN price_books pb ON pb.app_id = la.app_id AND pb.tier = la.license_tier
		WHERE la.user_id = $1
		  AND la.revoked_at IS NULL
		ORDER BY la.assigned_at DESC
	`

	result := d.db.WithContext(ctx).Raw(query, userID).Scan(&licenses)
	if result.Error != nil {
		return nil, result.Error
	}

	return licenses, nil
}

func New(db *gorm.DB) *databaseRepository {
	return &databaseRepository{
		db: db,
	}
}
