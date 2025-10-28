package repositories

import (
	"asset-pulse-api/entities"
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	// License Management
	GetLicenses(ctx context.Context, userID *int64, companyCode, status, search, category, licenseTier *string) ([]entities.LicenseWithInventory, error)
	GetLicensesWithLimit(ctx context.Context, userID *int64, companyCode, status, search, category, licenseTier *string, limit int) ([]entities.LicenseWithInventory, int64, error)

	// Request Management
	GetPendingRequests(ctx context.Context, userID *int64, companyCode *string, limit int) ([]entities.PendingRequestWithDetails, int64, error)

	// AI Recommendations
	GetUserJobProfile(ctx context.Context, userID int64) (*entities.JobProfile, error)
	GetActiveApps(ctx context.Context, companyCode *string) ([]entities.App, error)
	GetAppPrices(ctx context.Context, appIDs []int64) (map[int64]float64, error)

	// Grouped software license
	GetAllSoftwareLicenses(ctx context.Context) ([]entities.SoftwareLicense, error)
	InsertSoftwareLicense(ctx context.Context, license *entities.SoftwareLicense) error
	DeleteSoftwareLicense(ctx context.Context, id uint) error
	GetCurrentGroupedSoftware(ctx context.Context) ([]entities.CurrentGroupedSoftware, error)
	UpsertCurrentGroupedSoftware(ctx context.Context, jsonData []byte) error
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

func (d *databaseRepository) GetLicenses(ctx context.Context, userID *int64, companyCode, status, search, category, licenseTier *string) ([]entities.LicenseWithInventory, error) {
	var results []entities.LicenseWithInventory

	// Base query with type-safe selection
	query := d.db.WithContext(ctx).
		Table("license_assignments la").
		Select(`
			la.id as license_assignment_id,
			li.id as license_inventory_id,
			a.id as app_id,
			a.name as app_name,
			a.category as app_category,
			a.status as app_status,
			la.license_tier,
			la.assigned_at,
			la.revoked_at,
			CASE WHEN la.revoked_at IS NOT NULL THEN true ELSE false END as is_revoked,
			li.total_seats,
			li.reserved_seats,
			li.effective_date,
			li.expire_date,
			ue.last_used,
			ue.usage_frequency
		`).
		Joins("INNER JOIN apps a ON la.app_id = a.id").
		Joins("LEFT JOIN license_inventories li ON la.license_id = li.id").
		Joins(`LEFT JOIN (
			SELECT 
				app_id,
				MAX(event_at) as last_used,
				LEAST(100, ROUND((COUNT(DISTINCT DATE(event_at))::numeric / 30.0) * 100)) as usage_frequency
			FROM usage_events
			WHERE event_at >= NOW() - INTERVAL '30 days'
			GROUP BY app_id
		) ue ON a.id = ue.app_id`)

	// REQUIRED: Filter by user_id
	if userID == nil {
		return nil, fmt.Errorf("user_id is required")
	}
	query = query.Where("la.user_id = ?", *userID)

	// Filter by company_code if provided
	if companyCode != nil && *companyCode != "" {
		query = query.Where("la.company_code = ?", *companyCode)
	}

	// Filter by status
	if status != nil && *status != "" {
		switch *status {
		case "active":
			query = query.Where("la.revoked_at IS NULL")
		case "revoked":
			query = query.Where("la.revoked_at IS NOT NULL")
		case "expired":
			query = query.Where("la.revoked_at IS NULL AND li.expire_date < NOW()")
		}
	}

	// Search by app name
	if search != nil && *search != "" {
		searchPattern := "%" + *search + "%"
		query = query.Where("a.name ILIKE ?", searchPattern)
	}

	// Filter by category
	if category != nil && *category != "" {
		query = query.Where("a.category = ?", *category)
	}

	// Filter by license tier
	if licenseTier != nil && *licenseTier != "" {
		query = query.Where("la.license_tier = ?", *licenseTier)
	}

	// Order by most recent assignments
	query = query.Order("la.assigned_at DESC")

	// Execute query
	result := query.Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

func (d *databaseRepository) GetLicensesWithLimit(ctx context.Context, userID *int64, companyCode, status, search, category, licenseTier *string, limit int) ([]entities.LicenseWithInventory, int64, error) {
	var results []entities.LicenseWithInventory
	var total int64

	// REQUIRED: Filter by user_id
	if userID == nil {
		return nil, 0, fmt.Errorf("user_id is required")
	}

	// Build base query for counting unique apps (without Select)
	baseQuery := d.db.WithContext(ctx).
		Table("license_assignments la").
		Joins("LEFT JOIN apps a ON la.app_id = a.id").
		Joins("LEFT JOIN license_inventories li ON la.app_id = li.app_id AND la.license_tier = li.license_tier").
		Where("la.user_id = ?", *userID).
		Where("a.status = ?", "Active").
		Where("la.revoked_at IS NULL")

	// Filter by company_code if provided
	if companyCode != nil && *companyCode != "" {
		baseQuery = baseQuery.Where("la.company_code = ?", *companyCode)
	}

	// Count total unique apps (use DISTINCT on app_id)
	var countResult struct {
		Count int64
	}
	countQuery := baseQuery.Select("COUNT(DISTINCT a.id) as count").Scan(&countResult)
	if countQuery.Error != nil {
		return nil, 0, countQuery.Error
	}
	total = countResult.Count

	// Build query with DISTINCT ON to get one row per app (latest assignment)
	query := d.db.WithContext(ctx).
		Table("license_assignments la").
		Select(`DISTINCT ON (a.id)
			la.id as license_assignment_id,
			li.id as license_inventory_id,
			a.id as app_id,
			a.name as app_name,
			a.category as app_category,
			a.status as app_status,
			la.license_tier,
			la.assigned_at,
			la.revoked_at,
			CASE WHEN la.revoked_at IS NOT NULL THEN true ELSE false END as is_revoked,
			li.total_seats,
			li.reserved_seats,
			li.effective_date,
			li.expire_date,
			ue.last_used,
			ue.usage_frequency
		`).
		Joins("LEFT JOIN apps a ON la.app_id = a.id").
		Joins("LEFT JOIN license_inventories li ON la.app_id = li.app_id AND la.license_tier = li.license_tier").
		Joins(`LEFT JOIN (
			SELECT 
				app_id,
				MAX(event_at) as last_used,
				LEAST(100, ROUND((COUNT(DISTINCT DATE(event_at))::numeric / 30.0) * 100)) as usage_frequency
			FROM usage_events
			WHERE event_at >= NOW() - INTERVAL '30 days'
			GROUP BY app_id
		) ue ON a.id = ue.app_id`).
		Where("la.user_id = ?", *userID).
		Where("a.status = ?", "Active").
		Where("la.revoked_at IS NULL")

	// Apply same filters as count query
	if companyCode != nil && *companyCode != "" {
		query = query.Where("la.company_code = ?", *companyCode)
	}

	// Order by app_id first (required for DISTINCT ON), then by assigned_at DESC
	query = query.Order("a.id, la.assigned_at DESC")

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query
	result := query.Find(&results)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return results, total, nil
}

func (d *databaseRepository) GetPendingRequests(ctx context.Context, userID *int64, companyCode *string, limit int) ([]entities.PendingRequestWithDetails, int64, error) {
	var results []entities.PendingRequestWithDetails
	var total int64

	// Build query with ticket number generation
	query := d.db.WithContext(ctx).Debug().Table("requests r").
		Select(`
			r.id,
			CONCAT('#REQ-', EXTRACT(YEAR FROM r.created_at)::text, '-', LPAD(r.id::text, 3, '0')) AS ticket_no,
			r.company_code,
			r.type,
			r.requester_user_id,
			r.scope_level,
			r.scope_ref_id,
			r.payload_json,
			r.status,
			r.created_at,
			r.updated_at,
			COALESCE(u.username, u.email) AS requester_username,
			u.email AS requester_email,
			u.display_name AS requester_full_name,
			u.department_code AS requester_department_code,
			rs.step_no AS current_step_no,
			rs.approver_role AS current_approver_role,
			rs.approver_user_id AS current_approver_user_id,
			COALESCE(approver.username, approver.email) AS approver_username,
			approver.display_name AS approver_name,
			rs.sla_due_at AS step_sla_due_at,
			total_steps.total_steps,
			a.id AS app_id,
			a.name AS app_name
		`).
		Joins("LEFT JOIN users u ON r.requester_user_id = u.id").
		Joins("LEFT JOIN request_steps rs ON r.id = rs.request_id AND rs.status = 'pending'").
		Joins("LEFT JOIN users approver ON rs.approver_user_id = approver.id").
		Joins("LEFT JOIN (SELECT request_id, COUNT(*) as total_steps FROM request_steps GROUP BY request_id) total_steps ON r.id = total_steps.request_id").
		Joins("LEFT JOIN purchase_template_items pti ON pti.template_id = CAST(r.payload_json->>'template_id' AS BIGINT)").
		Joins("LEFT JOIN apps a ON pti.app_id = a.id").
		Where("r.status = ?", "pending")

	// Filter by user ID if provided
	if userID != nil {
		query = query.Where("r.requester_user_id = ?", *userID)
	}

	// Filter by company code if provided
	if companyCode != nil && *companyCode != "" {
		query = query.Where("r.company_code = ?", *companyCode)
	}

	// Count total before applying limit
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply limit and order
	query = query.Order("r.created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query
	result := query.Find(&results)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return results, total, nil
}

func (d *databaseRepository) GetUserJobProfile(ctx context.Context, userID int64) (*entities.JobProfile, error) {
	var jobProfile entities.JobProfile

	result := d.db.WithContext(ctx).
		Table("job_profiles jp").
		Select("jp.*").
		Joins("INNER JOIN user_job_profiles ujp ON jp.id = ujp.job_profile_id").
		Where("ujp.user_id = ?", userID).
		First(&jobProfile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &jobProfile, nil
}

func (d *databaseRepository) GetActiveApps(ctx context.Context, companyCode *string) ([]entities.App, error) {
	var apps []entities.App

	// Build query starting from license_assignments and license_inventories, then join apps
	query := d.db.WithContext(ctx).
		Table("license_assignments la").
		Select("DISTINCT a.*").
		Joins("INNER JOIN apps a ON la.app_id = a.id").
		Joins("LEFT JOIN license_inventories li ON la.app_id = li.app_id AND la.license_tier = li.license_tier").
		Where("a.status = ?", "Active").
		Where("la.revoked_at IS NULL")

	if companyCode != nil && *companyCode != "" {
		query = query.Where("la.company_code = ?", *companyCode)
	}

	result := query.Order("a.name ASC").Find(&apps)
	if result.Error != nil {
		return nil, result.Error
	}

	return apps, nil
}

func (d *databaseRepository) GetAppPrices(ctx context.Context, appIDs []int64) (map[int64]float64, error) {
	if len(appIDs) == 0 {
		return make(map[int64]float64), nil
	}

	var prices []struct {
		AppID     int64
		ListPrice float64
	}

	// Get current valid prices for the apps
	result := d.db.WithContext(ctx).
		Table("price_books").
		Select("app_id, list_price").
		Where("app_id IN ?", appIDs).
		Where("unit = ?", "seat").
		Where("(valid_from IS NULL OR valid_from <= CURRENT_DATE)").
		Where("(valid_to IS NULL OR valid_to >= CURRENT_DATE)").
		Order("app_id, valid_from DESC").
		Find(&prices)

	if result.Error != nil {
		return nil, result.Error
	}

	// Map app_id to price (taking the first/most recent valid price for each app)
	priceMap := make(map[int64]float64)
	for _, p := range prices {
		if _, exists := priceMap[p.AppID]; !exists {
			priceMap[p.AppID] = p.ListPrice
		}
	}

	return priceMap, nil

}

func (r *databaseRepository) GetAllSoftwareLicenses(ctx context.Context) ([]entities.SoftwareLicense, error) {
	var licenses []entities.SoftwareLicense
	if err := r.db.WithContext(ctx).Find(&licenses).Error; err != nil {
		return nil, err
	}
	return licenses, nil
}

func (r *databaseRepository) InsertSoftwareLicense(ctx context.Context, license *entities.SoftwareLicense) error {
	license.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(license).Error
}

func (r *databaseRepository) DeleteSoftwareLicense(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.SoftwareLicense{}, id).Error
}

func (r *databaseRepository) GetCurrentGroupedSoftware(ctx context.Context) ([]entities.CurrentGroupedSoftware, error) {
	var output entities.RawCurrentGroupedSoftware
	if err := r.db.WithContext(ctx).First(&output, 1).Error; err != nil {
		return nil, err
	}

	var groups []entities.CurrentGroupedSoftware
	if err := json.Unmarshal(output.JSONData, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *databaseRepository) UpsertCurrentGroupedSoftware(ctx context.Context, jsonData []byte) error {
	var groups []entities.CurrentGroupedSoftware
	if err := json.Unmarshal(jsonData, &groups); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	var output entities.RawCurrentGroupedSoftware
	err := r.db.WithContext(ctx).First(&output, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			output = entities.RawCurrentGroupedSoftware{
				ID:        1,
				JSONData:  jsonData,
				UpdatedAt: time.Now(),
			}
			return r.db.WithContext(ctx).Create(&output).Error
		}
		return err
	}

	output.JSONData = jsonData
	output.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(&output).Error
}

func New(db *gorm.DB) *databaseRepository {
	return &databaseRepository{
		db: db,
	}
}
