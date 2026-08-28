package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetArticle, verilen slug'a sahip makaleyi dondurur.
func GetArticle(ctx context.Context, slug string) (*models.Article, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var article models.Article
	if err := sanity.query(ctx, queries.ArticleQuery, map[string]any{"slug": slug}, &article); err != nil {
		return nil, err
	}

	return &article, nil
}
