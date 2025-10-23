# Asset Pulse API Setup Guide

## Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
# Application Settings
APP_ENV=development
APP_PORT=8080
APP_NAME=asset-pulse-api

# HTTP Settings
HTTP_TIMEOUT=30
HTTP_RETRY_COUNT=3

# JWT Settings
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Database Settings
DB_HOST=localhost
DB_PORT=5432
DB_NAME=asset_pulse
DB_USERNAME=postgres
DB_PASSWORD=password
SSL_MODE=disable
SSL_ROOT_CERT=
TIME_ZONE=Asia/Bangkok

# AI Service Settings (for future use)
OPENAI_API_KEY=your-openai-api-key
ANTHROPIC_API_KEY=your-anthropic-api-key

# Azure Settings (for future use)
AZURE_CLIENT_ID=your-azure-client-id
AZURE_CLIENT_SECRET=your-azure-client-secret
AZURE_TENANT_ID=your-azure-tenant-id
```

## Database Setup

1. Create PostgreSQL database:
```sql
CREATE DATABASE asset_pulse;
```

2. Run the SQL scripts in order:
```bash
psql -d asset_pulse -f scripts/0001_initial_tables.sql
psql -d asset_pulse -f scripts/0002_insert_mockdata.sql
psql -d asset_pulse -f scripts/0003_create_view_recommend_pack_for_user.sql
```

## Running the API

1. Install dependencies:
```bash
go mod tidy
```

2. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/me` - Get current user

### AI Services
- `POST /api/v1/ai/recommendations/jd-match` - Generate JD recommendations
- `POST /api/v1/ai/consolidation/memo` - Generate consolidation memo
- `POST /api/v1/ai/similarity` - Calculate software similarity

### Role-based Dashboards
- `GET /api/v1/employee/dashboard` - Employee dashboard
- `GET /api/v1/manager/dashboard` - Manager dashboard
- `GET /api/v1/cto/dashboard` - CTO dashboard
- `GET /api/v1/group-cto/dashboard` - Group CTO dashboard

### Users
- `GET /api/v1/users` - Get users list

## Features Implemented

✅ **Authentication & JWT**
- JWT token generation and validation
- Password hashing with bcrypt
- Role-based access control
- Authentication middleware

✅ **AI Services (Mock)**
- JD Recommendation service
- Consolidation memo generation
- Software similarity calculation
- Ready for real AI integration

✅ **Role-based Routing**
- Employee, Manager, CTO, Group CTO routes
- Permission-based access control
- Company-level access control

✅ **Database Integration**
- GORM ORM integration
- PostgreSQL connection
- User authentication tables
- Role management

## Next Steps

1. **Database Connection**: Ensure PostgreSQL is running and accessible
2. **Real AI Integration**: Replace mock AI services with OpenAI/Anthropic APIs
3. **Frontend Integration**: Connect Next.js frontend to these API endpoints
4. **Testing**: Add unit tests and integration tests
5. **Deployment**: Set up Docker containers and CI/CD pipeline
