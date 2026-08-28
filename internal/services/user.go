package services

import (
	"context"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// GetUser, verilen Logto kimligine sahip kullaniciyi ve listelerini dondurur.
func GetUser(ctx context.Context, id string) (*models.User, error) {
	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := sanity.query(ctx, queries.UserQuery, map[string]any{"id": id}, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
