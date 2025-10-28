package enums

const (
	ErrorMsgTemplate = "Code: %d | Message: %s | Cause: %s"
)

// User Status
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

// License Assignment Status
const (
	LicenseStatusActive  = "active"
	LicenseStatusRevoked = "revoked"
	LicenseStatusExpired = "expired"
	LicenseStatusPending = "pending"
)

// App Status
const (
	AppStatusActive   = "Active"
	AppStatusInactive = "Inactive"
)

// Scope Level
const (
	ScopeLevelGroup      = "group"
	ScopeLevelSubsidiary = "subsidiary"
	ScopeLevelDepartment = "department"
	ScopeLevelApp        = "app"
	ScopeLevelUser       = "user"
)

// Context Keys
const (
	ContextKeyUserID         = "user_id"
	ContextKeyUsername       = "username"
	ContextKeyEmail          = "email"
	ContextKeyCompanyCode    = "company_code"
	ContextKeyDepartmentCode = "department_code"
	ContextKeyRole           = "role"
)
