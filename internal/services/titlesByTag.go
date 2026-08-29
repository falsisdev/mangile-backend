package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetTitlesByTag, verilen etikete sahip yerel serileri (manga/lightNovel)
// döndürür; etiket, sorgu metnine serpistirilmeden $tag GROQ parametresi
// olarak gönderilir.
func GetTitlesByTag(ctx context.Context, tag string) ([]models.TitleByTag, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var titles []models.TitleByTag
	if err := sanity.query(ctx, queries.TitlesByTagQuery, map[string]any{"tag": tag}, &titles); err != nil {
		return nil, err
	}

	return titles, nil
}
