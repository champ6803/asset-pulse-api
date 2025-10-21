# Asset Pulse API - Project Structure

## Overview
This project follows the same clean architecture pattern as `datax-ai-app-invx-cs-assistant-api`, adapted for the Asset Pulse system based on the SQL schema in the scripts folder.

## Directory Structure

```
asset-pulse-api/
├── configs/                    # Configuration management
│   └── config.go              # App config struct with env bindings
│
├── entities/                   # Database entities (GORM models)
│   └── database_entity.go     # All database entities from SQL schema
│       ├── Org, Company, Department
│       ├── Role, User, UserRole
│       ├── App, AppFeature
│       ├── Vendor, Contract
│       └── LicenseInventory, LicenseAssignment
│
├── handler/                    # HTTP handlers (controllers)
│   ├── dto/                   # Data transfer objects
│   │   ├── base.response.go  # Base response structure
│   │   └── users.dto.go      # User-specific DTOs
│   ├── get_users_handler.go  # GET /users handler
│   └── route.go              # Route definitions
│
├── repositories/               # Data access layer
│   └── database/
│       └── database_repository.go  # Database operations
│           ├── GetUsers()
│           ├── GetUserByID()
│           └── CountUsers()
│
├── usecase/                    # Business logic layer
│   ├── models/                # Usecase models
│   │   └── users.model.go    # User usecase models
│   ├── get_users.go          # Get users business logic
│   └── usecase.go            # Usecase interface
│
├── utils/                      # Utility packages
│   ├── apperrs/              # Application error handling
│   │   └── errors.go
│   ├── gorm/                 # Database utilities
│   │   └── pgClient.go       # PostgreSQL connection
│   ├── logger/               # Logging utilities
│   │   └── logger.go
│   └── transformer/          # Response transformers
│       └── transformer.go
│
├── scripts/                    # SQL migration scripts
│   ├── 0001_initial_tables.sql
│   ├── 0002_insert_mockdata.sql
│   └── 0003_create_view_recommend_pack_for_user.sql
│
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── Dockerfile                  # Docker image definition
├── docker-compose.yaml         # Docker compose configuration
├── Makefile                    # Build and run commands
├── .env.example               # Environment variables template
├── .gitignore                 # Git ignore rules
└── README.md                  # Project documentation

```

## Key Components

### 1. Entities Layer (`entities/`)
- Database models mapped from SQL schema
- Uses GORM tags for database mapping
- Includes all tables: users, apps, contracts, licenses, etc.
- Schema-qualified table names (asset_pulse.*)

### 2. Repository Layer (`repositories/database/`)
- Interface-based design for testability
- Database operations abstraction
- Context-aware queries
- Example methods:
  - `GetUsers()` - with filtering and pagination
  - `GetUserByID()` - get single user
  - `CountUsers()` - total count for pagination

### 3. Usecase Layer (`usecase/`)
- Business logic implementation
- Input/output models separate from entities
- Error handling with custom AppError
- Pagination logic
- Example: `GetUsers()` with validation and transformation

### 4. Handler Layer (`handler/`)
- HTTP request/response handling
- Input validation
- Query parameter parsing
- Error recovery middleware
- DTOs for API responses

### 5. Utils Layer (`utils/`)
- **apperrs**: Custom error types with status codes
- **gorm**: Database connection management
- **logger**: Structured logging
- **transformer**: Response formatting

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### Users
```
GET /api/v1/users?company_code=SCBX&status=active&page=1&page_size=10
```

Query Parameters:
- `company_code` (optional) - Filter by company
- `status` (optional) - Filter by status (active/inactive)
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page

Response:
```json
{
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "company_code": "SCBX",
        "department_code": "IT",
        "email": "user@example.com",
        "display_name": "John Doe",
        "title": "Developer",
        "employee_id": "EMP001",
        "status": "active"
      }
    ],
    "total": 100
  }
}
```

## Differences from Reference Project

1. **Removed AI/Azure Dependencies**: No AI Core service integration
2. **Simplified Authentication**: No Azure Identity required
3. **Schema-Qualified Tables**: Uses `asset_pulse.*` schema
4. **Focus on CRUD**: RESTful API design vs. chat-based
5. **Pagination**: Built-in pagination support

## Database Setup

1. Create database:
```bash
createdb asset_pulse
```

2. Run migrations:
```bash
psql -U postgres -d asset_pulse -f scripts/0001_initial_tables.sql
psql -U postgres -d asset_pulse -f scripts/0002_insert_mockdata.sql
psql -U postgres -d asset_pulse -f scripts/0003_create_view_recommend_pack_for_user.sql
```

## Running the Application

### Using Go directly:
```bash
cp .env.example .env
# Edit .env with your database credentials
go run main.go
```

### Using Make:
```bash
make build    # Build binary
make run      # Run application
make test     # Run tests
make clean    # Clean artifacts
```

### Using Docker:
```bash
docker-compose up -d
```

## Next Steps to Extend

1. **Add More Endpoints**:
   - Apps management (GET, POST, PUT, DELETE /api/v1/apps)
   - License management
   - Contract management
   - Recommendations

2. **Add Authentication**:
   - JWT middleware
   - User session management
   - Role-based access control

3. **Add Testing**:
   - Unit tests for usecase
   - Repository tests with mocks
   - Integration tests

4. **Add Validation**:
   - Request validation middleware
   - Custom validators

5. **Add More Features**:
   - Search and filtering
   - Export to CSV/Excel
   - Bulk operations
   - Analytics endpoints

## Design Patterns Used

- **Clean Architecture**: Separation of concerns
- **Repository Pattern**: Data access abstraction
- **Dependency Injection**: Through constructors
- **Interface Segregation**: Small, focused interfaces
- **Error Wrapping**: Context-aware error handling

## Technologies

- **Go 1.21**: Programming language
- **Gin**: HTTP web framework
- **GORM**: ORM library
- **PostgreSQL**: Database
- **Docker**: Containerization

