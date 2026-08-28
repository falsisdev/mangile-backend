package services

import (
	"context"
	"errors"
	"strings"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

// sanityListSortFields, /api/sanityList ucundaki filterType degerlerini
// GROQ siralama alanlarina esler. Kullanici girisi dogrudan sorguya
// giremedigi icin degerler beyaz listede tutulur.
var sanityListSortFields = map[string]string{
	"createdAt": "_createdAt",
	"updatedAt": "_updatedAt",
}

// ErrInvalidSortField, gecersiz siralama alani istendiginde dondurulur.
var ErrInvalidSortField = errors.New("geçersiz sıralama alanı")

// GetSanityList, yerel tum serileri istenen alana gore azalan siralar.
func GetSanityList(ctx context.Context, filterType string) ([]models.SanityList, error) {
	sortField, ok := sanityListSortFields[filterType]
	if !ok {
		return nil, ErrInvalidSortField
	}

	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	query := strings.ReplaceAll(queries.SanityListQuery, "REPLACE_SORT_FIELD", sortField)

	var sanityList []models.SanityList
	if err := sanity.query(ctx, query, nil, &sanityList); err != nil {
		return nil, err
	}

	return sanityList, nil
}
