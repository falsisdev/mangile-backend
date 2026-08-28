package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetList, verilen kimlikteki listeyi sahibi ve begenileriyle birlikte dondurur.
func GetList(ctx context.Context, id string) (*models.List, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var list models.List
	if err := sanity.query(ctx, queries.ListQuery, map[string]any{"id": id}, &list); err != nil {
		return nil, err
	}

	return &list, nil
}
