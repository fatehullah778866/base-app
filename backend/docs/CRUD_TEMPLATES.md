# CRUD Templates Documentation

## Overview

The Base App includes pre-built CRUD templates that allow you to quickly create common entity types like portfolios, visa management systems, products, blog posts, events, and contacts.

## Available Templates

### 1. Portfolio (`portfolio`)
Manage portfolio items with projects, skills, and achievements.

**Fields:**
- `title` (string, required) - Project title
- `description` (string, required) - Project description
- `category` (string) - Project category: web, mobile, desktop, other
- `technologies` (array) - Technologies used
- `image_url` (string) - Project image URL
- `project_url` (string) - Project URL
- `github_url` (string) - GitHub repository URL
- `start_date` (date) - Project start date
- `end_date` (date) - Project end date
- `status` (string) - Status: completed, in-progress, planned
- `featured` (boolean) - Featured project

**Example Usage:**
```bash
# Get all templates
GET /v1/admin/cruds/templates

# Get portfolio template
GET /v1/admin/cruds/templates/portfolio

# Create portfolio entity from template
POST /v1/admin/cruds/templates/portfolio/create
{
  "display_name": "My Portfolio",
  "description": "Personal portfolio projects"
}
```

### 2. Visa Management (`visa`)
Manage visa applications and documents.

**Fields:**
- `applicant_name` (string, required) - Full name of applicant
- `passport_number` (string, required) - Passport number
- `country` (string, required) - Destination country
- `visa_type` (string, required) - Type: tourist, business, student, work, transit, other
- `application_date` (date, required) - Application submission date
- `status` (string) - Status: pending, under-review, approved, rejected, cancelled
- `expiry_date` (date) - Visa expiry date
- `documents` (array) - Attached documents
- `notes` (string) - Additional notes

**Example Usage:**
```bash
# Create visa management entity
POST /v1/admin/cruds/templates/visa/create
{
  "display_name": "Visa Applications",
  "description": "Manage all visa applications"
}
```

### 3. Products (`products`)
Manage product catalog with inventory.

**Fields:**
- `name` (string, required) - Product name
- `sku` (string, required) - SKU code
- `description` (string) - Product description
- `price` (number, required) - Product price
- `currency` (string) - Currency code (default: USD)
- `category` (string) - Product category
- `stock_quantity` (integer) - Stock quantity
- `images` (array) - Product images URLs
- `tags` (array) - Product tags
- `active` (boolean) - Product active status

### 4. Blog Posts (`blog_posts`)
Manage blog posts and articles.

**Fields:**
- `title` (string, required) - Post title
- `slug` (string, required) - URL slug
- `content` (string, required) - Post content (HTML or Markdown)
- `excerpt` (string) - Post excerpt
- `author` (string) - Author name
- `category` (string) - Post category
- `tags` (array) - Post tags
- `featured_image` (string) - Featured image URL
- `published` (boolean) - Published status
- `published_at` (date-time) - Publication date

### 5. Events (`events`)
Manage events and bookings.

**Fields:**
- `title` (string, required) - Event title
- `description` (string) - Event description
- `start_date` (date-time, required) - Event start date and time
- `end_date` (date-time, required) - Event end date and time
- `location` (string) - Event location
- `venue` (string) - Venue name
- `capacity` (integer) - Maximum capacity
- `price` (number) - Ticket price
- `status` (string) - Status: draft, published, cancelled, completed
- `image_url` (string) - Event image URL

### 6. Contacts (`contacts`)
Manage contact list and CRM.

**Fields:**
- `first_name` (string, required) - First name
- `last_name` (string, required) - Last name
- `email` (string, required) - Email address
- `phone` (string) - Phone number
- `company` (string) - Company name
- `position` (string) - Job position
- `address` (string) - Address
- `city` (string) - City
- `country` (string) - Country
- `tags` (array) - Contact tags
- `notes` (string) - Additional notes

## API Endpoints

### Get All Templates
```http
GET /v1/admin/cruds/templates
Authorization: Bearer <admin_token>
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "portfolio",
      "display_name": "Portfolio",
      "description": "Manage portfolio items...",
      "icon": "briefcase",
      "category": "business",
      "schema": { ... }
    },
    ...
  ]
}
```

### Get Specific Template
```http
GET /v1/admin/cruds/templates/{name}
Authorization: Bearer <admin_token>
```

### Create Entity from Template
```http
POST /v1/admin/cruds/templates/{name}/create
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "display_name": "Custom Display Name",  // Optional
  "description": "Custom Description"     // Optional
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "entity_name": "portfolio",
    "display_name": "Custom Display Name",
    "schema": { ... },
    ...
  },
  "message": "Entity created from template successfully"
}
```

## Working with CRUD Data

After creating an entity from a template, you can manage its data:

### Create Data
```http
POST /v1/admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "My Project",
  "description": "Project description",
  "category": "web",
  "technologies": ["React", "Node.js"],
  "status": "completed"
}
```

### List Data
```http
GET /v1/admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
```

### Update Data
```http
PUT /v1/admin/cruds/data/{data_id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "title": "Updated Title",
  "status": "in-progress"
}
```

### Delete Data
```http
DELETE /v1/admin/cruds/data/{data_id}
Authorization: Bearer <admin_token>
```

## Creating Custom Templates

You can also create custom CRUD entities without using templates:

```http
POST /v1/admin/cruds/entities
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "entity_name": "custom_entity",
  "display_name": "Custom Entity",
  "description": "Custom entity description",
  "schema": {
    "type": "object",
    "properties": {
      "field1": {
        "type": "string",
        "description": "Field description",
        "required": true
      }
    }
  }
}
```

## Best Practices

1. **Use Templates**: Start with templates for common use cases
2. **Customize Schema**: Modify schema after creation if needed
3. **Validate Data**: Always validate data against schema before saving
4. **Version Control**: Keep track of schema changes
5. **Documentation**: Document custom fields and their purposes

## Examples

### Example 1: Create Portfolio System
```bash
# 1. Create portfolio entity
curl -X POST http://localhost:8080/v1/admin/cruds/templates/portfolio/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"display_name": "My Portfolio"}'

# 2. Add portfolio item
curl -X POST http://localhost:8080/v1/admin/cruds/entities/{entity_id}/data \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "E-commerce Website",
    "description": "Full-stack e-commerce platform",
    "category": "web",
    "technologies": ["React", "Node.js", "PostgreSQL"],
    "status": "completed",
    "featured": true
  }'
```

### Example 2: Create Visa Management System
```bash
# 1. Create visa entity
curl -X POST http://localhost:8080/v1/admin/cruds/templates/visa/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"display_name": "Visa Applications"}'

# 2. Add visa application
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

## Schema Validation

The system validates data against the schema before saving. Required fields must be present, and field types must match the schema definition.

## Future Enhancements

- More pre-built templates
- Schema versioning
- Advanced validation rules
- Field relationships
- Bulk operations
- Import/Export functionality

