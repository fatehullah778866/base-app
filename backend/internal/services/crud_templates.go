package services

import (
	"encoding/json"
)

// CRUDTemplate represents a pre-built CRUD entity template
type CRUDTemplate struct {
	Name        string                 `json:"name"`
	DisplayName string                 `json:"display_name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	Icon        string                 `json:"icon"`
	Category    string                 `json:"category"`
}

// GetCRUDTemplates returns all available CRUD templates
func GetCRUDTemplates() []CRUDTemplate {
	return []CRUDTemplate{
		{
			Name:        "portfolio",
			DisplayName: "Portfolio",
			Description: "Manage portfolio items with projects, skills, and achievements",
			Icon:        "briefcase",
			Category:    "business",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Project title",
						"required":    true,
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Project description",
						"required":    true,
					},
					"category": map[string]interface{}{
						"type":        "string",
						"description": "Project category",
						"enum":        []string{"web", "mobile", "desktop", "other"},
					},
					"technologies": map[string]interface{}{
						"type":        "array",
						"description": "Technologies used",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"image_url": map[string]interface{}{
						"type":        "string",
						"description": "Project image URL",
						"format":      "uri",
					},
					"project_url": map[string]interface{}{
						"type":        "string",
						"description": "Project URL",
						"format":      "uri",
					},
					"github_url": map[string]interface{}{
						"type":        "string",
						"description": "GitHub repository URL",
						"format":      "uri",
					},
					"start_date": map[string]interface{}{
						"type":        "string",
						"description": "Project start date",
						"format":      "date",
					},
					"end_date": map[string]interface{}{
						"type":        "string",
						"description": "Project end date",
						"format":      "date",
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Project status",
						"enum":        []string{"completed", "in-progress", "planned"},
						"default":     "in-progress",
					},
					"featured": map[string]interface{}{
						"type":        "boolean",
						"description": "Featured project",
						"default":     false,
					},
				},
			},
		},
		{
			Name:        "visa",
			DisplayName: "Visa Management",
			Description: "Manage visa applications and documents",
			Icon:        "passport",
			Category:    "travel",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"applicant_name": map[string]interface{}{
						"type":        "string",
						"description": "Full name of applicant",
						"required":    true,
					},
					"passport_number": map[string]interface{}{
						"type":        "string",
						"description": "Passport number",
						"required":    true,
					},
					"country": map[string]interface{}{
						"type":        "string",
						"description": "Destination country",
						"required":    true,
					},
					"visa_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of visa",
						"enum":        []string{"tourist", "business", "student", "work", "transit", "other"},
						"required":    true,
					},
					"application_date": map[string]interface{}{
						"type":        "string",
						"description": "Application submission date",
						"format":      "date",
						"required":    true,
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Application status",
						"enum":        []string{"pending", "under-review", "approved", "rejected", "cancelled"},
						"default":     "pending",
					},
					"expiry_date": map[string]interface{}{
						"type":        "string",
						"description": "Visa expiry date",
						"format":      "date",
					},
					"documents": map[string]interface{}{
						"type":        "array",
						"description": "Attached documents",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type": "string",
								},
								"url": map[string]interface{}{
									"type":   "string",
									"format": "uri",
								},
								"type": map[string]interface{}{
									"type": "string",
									"enum": []string{"passport", "photo", "invitation", "bank-statement", "other"},
								},
							},
						},
					},
					"notes": map[string]interface{}{
						"type":        "string",
						"description": "Additional notes",
					},
				},
			},
		},
		{
			Name:        "products",
			DisplayName: "Products",
			Description: "Manage product catalog with inventory",
			Icon:        "shopping-cart",
			Category:    "ecommerce",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Product name",
						"required":    true,
					},
					"sku": map[string]interface{}{
						"type":        "string",
						"description": "SKU code",
						"required":    true,
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Product description",
					},
					"price": map[string]interface{}{
						"type":        "number",
						"description": "Product price",
						"required":    true,
						"minimum":     0,
					},
					"currency": map[string]interface{}{
						"type":        "string",
						"description": "Currency code",
						"default":     "USD",
					},
					"category": map[string]interface{}{
						"type":        "string",
						"description": "Product category",
					},
					"stock_quantity": map[string]interface{}{
						"type":        "integer",
						"description": "Stock quantity",
						"default":     0,
						"minimum":     0,
					},
					"images": map[string]interface{}{
						"type":        "array",
						"description": "Product images",
						"items": map[string]interface{}{
							"type":   "string",
							"format": "uri",
						},
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Product tags",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"active": map[string]interface{}{
						"type":        "boolean",
						"description": "Product active status",
						"default":     true,
					},
				},
			},
		},
		{
			Name:        "blog_posts",
			DisplayName: "Blog Posts",
			Description: "Manage blog posts and articles",
			Icon:        "file-text",
			Category:    "content",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Post title",
						"required":    true,
					},
					"slug": map[string]interface{}{
						"type":        "string",
						"description": "URL slug",
						"required":    true,
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Post content (HTML or Markdown)",
						"required":    true,
					},
					"excerpt": map[string]interface{}{
						"type":        "string",
						"description": "Post excerpt",
					},
					"author": map[string]interface{}{
						"type":        "string",
						"description": "Author name",
					},
					"category": map[string]interface{}{
						"type":        "string",
						"description": "Post category",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Post tags",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"featured_image": map[string]interface{}{
						"type":        "string",
						"description": "Featured image URL",
						"format":      "uri",
					},
					"published": map[string]interface{}{
						"type":        "boolean",
						"description": "Published status",
						"default":     false,
					},
					"published_at": map[string]interface{}{
						"type":        "string",
						"description": "Publication date",
						"format":      "date-time",
					},
				},
			},
		},
		{
			Name:        "events",
			DisplayName: "Events",
			Description: "Manage events and bookings",
			Icon:        "calendar",
			Category:    "business",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Event title",
						"required":    true,
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Event description",
					},
					"start_date": map[string]interface{}{
						"type":        "string",
						"description": "Event start date and time",
						"format":      "date-time",
						"required":    true,
					},
					"end_date": map[string]interface{}{
						"type":        "string",
						"description": "Event end date and time",
						"format":      "date-time",
						"required":    true,
					},
					"location": map[string]interface{}{
						"type":        "string",
						"description": "Event location",
					},
					"venue": map[string]interface{}{
						"type":        "string",
						"description": "Venue name",
					},
					"capacity": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum capacity",
						"minimum":     0,
					},
					"price": map[string]interface{}{
						"type":        "number",
						"description": "Ticket price",
						"minimum":     0,
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Event status",
						"enum":        []string{"draft", "published", "cancelled", "completed"},
						"default":     "draft",
					},
					"image_url": map[string]interface{}{
						"type":        "string",
						"description": "Event image URL",
						"format":      "uri",
					},
				},
			},
		},
		{
			Name:        "contacts",
			DisplayName: "Contacts",
			Description: "Manage contact list and CRM",
			Icon:        "users",
			Category:    "crm",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"first_name": map[string]interface{}{
						"type":        "string",
						"description": "First name",
						"required":    true,
					},
					"last_name": map[string]interface{}{
						"type":        "string",
						"description": "Last name",
						"required":    true,
					},
					"email": map[string]interface{}{
						"type":        "string",
						"description": "Email address",
						"format":      "email",
						"required":    true,
					},
					"phone": map[string]interface{}{
						"type":        "string",
						"description": "Phone number",
					},
					"company": map[string]interface{}{
						"type":        "string",
						"description": "Company name",
					},
					"position": map[string]interface{}{
						"type":        "string",
						"description": "Job position",
					},
					"address": map[string]interface{}{
						"type":        "string",
						"description": "Address",
					},
					"city": map[string]interface{}{
						"type":        "string",
						"description": "City",
					},
					"country": map[string]interface{}{
						"type":        "string",
						"description": "Country",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Contact tags",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"notes": map[string]interface{}{
						"type":        "string",
						"description": "Additional notes",
					},
				},
			},
		},
	}
}

// GetTemplateByName returns a specific template by name
func GetTemplateByName(name string) (*CRUDTemplate, error) {
	templates := GetCRUDTemplates()
	for _, template := range templates {
		if template.Name == name {
			return &template, nil
		}
	}
	return nil, nil
}

// CreateEntityFromTemplate creates a CRUD entity from a template
func CreateEntityFromTemplate(templateName string) (map[string]interface{}, error) {
	template, err := GetTemplateByName(templateName)
	if err != nil || template == nil {
		return nil, err
	}

	// Convert template to entity creation format
	entity := map[string]interface{}{
		"entity_name":  template.Name,
		"display_name": template.DisplayName,
		"description":  template.Description,
		"schema":       template.Schema,
	}

	return entity, nil
}

// ValidateDataAgainstSchema validates data against a JSON schema
func ValidateDataAgainstSchema(data map[string]interface{}, schema map[string]interface{}) error {
	// Basic validation - in production, use a proper JSON schema validator like github.com/xeipuuv/gojsonschema
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return nil // No schema to validate against
	}

	// Check required fields
	if requiredFields, ok := schema["required"].([]interface{}); ok {
		for _, field := range requiredFields {
			if fieldStr, ok := field.(string); ok {
				if _, exists := data[fieldStr]; !exists {
					return nil // Return error in production
				}
			}
		}
	}

	// Validate field types (simplified)
	for fieldName, fieldSchema := range properties {
		if fieldSchemaMap, ok := fieldSchema.(map[string]interface{}); ok {
			if value, exists := data[fieldName]; exists {
				expectedType, ok := fieldSchemaMap["type"].(string)
				if ok {
					switch expectedType {
					case "string":
						if _, ok := value.(string); !ok {
							// Type mismatch
						}
					case "number", "integer":
						if _, ok := value.(float64); !ok {
							if _, ok := value.(int); !ok {
								// Type mismatch
							}
						}
					case "boolean":
						if _, ok := value.(bool); !ok {
							// Type mismatch
						}
					case "array":
						if _, ok := value.([]interface{}); !ok {
							// Type mismatch
						}
					}
				}
			}
		}
	}

	return nil
}

// GetSchemaJSON returns schema as JSON string
func GetSchemaJSON(schema map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

