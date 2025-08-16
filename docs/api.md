# Alarm Agent Web API Documentation

## Overview
The Alarm Agent Web API provides REST endpoints for managing calendar events through a web interface, offering the same functionality as the WhatsApp LLM bot.

## Authentication
The API uses WhatsApp number-based authentication. Include the WhatsApp number in the `X-WA-Number` header for all authenticated endpoints.

```
X-WA-Number: +5511999999999
```

## Base URL
```
http://localhost:8080/api/v1
```

## Endpoints

### Authentication

#### Authenticate User
Verify user exists and get profile information.

```http
POST /api/v1/auth
```

**Request Body:**
```json
{
  "wa_number": "+5511999999999"
}
```

**Response:**
```json
{
  "message": "User authenticated successfully",
  "data": {
    "id": 1,
    "wa_number": "+5511999999999",
    "name": "João Silva",
    "timezone": "America/Sao_Paulo",
    "default_remind_before_minutes": 30,
    "default_remind_frequency_minutes": 15,
    "default_require_confirmation": true,
    "llm_provider": "anthropic",
    "llm_model": "claude-3-haiku-20240307",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### User Profile

#### Get Profile
Get current user's profile information.

```http
GET /api/v1/profile
Headers: X-WA-Number: +5511999999999
```

**Response:**
```json
{
  "id": 1,
  "wa_number": "+5511999999999",
  "name": "João Silva",
  "timezone": "America/Sao_Paulo",
  "default_remind_before_minutes": 30,
  "default_remind_frequency_minutes": 15,
  "default_require_confirmation": true,
  "llm_provider": "anthropic",
  "llm_model": "claude-3-haiku-20240307",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

#### Update Profile
Update user's profile settings.

```http
PUT /api/v1/profile
Headers: X-WA-Number: +5511999999999
```

**Request Body:**
```json
{
  "name": "João Silva Santos",
  "timezone": "America/Sao_Paulo",
  "default_remind_before_minutes": 45,
  "default_remind_frequency_minutes": 10,
  "default_require_confirmation": false,
  "llm_provider": "openai",
  "llm_model": "gpt-4"
}
```

### Events Management

#### Create Event
Create a new calendar event.

```http
POST /api/v1/events
Headers: X-WA-Number: +5511999999999
```

**Request Body:**
```json
{
  "title": "Team Meeting",
  "location": "Conference Room A",
  "starts_at": "2024-01-15T14:30:00Z",
  "remind_before_minutes": 30,
  "remind_frequency_minutes": 15,
  "require_confirmation": true,
  "max_notifications": 3
}
```

**Response:**
```json
{
  "message": "Event created successfully",
  "data": {
    "id": 123,
    "title": "Team Meeting",
    "location": "Conference Room A",
    "starts_at": "2024-01-15T14:30:00Z",
    "remind_before_minutes": 30,
    "remind_frequency_minutes": 15,
    "require_confirmation": true,
    "max_notifications": 3,
    "status": "scheduled",
    "notifications_sent": 0,
    "last_notified_at": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### List Events
Get all events for the authenticated user.

```http
GET /api/v1/events?start_date=2024-01-01&end_date=2024-01-31&status=scheduled&limit=50&offset=0
Headers: X-WA-Number: +5511999999999
```

**Query Parameters:**
- `start_date` (optional): Filter events starting from this date (ISO 8601)
- `end_date` (optional): Filter events up to this date (ISO 8601)
- `status` (optional): Filter by status (scheduled, confirmed, canceled, completed)
- `limit` (optional): Maximum number of events to return (default: 50, max: 100)
- `offset` (optional): Number of events to skip for pagination (default: 0)

**Response:**
```json
{
  "events": [
    {
      "id": 123,
      "title": "Team Meeting",
      "location": "Conference Room A",
      "starts_at": "2024-01-15T14:30:00Z",
      "remind_before_minutes": 30,
      "remind_frequency_minutes": 15,
      "require_confirmation": true,
      "max_notifications": 3,
      "status": "scheduled",
      "notifications_sent": 0,
      "last_notified_at": null,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "total_count": 1,
  "limit": 50,
  "offset": 0
}
```

#### Get Event
Get a specific event by ID.

```http
GET /api/v1/events/123
Headers: X-WA-Number: +5511999999999
```

#### Update Event
Update an existing event.

```http
PUT /api/v1/events/123
Headers: X-WA-Number: +5511999999999
```

**Request Body:**
```json
{
  "title": "Updated Team Meeting",
  "starts_at": "2024-01-15T15:00:00Z",
  "status": "confirmed"
}
```

#### Delete Event
Cancel/delete an event.

```http
DELETE /api/v1/events/123
Headers: X-WA-Number: +5511999999999
```

#### Confirm Event
Confirm an event (change status to confirmed).

```http
POST /api/v1/events/123/confirm
Headers: X-WA-Number: +5511999999999
```

### LLM Configuration

#### Get Available Providers
Get all available LLM providers and their models.

```http
GET /api/v1/llm/providers
```

**Response:**
```json
{
  "message": "Providers retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "anthropic",
      "display_name": "Anthropic",
      "description": "Anthropic Claude models",
      "is_active": true,
      "models": [
        {
          "id": 1,
          "name": "claude-3-haiku-20240307",
          "display_name": "Claude 3 Haiku",
          "description": "Fast and efficient Claude model",
          "is_default": true,
          "is_active": true
        },
        {
          "id": 2,
          "name": "claude-3-sonnet-20240229",
          "display_name": "Claude 3 Sonnet", 
          "description": "Balanced Claude model",
          "is_default": false,
          "is_active": true
        }
      ]
    }
  ]
}
```

#### Get Models for Provider
Get models for a specific provider.

```http
GET /api/v1/llm/providers/anthropic/models
```

#### Get Default Model
Get the system default LLM model.

```http
GET /api/v1/llm/default
```

## Error Responses

All endpoints return errors in this format:

```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "code": 400
}
```

### Common Error Codes
- `unauthorized`: Missing or invalid authentication
- `number_not_whitelisted`: WhatsApp number not authorized
- `user_not_found`: User doesn't exist
- `invalid_request`: Request validation failed
- `event_not_found`: Event doesn't exist or access denied
- `create_failed`: Failed to create resource
- `update_failed`: Failed to update resource

## Status Codes
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Example Usage

### Creating a Complete Event Workflow

1. **Authenticate User**
```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{"wa_number": "+5511999999999"}'
```

2. **Create Event**
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "X-WA-Number: +5511999999999" \
  -d '{
    "title": "Doctor Appointment",
    "location": "Medical Center",
    "starts_at": "2024-01-15T10:00:00Z",
    "remind_before_minutes": 60
  }'
```

3. **List Events**
```bash
curl -X GET "http://localhost:8080/api/v1/events?limit=10" \
  -H "X-WA-Number: +5511999999999"
```

4. **Update Event**
```bash
curl -X PUT http://localhost:8080/api/v1/events/123 \
  -H "Content-Type: application/json" \
  -H "X-WA-Number: +5511999999999" \
  -d '{"starts_at": "2024-01-15T10:30:00Z"}'
```

5. **Confirm Event**
```bash
curl -X POST http://localhost:8080/api/v1/events/123/confirm \
  -H "X-WA-Number: +5511999999999"
```