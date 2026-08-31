package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// LocalTitlesFilter, yerel serileri sorgularken kullanılabilecek filtre ve sayfalama parametrelerini tutar.
type LocalTitlesFilter struct {
	Search string
	Type   string
	Tag    string
	Status string
	Page   int
	Limit  int
}

// GetLocalTitles, filtreleme ve sayfalama ile Sanity veritabanından serileri çeker.
func GetLocalTitles(ctx context.Context, filter LocalTitlesFilter) (*models.LocalTitlesResponse, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	start := (filter.Page - 1) * filter.Limit
	end := start + filter.Limit - 1

	params := map[string]any{
		"search": filter.Search,
		"type":   filter.Type,
		"tag":    filter.Tag,
		"status": filter.Status,
		"start":  start,
		"end":    end,
	}

	var result models.LocalTitlesQueryResult
	if err := sanity.query(ctx, queries.LocalTitlesQuery, params, &result); err != nil {
		return nil, err
	}

	if result.Data == nil {
		result.Data = []models.LocalTitle{}
	}

	return &models.LocalTitlesResponse{
		Data:  result.Data,
		Total: result.Total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}
