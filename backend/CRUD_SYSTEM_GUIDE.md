# Modern CRUD System Guide

## Overview

The Base App now includes a **modern, flexible CRUD system** that allows you to create multiple CRUD entities (like Portfolio, Visa System, Products, etc.) without changing existing functionality. The system is fully compatible with any project and can be extended easily.

## Key Features

✅ **Pre-built Templates** - 6 ready-to-use templates:
- Portfolio
- Visa Management
- Products
- Blog Posts
- Events
- Contacts

✅ **Schema Validation** - Automatic validation against JSON schemas
✅ **Flexible Schema** - Create custom schemas for any use case
✅ **Full CRUD Operations** - Create, Read, Update, Delete for all entities
✅ **Backward Compatible** - Doesn't affect existing functionality
✅ **Modern Architecture** - Clean, maintainable code structure

## Quick Start

### 1. Get Available Templates

```bash
GET /v1/admin/cruds/templates
Authorization: Bearer <admin_token>
```

### 2. Create Entity from Template (Portfolio Example)

```bash
POST /v1/admin/cruds/templates/portfolio/create
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "display_name": "My Portfolio",
  "description": "Personal portfolio projects"
}
```

### 3. Add Portfolio Item

```bash
POST /v1/admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "E-commerce Website",
  "description": "Full-stack e-commerce platform",
  "category": "web",
  "technologies": ["React", "Node.js", "PostgreSQL"],
  "status": "completed",
  "featured": true,
  "project_url": "https://example.com",
  "github_url": "https://github.com/user/repo"
}
```

### 4. List Portfolio Items

```bash
GET /v1/admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
```

## Creating Custom CRUD Entities

You can also create custom entities without templates:

```bash
POST /v1/admin/cruds/entities
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "entity_name": "custom_entity",
  "display_name": "Custom Entity",
  "description": "Description here",
  "schema": {
    "type": "object",
    "properties": {
      "field1": {
        "type": "string",
        "description": "Field description",
        "required": true
      },
      "field2": {
        "type": "number",
        "description": "Numeric field"
      }
    }
  }
}
```

## Available Templates

### 1. Portfolio (`portfolio`)
Perfect for showcasing projects, skills, and achievements.

**Use Cases:**
- Personal portfolios
- Company project showcases
- Developer portfolios
- Design portfolios

### 2. Visa Management (`visa`)
Manage visa applications and travel documents.

**Use Cases:**
- Travel agencies
- Immigration services
- HR departments managing employee visas
- Personal visa tracking

### 3. Products (`products`)
E-commerce product catalog with inventory management.

**Use Cases:**
- Online stores
- Inventory management
- Product catalogs
- Marketplace platforms

### 4. Blog Posts (`blog_posts`)
Content management for articles and blog posts.

**Use Cases:**
- Blog platforms
- News websites
- Content management
- Documentation sites

### 5. Events (`events`)
Event management and booking system.

**Use Cases:**
- Event planning
- Conference management
- Ticket sales
- Calendar systems

### 6. Contacts (`contacts`)
CRM and contact management system.

**Use Cases:**
- Customer relationship management
- Address books
- Lead management
- Business directories

## API Endpoints Summary

### Template Endpoints
- `GET /v1/admin/cruds/templates` - List all templates
- `GET /v1/admin/cruds/templates/{name}` - Get specific template
- `POST /v1/admin/cruds/templates/{name}/create` - Create entity from template

### Entity Endpoints
- `POST /v1/admin/cruds/entities` - Create custom entity
- `GET /v1/admin/cruds/entities` - List all entities
- `GET /v1/admin/cruds/entities/{id}` - Get entity details
- `PUT /v1/admin/cruds/entities/{id}` - Update entity
- `DELETE /v1/admin/cruds/entities/{id}` - Delete entity

### Data Endpoints
- `POST /v1/admin/cruds/entities/{id}/data` - Create data entry
- `GET /v1/admin/cruds/entities/{id}/data` - List data entries
- `GET /v1/admin/cruds/data/{id}` - Get specific data entry
- `PUT /v1/admin/cruds/data/{id}` - Update data entry
- `DELETE /v1/admin/cruds/data/{id}` - Delete data entry

## Examples

### Example 1: Portfolio System

```bash
# Step 1: Create portfolio entity
curl -X POST http://localhost:8080/v1/admin/cruds/templates/portfolio/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"display_name": "My Portfolio"}'

# Step 2: Add project
curl -X POST http://localhost:8080/v1/admin/cruds/entities/{entity_id}/data \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "E-commerce Platform",
    "description": "Modern e-commerce solution",
    "category": "web",
    "technologies": ["React", "Node.js"],
    "status": "completed"
  }'
```

### Example 2: Visa Management System

```bash
# Step 1: Create visa entity
curl -X POST http://localhost:8080/v1/admin/cruds/templates/visa/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"display_name": "Visa Applications"}'

# Step 2: Add visa application
curl -X POST http://localhost:8080/v1/admin/cruds/entities/{entity_id}/data \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "applicant_name": "John Doe",
    "passport_number": "AB123456",
    "country": "USA",
    "visa_type": "tourist",
    "application_date": "2024-01-15",
    "status": "pending"
  }'
```

## Architecture

The CRUD system follows clean architecture principles:

```
backend/
├── internal/
│   ├── models/
│   │   └── admin_settings.go (CustomCRUDEntity, CustomCRUDData)
│   ├── repositories/
│   │   └── custom_crud_repository.go
│   ├── services/
│   │   ├── custom_crud_service.go (Core CRUD logic)
│   │   └── crud_templates.go (Pre-built templates)
│   └── handlers/
│       ├── admin.go (CRUD entity handlers)
│       └── crud_templates.go (Template handlers)
└── migrations/
    └── 005_admin_settings_and_cruds.up.sql
```

## Benefits

1. **No Code Changes Required** - Create new CRUDs via API
2. **Schema Validation** - Automatic data validation
3. **Flexible** - Support any data structure
4. **Scalable** - Handle unlimited entities
5. **Modern** - JSON-based, RESTful API
6. **Compatible** - Works with existing codebase

## Best Practices

1. **Use Templates** - Start with templates for common use cases
2. **Validate Early** - Test schemas before production use
3. **Document Schemas** - Keep schema documentation updated
4. **Version Control** - Track schema changes
5. **Error Handling** - Implement proper error handling in frontend

## Future Enhancements

- More pre-built templates
- Advanced schema validation
- Field relationships
- Bulk operations
- Import/Export functionality
- Search and filtering
- Pagination improvements

## Support

For detailed API documentation, see:
- `backend/docs/CRUD_TEMPLATES.md` - Complete template documentation
- `backend/API_ENDPOINTS.md` - All API endpoints

