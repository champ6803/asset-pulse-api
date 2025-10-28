# Asset Pulse API - Postman Collection

This directory contains Postman collections and environments for testing the Asset Pulse API.

## Files

- **Asset_Pulse_API.postman_collection.json** - Complete API collection with all endpoints
- **Asset_Pulse_Local.postman_environment.json** - Local development environment
- **Asset_Pulse_Production.postman_environment.json** - Production environment

## How to Import

### Import Collection
1. Open Postman
2. Click **Import** button
3. Select `Asset_Pulse_API.postman_collection.json`
4. Click **Import**

### Import Environment
1. Click the **Environments** tab in Postman
2. Click **Import**
3. Select the environment file (Local or Production)
4. Click **Import**

## Usage

### 1. Select Environment
- Click the environment dropdown in the top-right corner
- Select **Asset Pulse - Local** for local development
- Select **Asset Pulse - Production** for production

### 2. Authentication Flow

#### Step 1: Login
1. Open the **Authentication** folder
2. Select **Login** request
3. Update the request body with your credentials:
   ```json
   {
       "username": "employee@scb.com",
       "password": "password"
   }
   ```
4. Click **Send**
5. The token will be automatically saved to the environment variable `auth_token`

#### Step 2: Use Protected Endpoints
All protected endpoints will automatically use the `{{auth_token}}` variable in the Authorization header.

### 3. Available Endpoints

#### Authentication
- **POST** `/api/v1/auth/login` - Login and get JWT token
- **POST** `/api/v1/auth/logout` - Logout
- **GET** `/api/v1/me` - Get current user info

#### Licenses
- **GET** `/api/v1/licenses` - Get licenses with filters
  - Query params: `status`, `search`, `category`, `license_tier`
  - Automatically filters by user ID from JWT token
- **GET** `/api/v1/licenses/active` - Get active licenses for authenticated user
  - Automatically uses user ID from JWT token
  - No query parameters needed

#### Requests
- **GET** `/api/v1/requests/pending` - Get pending requests for authenticated user
  - Query params: `limit` (default: 2)
  - Returns requests with auto-generated ticket numbers (#REQ-{YEAR}-{ID})
  - Returns total count of pending requests

#### Users
- **GET** `/api/v1/users` - Get list of users (with pagination and filters)

#### AI & Recommendations
- **POST** `/api/v1/ai/recommendations/jd-match` - Generate JD-based software recommendations
- **POST** `/api/v1/ai/consolidation/memo` - Generate consolidation memo
- **POST** `/api/v1/ai/similarity` - Calculate software similarity

#### Health Check
- **GET** `/api/v1/health` - Check API health status

## Environment Variables

### Local Environment
```
base_url: http://localhost:8080
auth_token: (automatically set after login)
company_code: SCB
```

### Production Environment
```
base_url: https://api.assetpulse.scbx.com
auth_token: (automatically set after login)
company_code: (empty, set based on your company)
```

## Demo Accounts

### Local Development
| Username | Password | Role |
|----------|----------|------|
| employee@scb.com | password | Employee/HR |
| manager@scb.com | password | Department Manager |
| cto@scb.com | password | Subsidiary CTO |

## Example Responses

### Get Active Licenses
```json
{
    "msg": "success",
    "data": {
        "licenses": [
            {
                "app_id": 1,
                "app_name": "Slack",
                "app_category": "Communication",
                "app_status": "Active",
                "license_tier": "Pro",
                "assigned_at": "2025-01-15T10:30:00Z",
                "total_seats": 100,
                "reserved_seats": 10,
                "effective_date": "2025-01-01T00:00:00Z",
                "expire_date": "2025-12-31T23:59:59Z"
            }
        ],
        "total": 1
    }
}
```

### Error Response (401 Unauthorized)
```json
{
    "error": "Invalid or expired token"
}
```

### Error Response (500 Internal Server Error)
```json
{
    "msg": "error",
    "error": {
        "code": "INTERNAL_ERROR",
        "message": "internal server error"
    }
}
```

## Notes

- The **Login** request has a test script that automatically saves the JWT token to the environment
- All protected endpoints require the `Authorization: Bearer {{auth_token}}` header
- The **Get Active Licenses** endpoint uses the authenticated user's ID from the JWT token automatically
- Token expiration is set to 24 hours by default

## Troubleshooting

### Token Not Saved
If the token is not automatically saved after login:
1. Check the **Tests** tab in the Login request
2. Ensure the test script is present
3. Manually copy the token from the response and set it in the environment

### 401 Unauthorized
- Check if you're logged in (token exists in environment)
- Token may have expired (login again)
- Ensure you selected the correct environment

### 500 Internal Server Error
- Check if the API server is running
- Check server logs for detailed error messages
- Verify database connection

## Support

For issues or questions, please contact the development team or check the API documentation at `/docs` endpoint (if Swagger is enabled).
