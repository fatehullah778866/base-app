package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type SearchService struct {
	searchRepo        repositories.SearchRepository
	dashboardRepo     repositories.DashboardRepository
	messageRepo       repositories.MessageRepository
	notificationRepo  repositories.NotificationRepository
	customCRUDRepo    repositories.CustomCRUDRepository
	userRepo          repositories.UserRepository
	logger            *zap.Logger
}

func NewSearchService(
	searchRepo repositories.SearchRepository,
	dashboardRepo repositories.DashboardRepository,
	messageRepo repositories.MessageRepository,
	notificationRepo repositories.NotificationRepository,
	customCRUDRepo repositories.CustomCRUDRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) *SearchService {
	return &SearchService{
		searchRepo:       searchRepo,
		dashboardRepo:    dashboardRepo,
		messageRepo:      messageRepo,
		notificationRepo: notificationRepo,
		customCRUDRepo:   customCRUDRepo,
		userRepo:         userRepo,
		logger:           logger,
	}
}

// SearchRequest represents advanced search parameters
type SearchRequest struct {
	Query       string    `json:"query"`
	Type        string    `json:"type"`        // all, users, dashboard_items, messages, notifications, cruds, locations
	Limit       int       `json:"limit"`
	Offset      int       `json:"offset"`
	Location    *string   `json:"location"`    // City or country
	Country     *string   `json:"country"`
	City        *string   `json:"city"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	Category    *string   `json:"category"`
	Status      *string   `json:"status"`
	EntityID    *uuid.UUID `json:"entity_id"`  // For custom CRUD search
}

// Search performs a comprehensive search across all searchable entities
func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, req SearchRequest) (*models.SearchResult, error) {
	if req.Query == "" && req.Location == nil && req.Country == nil && req.City == nil {
		return &models.SearchResult{
			Type: "search_results",
			ID:   uuid.New().String(),
			Title: "Search Results",
			Data: map[string]interface{}{
				"results": []interface{}{},
				"count":   0,
				"query":   req.Query,
				"type":    req.Type,
			},
		}, nil
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Type == "" {
		req.Type = "all"
	}

	var results []models.SearchResult
	var totalCount int

	// Global search - search across all entities
	if req.Type == "all" || req.Type == "dashboard_items" {
		items, err := s.searchDashboardItems(ctx, userID, req)
		if err == nil {
			for _, item := range items {
				results = append(results, models.SearchResult{
					Type:        "dashboard_item",
					ID:          item.ID.String(),
					Title:       item.Title,
					Description: item.Description,
					Data:        item,
				})
			}
			totalCount += len(items)
		}
	}

	if req.Type == "all" || req.Type == "messages" {
		messages, err := s.searchMessages(ctx, userID, req)
		if err == nil {
			for _, msg := range messages {
				title := "Message"
				if msg.Subject != nil {
					title = *msg.Subject
				}
				results = append(results, models.SearchResult{
					Type:        "message",
					ID:          msg.ID.String(),
					Title:       title,
					Description: &msg.Content,
					Data:        msg,
				})
			}
			totalCount += len(messages)
		}
	}

	if req.Type == "all" || req.Type == "users" {
		users, err := s.searchUsers(ctx, userID, req)
		if err == nil {
			for _, user := range users {
				results = append(results, models.SearchResult{
					Type:  "user",
					ID:    user.ID.String(),
					Title: user.Name,
					Data:  user,
				})
			}
			totalCount += len(users)
		}
	}

	if req.Type == "all" || req.Type == "notifications" {
		notifications, err := s.searchNotifications(ctx, userID, req)
		if err == nil {
			for _, notif := range notifications {
				msg := notif.Message
				results = append(results, models.SearchResult{
					Type:        "notification",
					ID:          notif.ID.String(),
					Title:       notif.Title,
					Description: &msg,
					Data:        notif,
				})
			}
			totalCount += len(notifications)
		}
	}

	if req.Type == "all" || req.Type == "cruds" || req.Type == "custom_cruds" {
		crudData, err := s.searchCustomCRUDs(ctx, userID, req)
		if err == nil {
			for _, data := range crudData {
				var dataMap map[string]interface{}
				if err := json.Unmarshal([]byte(data.Data), &dataMap); err == nil {
					title := "CRUD Item"
					if titleVal, ok := dataMap["title"].(string); ok {
						title = titleVal
					} else if nameVal, ok := dataMap["name"].(string); ok {
						title = nameVal
					}
					results = append(results, models.SearchResult{
						Type:  "crud_data",
						ID:    data.ID.String(),
						Title: title,
						Data:  dataMap,
					})
				}
			}
			totalCount += len(crudData)
		}
	}

	if req.Type == "all" || req.Type == "locations" {
		locationResults, err := s.searchByLocation(ctx, userID, req)
		if err == nil {
			results = append(results, locationResults...)
			totalCount += len(locationResults)
		}
	}

	// Save search history
	history := &models.SearchHistory{
		ID:           uuid.New(),
		UserID:       userID,
		Query:        req.Query,
		SearchType:   &req.Type,
		ResultsCount: totalCount,
		CreatedAt:    time.Now(),
	}
	_ = s.searchRepo.SaveSearchHistory(ctx, history)

	return &models.SearchResult{
		Type: "search_results",
		ID:   uuid.New().String(),
		Title: "Search Results",
		Data: map[string]interface{}{
			"results": results,
			"count":   totalCount,
			"query":   req.Query,
			"type":    req.Type,
			"limit":   req.Limit,
			"offset":  req.Offset,
		},
	}, nil
}

// searchDashboardItems searches dashboard items with advanced filters
func (s *SearchService) searchDashboardItems(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.DashboardItem, error) {
	if req.Query != "" {
		return s.searchRepo.SearchDashboardItems(ctx, userID, req.Query, req.Limit)
	}
	
	// Search by filters (category, status, date range)
	items, err := s.dashboardRepo.ListItems(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var filtered []*models.DashboardItem
	for _, item := range items {
		// Filter by category
		if req.Category != nil && (item.Category == nil || *item.Category != *req.Category) {
			continue
		}
		// Filter by status
		if req.Status != nil && item.Status != *req.Status {
			continue
		}
		// Filter by date range
		if req.DateFrom != nil && item.CreatedAt.Before(*req.DateFrom) {
			continue
		}
		if req.DateTo != nil && item.CreatedAt.After(*req.DateTo) {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered, nil
}

// searchMessages searches messages with advanced filters
func (s *SearchService) searchMessages(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.Message, error) {
	if req.Query != "" {
		return s.searchRepo.SearchMessages(ctx, userID, req.Query, req.Limit)
	}
	
	// Get all messages and filter
	conversations, err := s.messageRepo.GetConversations(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMessages []*models.Message
	for _, conv := range conversations {
		messages, err := s.messageRepo.GetMessages(ctx, conv.ID, userID, 100, 0)
		if err == nil {
			allMessages = append(allMessages, messages...)
		}
	}

	// Apply filters
	var filtered []*models.Message
	for _, msg := range allMessages {
		// Filter by date range
		if req.DateFrom != nil && msg.CreatedAt.Before(*req.DateFrom) {
			continue
		}
		if req.DateTo != nil && msg.CreatedAt.After(*req.DateTo) {
			continue
		}
		// Filter by read status
		if req.Status != nil {
			if *req.Status == "read" && !msg.IsRead {
				continue
			}
			if *req.Status == "unread" && msg.IsRead {
				continue
			}
		}
		filtered = append(filtered, msg)
		if len(filtered) >= req.Limit {
			break
		}
	}

	return filtered, nil
}

// searchUsers searches users with location filters
func (s *SearchService) searchUsers(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.User, error) {
	if req.Query != "" {
		return s.searchRepo.SearchUsers(ctx, userID, req.Query, req.Limit)
	}
	
	// Location-based search
	if req.Location != nil || req.Country != nil || req.City != nil {
		return s.searchUsersByLocation(ctx, userID, req)
	}

	// For now, return empty - user search by location would need location data
	// In production, you'd query a location table or user settings
	return []*models.User{}, nil
}

// searchUsersByLocation searches users by location
func (s *SearchService) searchUsersByLocation(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.User, error) {
	// Use repository method for location-based search
	return s.searchRepo.SearchUsersByLocation(ctx, userID, req.Country, req.City, req.Limit)
}

// searchNotifications searches notifications
func (s *SearchService) searchNotifications(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.Notification, error) {
	if req.Query != "" {
		return s.searchRepo.SearchNotifications(ctx, userID, req.Query, req.Limit)
	}
	
	notifications, err := s.notificationRepo.GetByUserID(ctx, userID, false, req.Limit)
	if err != nil {
		return nil, err
	}

	if req.Query == "" {
		return notifications, nil
	}

	// Filter by query
	var filtered []*models.Notification
	queryLower := strings.ToLower(req.Query)
	for _, notif := range notifications {
		if strings.Contains(strings.ToLower(notif.Title), queryLower) ||
			strings.Contains(strings.ToLower(notif.Message), queryLower) {
			filtered = append(filtered, notif)
		}
	}

	return filtered, nil
}

// searchCustomCRUDs searches custom CRUD data
func (s *SearchService) searchCustomCRUDs(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]*models.CustomCRUDData, error) {
	if req.EntityID != nil {
		// Search in specific entity
		data, err := s.customCRUDRepo.ListDataByEntity(ctx, *req.EntityID, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}

		if req.Query == "" {
			return data, nil
		}

		// Filter by query in JSON data
		var filtered []*models.CustomCRUDData
		queryLower := strings.ToLower(req.Query)
		for _, item := range data {
			if strings.Contains(strings.ToLower(item.Data), queryLower) {
				filtered = append(filtered, item)
			}
		}
		return filtered, nil
	}

	// Search across all entities
	entities, err := s.customCRUDRepo.ListEntities(ctx, nil, true)
	if err != nil {
		return nil, err
	}

	var allData []*models.CustomCRUDData
	for _, entity := range entities {
		data, err := s.customCRUDRepo.ListDataByEntity(ctx, entity.ID, 50, 0)
		if err == nil {
			allData = append(allData, data...)
		}
	}

	// Filter by query
	if req.Query != "" {
		queryLower := strings.ToLower(req.Query)
		var filtered []*models.CustomCRUDData
		for _, item := range allData {
			if strings.Contains(strings.ToLower(item.Data), queryLower) {
				filtered = append(filtered, item)
			}
			if len(filtered) >= req.Limit {
				break
			}
		}
		return filtered, nil
	}

	// Return limited results
	if len(allData) > req.Limit {
		return allData[:req.Limit], nil
	}
	return allData, nil
}

// searchByLocation searches across all entities by location
func (s *SearchService) searchByLocation(ctx context.Context, userID uuid.UUID, req SearchRequest) ([]models.SearchResult, error) {
	var results []models.SearchResult

	// Search users by location
	if req.Country != nil || req.City != nil || req.Location != nil {
		users, err := s.searchUsersByLocation(ctx, userID, req)
		if err == nil {
			for _, user := range users {
				results = append(results, models.SearchResult{
					Type:  "user",
					ID:    user.ID.String(),
					Title: user.Name,
					Data:  user,
				})
			}
		}
	}

	// Search dashboard items by location (if metadata contains location)
	if req.Location != nil {
		items, err := s.dashboardRepo.ListItems(ctx, userID, 100, 0)
		if err == nil {
			locationLower := strings.ToLower(*req.Location)
			for _, item := range items {
				if item.Metadata != nil && strings.Contains(strings.ToLower(*item.Metadata), locationLower) {
					results = append(results, models.SearchResult{
						Type:        "dashboard_item",
						ID:          item.ID.String(),
						Title:       item.Title,
						Description: item.Description,
						Data:        item,
					})
				}
			}
		}
	}

	return results, nil
}

// GetSearchHistory retrieves user's search history
func (s *SearchService) GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.SearchHistory, error) {
	return s.searchRepo.GetSearchHistory(ctx, userID, limit)
}

// ClearSearchHistory clears user's search history
func (s *SearchService) ClearSearchHistory(ctx context.Context, userID uuid.UUID) error {
	return s.searchRepo.ClearSearchHistory(ctx, userID)
}
