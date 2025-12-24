package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type SearchService struct {
	searchRepo repositories.SearchRepository
	logger     *zap.Logger
}

func NewSearchService(searchRepo repositories.SearchRepository, logger *zap.Logger) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
		logger:     logger,
	}
}

// Search performs a comprehensive search across all searchable entities
func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, query string, searchType string, limit int) (*models.SearchResult, error) {
	if query == "" {
		return nil, nil
	}

	if limit <= 0 {
		limit = 20
	}

	var results []models.SearchResult
	var totalCount int

	switch searchType {
	case "dashboard_items", "all":
		items, err := s.searchRepo.SearchDashboardItems(ctx, userID, query, limit)
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

	switch searchType {
	case "messages", "all":
		messages, err := s.searchRepo.SearchMessages(ctx, userID, query, limit)
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

	switch searchType {
	case "users", "all":
		users, err := s.searchRepo.SearchUsers(ctx, userID, query, limit)
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

	// Save search history
	history := &models.SearchHistory{
		ID:           uuid.New(),
		UserID:       userID,
		Query:        query,
		SearchType:   &searchType,
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
			"query":   query,
			"type":    searchType,
		},
	}, nil
}

