package entities

import (
	"time"

	"github.com/google/uuid"
)

// Organization layer entities
type Org struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Code        string     `gorm:"unique;not null" json:"code"`
	Name        string     `gorm:"not null" json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by"`
	Companies   []Company  `gorm:"foreignKey:OrgID" json:"companies,omitempty"`
}

type Company struct {
	ID          int64        `gorm:"primaryKey" json:"id"`
	OrgID       *int64       `json:"org_id"`
	Code        string       `gorm:"unique;not null" json:"code"`
	Name        string       `gorm:"not null" json:"name"`
	Description *string      `json:"description"`
	CreatedAt   time.Time    `gorm:"default:now()" json:"created_at"`
	CreatedBy   *int64       `json:"created_by"`
	UpdatedAt   *time.Time   `json:"updated_at"`
	UpdatedBy   *int64       `json:"updated_by"`
	Org         *Org         `gorm:"foreignKey:OrgID" json:"org,omitempty"`
	Departments []Department `gorm:"foreignKey:CompanyID" json:"departments,omitempty"`
}

type Department struct {
	ID            int64      `gorm:"primaryKey" json:"id"`
	CompanyID     *int64     `json:"company_id"`
	Code          string     `gorm:"unique;not null" json:"code"`
	Name          string     `gorm:"not null" json:"name"`
	Description   *string    `json:"description"`
	ManagerUserID *int64     `json:"manager_user_id"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy     *int64     `json:"created_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by"`
	Company       *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Manager       *User      `gorm:"foreignKey:ManagerUserID" json:"manager,omitempty"`
}

type Role struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	CompanyCode *string    `json:"company_code"`
	Key         string     `gorm:"unique;not null" json:"key"`
	Name        string     `gorm:"not null" json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by"`
}

type User struct {
	ID             int64      `gorm:"primaryKey" json:"id"`
	Username       *string    `gorm:"unique" json:"username"`
	CompanyCode    *string    `json:"company_code"`
	DepartmentCode *string    `json:"department_code"`
	EntraID        *uuid.UUID `gorm:"unique" json:"entra_id"`
	Email          *string    `gorm:"unique" json:"email"`
	Password       *string    `json:"password"`
	Salt           *string    `json:"salt"`
	DisplayName    *string    `json:"display_name"`
	Title          *string    `json:"title"`
	EmployeeID     *string    `json:"employee_id"`
	Status         *string    `json:"status"`
	CreatedAt      time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy      *int64     `json:"created_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedBy      *int64     `json:"updated_by"`
	UserRoles      []UserRole `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

type UserRole struct {
	ID         int64      `gorm:"primaryKey" json:"id"`
	UserID     int64      `gorm:"not null" json:"user_id"`
	RoleID     int64      `gorm:"not null" json:"role_id"`
	ScopeLevel string     `gorm:"not null" json:"scope_level"`
	ScopeRefID *int64     `json:"scope_ref_id"`
	AssignedAt time.Time  `gorm:"default:now()" json:"assigned_at"`
	AssignedBy *int64     `json:"assigned_by"`
	RevokedAt  *time.Time `json:"revoked_at"`
	RevokedBy  *int64     `json:"revoked_by"`
	CreatedAt  time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy  *int64     `json:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by"`
	User       *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role       *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// Application entities
type App struct {
	ID                 int64        `gorm:"primaryKey" json:"id"`
	CompanyCode        *string      `json:"company_code"`
	Key                *string      `gorm:"unique" json:"key"`
	Name               string       `gorm:"not null" json:"name"`
	Alias              *string      `json:"alias"`
	Status             *string      `json:"status"`
	Description        *string      `json:"description"`
	Category           *string      `json:"category"`
	DataClassification *string      `json:"data_classification"`
	ApplicationTier    *string      `json:"application_tier"`
	ASVL               *string      `json:"asvl"`
	ServiceHour        *string      `json:"service_hour"`
	ServiceLifeYears   *int         `json:"service_life_years"`
	OwnerDept          *string      `json:"owner_dept"`
	OwnerBusiness      *string      `json:"owner_business"`
	OwnerIT            *string      `json:"owner_it"`
	UserGroup          *string      `json:"user_group"`
	LifecycleStatus    *string      `json:"lifecycle_status"`
	CreatedAt          time.Time    `gorm:"default:now()" json:"created_at"`
	CreatedBy          *int64       `json:"created_by"`
	UpdatedAt          *time.Time   `json:"updated_at"`
	UpdatedBy          *int64       `json:"updated_by"`
	AppFeatures        []AppFeature `gorm:"foreignKey:AppID" json:"app_features,omitempty"`
}

type AppFeature struct {
	ID          int64   `gorm:"primaryKey" json:"id"`
	AppID       int64   `gorm:"not null" json:"app_id"`
	FeatureKey  *string `json:"feature_key"`
	FeatureName *string `json:"feature_name"`
	Description *string `json:"description"`
	App         *App    `gorm:"foreignKey:AppID" json:"app,omitempty"`
}

// Vendor and Contract entities
type Vendor struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	TaxID     *string    `json:"tax_id"`
	Country   *string    `json:"country"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy *int64     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int64     `json:"updated_by"`
}

type Contract struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	CompanyCode      *string    `json:"company_code"`
	VendorID         *int64     `json:"vendor_id"`
	ContractNumber   *string    `json:"contract_number"`
	Title            *string    `json:"title"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	Status           *string    `json:"status"`
	ParentContractID *int64     `json:"parent_contract_id"`
	Currency         *string    `json:"currency"`
	CreatedAt        time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy        *int64     `json:"created_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int64     `json:"updated_by"`
	Vendor           *Vendor    `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
}

// License entities
type LicenseInventory struct {
	ID            int64      `gorm:"primaryKey" json:"id"`
	CompanyCode   *string    `json:"company_code"`
	AppID         *int64     `json:"app_id"`
	ContractID    *int64     `json:"contract_id"`
	LicenseTier   *string    `json:"license_tier"`
	TotalSeats    *int       `json:"total_seats"`
	ReservedSeats *int       `gorm:"default:0" json:"reserved_seats"`
	EffectiveDate *time.Time `json:"effective_date"`
	ExpireDate    *time.Time `json:"expire_date"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy     *int64     `json:"created_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by"`
	App           *App       `gorm:"foreignKey:AppID" json:"app,omitempty"`
	Contract      *Contract  `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

type LicenseAssignment struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	CompanyCode      *string    `json:"company_code"`
	UserID           *int64     `json:"user_id"`
	AppID            *int64     `json:"app_id"`
	LicenseTier      *string    `json:"license_tier"`
	AssignmentSource *string    `json:"assignment_source"`
	AssignedAt       time.Time  `gorm:"default:now()" json:"assigned_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
	Reason           *string    `json:"reason"`
	CreatedAt        time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy        *int64     `json:"created_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int64     `json:"updated_by"`
	User             *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	App              *App       `gorm:"foreignKey:AppID" json:"app,omitempty"`
}

// Recommendation entities
type OptimizationOpportunity struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	AppName          string     `gorm:"not null" json:"app_name"`
	Department       string     `gorm:"not null" json:"department"`
	Action           string     `gorm:"not null" json:"action"` // revoke, reallocate, downgrade
	InactiveUsers    int        `json:"inactive_users"`
	PendingRequests  int        `json:"pending_requests"`
	CanReallocate    int        `json:"can_reallocate"`
	PotentialSavings float64    `json:"potential_savings"`
	RiskLevel        string     `json:"risk_level"`
	Rationale        string     `json:"rationale"`
	CompanyCode      string     `json:"company_code"`
	CreatedAt        time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type PurchaseTemplate struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description *string    `json:"description"`
	Apps        string     `gorm:"type:text" json:"apps"` // JSON array
	Cost        float64    `json:"cost"`
	Department  string     `json:"department"`
	CompanyCode string     `json:"company_code"`
	TimesUsed   int        `gorm:"default:0" json:"times_used"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy   *int64     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by"`
}

type SimilarApp struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	AppName         string     `gorm:"not null" json:"app_name"`
	CompanyCode     string     `gorm:"not null" json:"company_code"`
	SimilarityScore float64    `json:"similarity_score"`
	JaccardIndex    float64    `json:"jaccard_index"`
	LLMScore        float64    `json:"llm_score"`
	Reasoning       string     `gorm:"type:text" json:"reasoning"`
	CreatedAt       time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
