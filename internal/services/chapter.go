package services

import (
	"context"
	"errors"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// ErrChapterNotFound, istenen bölüm dokümanı Sanity'de bulunamadığında döndürülür.
var ErrChapterNotFound = errors.New("bölüm bulunamadı")

// GetChapter, verilen kimlikteki manga/novel bölümünü, bağlı olduğu eserle ve
// aynı eserin tüm bölümleriyle birlikte döndürür.
func GetChapter(ctx context.Context, id string) (models.ChapterDetails, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return models.ChapterDetails{}, err
	}

	var results []models.ChapterDetails
	if err := sanity.query(ctx, queries.ChapterQuery, map[string]any{"id": id}, &results); err != nil {
		return models.ChapterDetails{}, err
	}

	if len(results) == 0 {
		return models.ChapterDetails{}, ErrChapterNotFound
	}

	return results[0], nil
}
