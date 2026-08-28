package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetScan, verilen kimlikteki ceviri ekibini uyeleriyle birlikte dondurur.
func GetScan(ctx context.Context, id string) (*models.Scan, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var scan models.Scan
	if err := sanity.query(ctx, queries.ScanQuery, map[string]any{"id": id}, &scan); err != nil {
		return nil, err
	}

	return &scan, nil
}
