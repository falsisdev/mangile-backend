package services

import (
	"context"
	"strings"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// LocalTitlesFilter, yerel serileri sorgularken kullanılabilecek filtre ve sayfalama parametrelerini tutar.
type LocalTitlesFilter struct {
	Search string
	Type   string
	Tag    string
	Status string
	Sort   string
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

	var orderClause string
	switch filter.Sort {
	case "title_asc":
		orderClause = "order(title asc)"
	case "title_desc":
		orderClause = "order(title desc)"
	case "date_asc":
		orderClause = "order(_createdAt asc)"
	default:
		orderClause = "order(_createdAt desc)"
	}

	query := strings.Replace(queries.LocalTitlesQuery, "%ORDER%", orderClause, 1)

	searchParam := filter.Search
	if searchParam != "" {
		searchParam = "*" + searchParam + "*"
	}

	params := map[string]any{
		"search": searchParam,
		"type":   filter.Type,
		"tag":    filter.Tag,
		"status": filter.Status,
		"start":  start,
		"end":    end,
	}

	var result models.LocalTitlesQueryResult
	if err := sanity.query(ctx, query, params, &result); err != nil {
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
