package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetLatestTitles, yerel tum serileri en son eklenene gore azalan siralar.
func GetLatestTitles(ctx context.Context) ([]models.LatestTitle, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var latestTitles []models.LatestTitle
	if err := sanity.query(ctx, queries.LatestTitlesQuery, nil, &latestTitles); err != nil {
		return nil, err
	}

	return latestTitles, nil
}
