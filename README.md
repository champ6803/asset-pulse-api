# Asset Pulse API

Asset Pulse API is a comprehensive application and license management system built with Go and the Gin framework.

## Features

- User management and authentication
- Application inventory tracking
- License management and assignment
- Contract and vendor management
- Usage tracking and analytics
- Recommendation engine for license optimization

## Project Structure

```
asset-pulse-api/
├── configs/          # Configuration files
├── entities/         # Database entities/models
├── handler/          # HTTP handlers (controllers)
│   └── dto/         # Data transfer objects
├── repositories/     # Data access layer
│   └── database/    # Database repository implementations
├── usecase/          # Business logic layer
│   └── models/      # Usecase models
├── utils/            # Utility packages
│   ├── apperrs/     # Application error handling
│   ├── gorm/        # Database utilities
│   ├── logger/      # Logging utilities
│   └── transformer/ # Response transformers
├── scripts/          # SQL migration scripts
├── main.go          # Application entry point
└── go.mod           # Go module dependencies
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher

## Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE asset_pulse;
```

2. Run the migration scripts in order:
```bash
psql -U postgres -d asset_pulse -f scripts/0001_initial_tables.sql
psql -U postgres -d asset_pulse -f scripts/0002_insert_mockdata.sql
psql -U postgres -d asset_pulse -f scripts/0003_create_view_recommend_pack_for_user.sql
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd asset-pulse-api
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file and configure it:
```bash
cp .env.example .env
```

4. Edit `.env` with your database credentials and configuration.

5. **Important**: Set up your Azure OpenAI credentials:
   - Create an Azure OpenAI resource in the [Azure Portal](https://portal.azure.com)
   - Deploy a GPT-4 model in your Azure OpenAI resource
   - Get your API key from Keys and Endpoint section
   - Add it to your `.env` file:
     ```bash
     AZURE_OPENAI_KEY=your-actual-api-key-here
     AZURE_OPENAI_ENDPOINT=https://your-resource-name.openai.azure.com/
     AZURE_OPENAI_MODEL=gpt-4
     ```
   - The application uses Azure OpenAI for AI-powered features including:
     - Job description-based software recommendations
     - Consolidation memo generation
     - Software similarity analysis
     - App recommendation scoring
     - Short description generation

## Running the Application

### Development

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Production Build

```bash
go build -o asset-pulse-api
./asset-pulse-api
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check API health status

### Users
- `GET /api/v1/users` - Get list of users
  - Query parameters:
    - `company_code` (optional) - Filter by company code
    - `status` (optional) - Filter by status (active/inactive)
    - `page` (optional, default: 1) - Page number
    - `page_size` (optional, default: 10) - Items per page

## Example Usage

### Get Users
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

### Get Users by Company
```bash
curl -X GET "http://localhost:8080/api/v1/users?company_code=SCBX&status=active"
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_ENV | Application environment | development |
| APP_PORT | Server port | 8080 |
| APP_NAME | Application name | asset-pulse-api |
| JWT_SECRET | Secret key for JWT token generation | - |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_NAME | Database name | asset_pulse |
| DB_USERNAME | Database username | postgres |
| DB_PASSWORD | Database password | postgres |
| SSL_MODE | Database SSL mode | disable |
| TIME_ZONE | Application timezone | Asia/Bangkok |
| **AZURE_OPENAI_KEY** | **Azure OpenAI API key** | **Required** |
| **AZURE_OPENAI_ENDPOINT** | **Azure OpenAI endpoint URL** | **Required** |
| **AZURE_OPENAI_MODEL** | **Azure OpenAI model deployment name** | **gpt-4** |

## Development

### Code Organization

This project follows clean architecture principles:

1. **Entities Layer**: Database models and entities
2. **Repository Layer**: Data access and persistence
3. **Usecase Layer**: Business logic and orchestration
4. **Handler Layer**: HTTP request handling and routing

### Adding New Endpoints

1. Create entity in `entities/database_entity.go`
2. Add repository methods in `repositories/database/`
3. Implement business logic in `usecase/`
4. Create handler in `handler/`
5. Register route in `handler/route.go`

## Testing

```bash
go test ./...
```

## License

[Your License Here]

## Contributing

[Contributing Guidelines]
