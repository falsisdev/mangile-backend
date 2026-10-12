package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

func GetLatestChapters(ctx context.Context) ([]models.LatestChapter, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var latestChapters []models.LatestChapter
	if err := sanity.query(ctx, queries.LatestChaptersQuery, nil, &latestChapters); err != nil {
		return nil, err
	}

	return latestChapters, nil
}
