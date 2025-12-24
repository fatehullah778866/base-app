# Comprehensive Search System

## Overview

The backend now includes a comprehensive search system that supports:
- **Global Search** - Search across all entities
- **Location-based Search** - Search by country, city, or location
- **Message Search** - Search in messages and conversations
- **Dashboard Search** - Search dashboard items with filters
- **User Search** - Search users by name, email, or location
- **Notification Search** - Search notifications
- **Custom CRUD Search** - Search in custom CRUD entities
- **Search History** - Track and manage search history

## API Endpoints

### Search
```http
GET /v1/search?q=query&type=all&limit=20&offset=0
POST /v1/search
```

**Query Parameters:**
- `q` - Search query (required for text search)
- `type` - Search type: `all`, `users`, `dashboard_items`, `messages`, `notifications`, `cruds`, `locations`
- `limit` - Results limit (default: 20)
- `offset` - Pagination offset (default: 0)
- `location` - Location filter
- `country` - Country filter
- `city` - City filter
- `date_from` - Start date (RFC3339 format)
- `date_to` - End date (RFC3339 format)
- `category` - Category filter
- `status` - Status filter
- `entity_id` - Custom CRUD entity ID

**JSON Body (POST):**
```json
{
  "query": "search term",
  "type": "all",
  "limit": 20,
  "offset": 0,
  "location": "New York",
  "country": "USA",
  "city": "New York",
  "date_from": "2024-01-01T00:00:00Z",
  "date_to": "2024-12-31T23:59:59Z",
  "category": "work",
  "status": "active",
  "entity_id": "uuid"
}
```

### Search History
```http
GET /v1/search/history?limit=50
DELETE /v1/search/history
```

## Search Types

### 1. Global Search (`type=all`)
Searches across all entities:
- Dashboard items
- Messages
- Users
- Notifications
- Custom CRUDs
- Location-based results

### 2. Dashboard Search (`type=dashboard_items`)
- Full-text search in titles and descriptions
- Filter by category
- Filter by status
- Filter by date range
- Search in metadata

### 3. Message Search (`type=messages`)
- Full-text search in message content and subjects
- Filter by date range
- Filter by read/unread status
- Search across all user's conversations

### 4. User Search (`type=users`)
- Search by name or email
- Location-based search (country, city)
- Excludes current user
- Only active/pending users

### 5. Notification Search (`type=notifications`)
- Search in notification titles and messages
- Filter by notification type
- User-specific results

### 6. Custom CRUD Search (`type=cruds`)
- Search in custom CRUD data
- Filter by entity ID
- JSON data search
- Schema-aware search

### 7. Location Search (`type=locations`)
- Search users by location
- Search dashboard items with location metadata
- Country and city filters

## Advanced Features

### Full-Text Search (FTS5)
- Uses SQLite FTS5 for fast full-text search
- Falls back to LIKE search if FTS5 unavailable
- Supports multiple search terms

### Search History
- Automatically saves search queries
- Tracks search type and result count
- Can retrieve and clear history

### Filtering
- Date range filtering
- Category filtering
- Status filtering
- Location filtering
- Pagination support

## Response Format

```json
{
  "success": true,
  "data": {
    "type": "search_results",
    "id": "uuid",
    "title": "Search Results",
    "data": {
      "results": [
        {
          "type": "dashboard_item",
          "id": "uuid",
          "title": "Item Title",
          "description": "Item description",
          "data": { ... }
        }
      ],
      "count": 10,
      "query": "search term",
      "type": "all",
      "limit": 20,
      "offset": 0
    }
  }
}
```

## Examples

### Basic Search
```bash
curl -X GET "http://localhost:8080/v1/search?q=project&type=dashboard_items" \
  -H "Authorization: Bearer <token>"
```

### Advanced Search with Filters
```bash
curl -X POST "http://localhost:8080/v1/search" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "meeting",
    "type": "messages",
    "date_from": "2024-01-01T00:00:00Z",
    "date_to": "2024-12-31T23:59:59Z",
    "limit": 50
  }'
```

### Location-based Search
```bash
curl -X GET "http://localhost:8080/v1/search?type=locations&country=USA&city=New York" \
  -H "Authorization: Bearer <token>"
```

### Search History
```bash
curl -X GET "http://localhost:8080/v1/search/history?limit=20" \
  -H "Authorization: Bearer <token>"
```

## Implementation Details

### Backend Components

1. **SearchService** (`internal/services/search_service.go`)
   - Main search logic
   - Coordinates multiple repository calls
   - Applies filters and pagination
   - Manages search history

2. **SearchRepository** (`internal/repositories/search_repository_impl.go`)
   - Database queries
   - FTS5 full-text search
   - Location-based queries
   - Search history management

3. **SearchHandler** (`internal/handlers/search.go`)
   - HTTP request handling
   - Parameter parsing (query string and JSON)
   - Response formatting

### Database

- Uses SQLite FTS5 for full-text search
- Search history table tracks user searches
- Location data from sessions table
- Metadata stored as JSON for flexibility

## Performance

- FTS5 indexes for fast text search
- Pagination to limit result sets
- Efficient queries with proper indexing
- Fallback to LIKE search if FTS5 unavailable

## Future Enhancements

- Search result ranking
- Search suggestions/autocomplete
- Search analytics
- Advanced location search with geolocation
- Search result caching
- Multi-language search support

