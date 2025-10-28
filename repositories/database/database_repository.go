package repositories

import (
	"asset-pulse-api/entities"
	"context"
	"fmt"

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
	GetUserLicenseAssignments(ctx context.Context, userID int64, search, status string) ([]UserLicenseAssignment, error)

	// Seat Optimization Analytics
	GetLicenseUsageAnalytics(ctx context.Context, companyCode, departmentCode, appName string) ([]LicenseUsageAnalytic, error)
	GetInactiveUsers(ctx context.Context, appID int64, days int) ([]InactiveUserAnalytic, error)
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

// LicenseUsageAnalytic represents license usage analytics data
type LicenseUsageAnalytic struct {
	AppID          int64    `json:"app_id"`
	AppName        string   `json:"app_name"`
	AppCategory    string   `json:"app_category"`
	LicenseTier    *string  `json:"license_tier"`
	DepartmentCode string   `json:"department_code"`
	DepartmentName string   `json:"department_name"`
	TotalUsers     int      `json:"total_users"`
	InactiveUsers  int      `json:"inactive_users"`
	AvgCostPerUser *float64 `json:"avg_cost_per_user"`
	TotalCost      *float64 `json:"total_cost"`
}

// InactiveUserAnalytic represents inactive user data
type InactiveUserAnalytic struct {
	UserID         int64   `json:"user_id"`
	DisplayName    *string `json:"display_name"`
	DepartmentCode string  `json:"department_code"`
	DepartmentName string  `json:"department_name"`
	AppID          int64   `json:"app_id"`
	AppName        string  `json:"app_name"`
	LicenseTier    *string `json:"license_tier"`
	LastUsedAt     string  `json:"last_used_at"`
	DaysInactive   int     `json:"days_inactive"`
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

func (d *databaseRepository) GetUserLicenseAssignments(ctx context.Context, userID int64, search, status string) ([]UserLicenseAssignment, error) {
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
	`

	args := []interface{}{userID}
	argIdx := 2

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(" AND a.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Add status filter (based on expire_date)
	if status != "" && status != "all" {
		if status == "active" {
			query += " AND (li.expire_date IS NULL OR li.expire_date > NOW() + INTERVAL '30 days')"
		} else if status == "expiring" {
			query += " AND li.expire_date IS NOT NULL AND li.expire_date <= NOW() + INTERVAL '30 days' AND li.expire_date > NOW()"
		} else if status == "expired" {
			query += " AND li.expire_date IS NOT NULL AND li.expire_date <= NOW()"
		}
	}

	query += " ORDER BY la.assigned_at DESC"

	result := d.db.WithContext(ctx).Raw(query, args...).Scan(&licenses)
	if result.Error != nil {
		return nil, result.Error
	}

	return licenses, nil
}

func (d *databaseRepository) GetLicenseUsageAnalytics(ctx context.Context, companyCode, departmentCode, appName string) ([]LicenseUsageAnalytic, error) {
	var analytics []LicenseUsageAnalytic

	query := `
		SELECT 
			a.id as app_id,
			a.name as app_name,
			COALESCE(a.category, 'Other') as app_category,
			la.license_tier,
			d.code as department_code,
			d.name as department_name,
			COUNT(DISTINCT la.user_id) as total_users,
			COUNT(DISTINCT CASE WHEN ue.id IS NULL THEN la.user_id END) as inactive_users,
			AVG(pb.list_price) as avg_cost_per_user,
			SUM(pb.list_price) as total_cost
		FROM license_assignments la
		JOIN apps a ON a.id = la.app_id
		JOIN users u ON u.id = la.user_id
		JOIN departments d ON d.code = u.department_code
		LEFT JOIN usage_events ue ON ue.app_id = la.app_id 
			AND ue.user_id = la.user_id 
			AND ue.event_at > NOW() - INTERVAL '90 days'
		LEFT JOIN license_inventories li ON li.id = la.license_id
		LEFT JOIN price_books pb ON pb.app_id = la.app_id AND pb.tier = la.license_tier
		WHERE la.revoked_at IS NULL
	`

	args := []interface{}{}
	argIdx := 1

	if companyCode != "" {
		query += fmt.Sprintf(" AND la.company_code = $%d", argIdx)
		args = append(args, companyCode)
		argIdx++
	}

	if departmentCode != "" {
		query += fmt.Sprintf(" AND d.code = $%d", argIdx)
		args = append(args, departmentCode)
		argIdx++
	}

	if appName != "" {
		query += fmt.Sprintf(" AND a.name ILIKE $%d", argIdx)
		args = append(args, "%"+appName+"%")
		argIdx++
	}

	query += " GROUP BY a.id, a.name, a.category, la.license_tier, d.code, d.name ORDER BY total_cost DESC"

	result := d.db.WithContext(ctx).Raw(query, args...).Scan(&analytics)
	if result.Error != nil {
		return nil, result.Error
	}

	return analytics, nil
}

func (d *databaseRepository) GetInactiveUsers(ctx context.Context, appID int64, days int) ([]InactiveUserAnalytic, error) {
	var inactiveUsers []InactiveUserAnalytic

	query := `
		SELECT 
			la.user_id,
			u.display_name,
			COALESCE(u.department_code, 'Unknown') as department_code,
			COALESCE(d.name, 'Unknown') as department_name,
			la.app_id,
			a.name as app_name,
			la.license_tier,
			COALESCE(MAX(ue.event_at)::text, la.assigned_at::text) as last_used_at,
			EXTRACT(DAY FROM (NOW() - COALESCE(MAX(ue.event_at), la.assigned_at)))::int as days_inactive
		FROM license_assignments la
		JOIN users u ON u.id = la.user_id
		LEFT JOIN departments d ON d.code = u.department_code
		JOIN apps a ON a.id = la.app_id
		LEFT JOIN usage_events ue ON ue.app_id = la.app_id 
			AND ue.user_id = la.user_id
		WHERE la.revoked_at IS NULL
			AND ($1 = 0 OR la.app_id = $1)
		GROUP BY la.user_id, u.display_name, u.department_code, d.name, 
				 la.app_id, a.name, la.license_tier, la.assigned_at
		HAVING COALESCE(MAX(ue.event_at), la.assigned_at) < NOW() - INTERVAL '90 days'
			AND EXTRACT(DAY FROM (NOW() - COALESCE(MAX(ue.event_at), la.assigned_at))) >= $2
		ORDER BY days_inactive DESC
	`

	result := d.db.WithContext(ctx).Raw(query, appID, days).Scan(&inactiveUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return inactiveUsers, nil
}

func New(db *gorm.DB) *databaseRepository {
	return &databaseRepository{
		db: db,
	}
}
