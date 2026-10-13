package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

const maxHistoryItems = 100

func itemToMap(item models.ReadingHistoryItem) map[string]any {
	key := item.Key
	if key == "" {
		key = generateKey()
	}
	m := map[string]any{
		"_key":          key,
		"seriesId":      item.SeriesID,
		"seriesTitle":   item.SeriesTitle,
		"seriesType":    item.SeriesType,
		"chapterId":     item.ChapterID,
		"chapterNumber": item.ChapterNumber,
		"readAt":        item.ReadAt,
	}
	if item.SeriesFormat != "" {
		m["seriesFormat"] = item.SeriesFormat
	}
	if item.CoverURL != "" {
		m["coverUrl"] = item.CoverURL
	}
	if item.VolumeNumber > 0 {
		m["volumeNumber"] = item.VolumeNumber
	}
	if item.ChapterTitle != "" {
		m["chapterTitle"] = item.ChapterTitle
	}
	if item.CurrentPage > 0 {
		m["currentPage"] = item.CurrentPage
	}
	if item.TotalPages > 0 {
		m["totalPages"] = item.TotalPages
	}
	if item.Progress > 0 {
		m["progress"] = item.Progress
	}
	return m
}

func GetReadingHistory(ctx context.Context, logtoID string) (*models.ReadingHistoryDocument, error) {
	if logtoID == "" {
		return nil, ErrInvalidRequest
	}

	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var history models.ReadingHistoryDocument
	if err := sanity.query(ctx, queries.HistoryQuery, map[string]any{"logtoId": logtoID}, &history); err != nil {
		return nil, err
	}

	if history.Items == nil {
		history.Items = []models.ReadingHistoryItem{}
	}

	return &history, nil
}

func UpsertReadingHistoryItem(ctx context.Context, logtoID string, item models.ReadingHistoryItem) (*models.ReadingHistoryDocument, error) {
	if logtoID == "" || item.SeriesID == "" || item.ChapterID == "" {
		return nil, fmt.Errorf("%w: logtoId, seriesId ve chapterId zorunludur", ErrInvalidRequest)
	}

	user, err := GetUser(ctx, logtoID)
	if err != nil || user == nil || user.ID == "" {
		return nil, ErrUserNotFound
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if item.ReadAt == "" {
		item.ReadAt = now
	}
	if item.Key == "" {
		item.Key = generateKey()
	}

	existing, err := GetReadingHistory(ctx, logtoID)
	if err != nil {
		return nil, err
	}

	// Filter out any existing item with the same seriesId
	var updatedItems []models.ReadingHistoryItem
	updatedItems = append(updatedItems, item)
	for _, it := range existing.Items {
		if it.SeriesID != item.SeriesID {
			updatedItems = append(updatedItems, it)
		}
	}

	if len(updatedItems) > maxHistoryItems {
		updatedItems = updatedItems[:maxHistoryItems]
	}

	var itemsPayload []any
	for _, it := range updatedItems {
		itemsPayload = append(itemsPayload, itemToMap(it))
	}

	if existing.ID == "" {
		// Create new readingHistory document
		doc := map[string]any{
			"_type":     "readingHistory",
			"updatedAt": now,
			"user": map[string]any{
				"_type": "reference",
				"_ref":  user.ID,
			},
			"items": itemsPayload,
		}

		res, err := MutateSanity(ctx, []map[string]any{{"create": doc}})
		if err != nil {
			return nil, fmt.Errorf("okuma gecmisi olusturulamadi: %w", err)
		}
		docID := ""
		if len(res.Results) > 0 {
			docID = res.Results[0].ID
		}
		return &models.ReadingHistoryDocument{
			ID:        docID,
			Type:      "readingHistory",
			UpdatedAt: now,
			User: &models.ReadingHistoryUser{
				ID:       user.ID,
				LogtoID:  user.LogtoID,
				Name:     user.Name,
				Username: user.Username,
			},
			Items: updatedItems,
		}, nil
	}

	// Patch existing document
	patch := map[string]any{
		"patch": map[string]any{
			"id": existing.ID,
			"set": map[string]any{
				"items":     itemsPayload,
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patch}); err != nil {
		return nil, fmt.Errorf("okuma gecmisi guncellenemedi: %w", err)
	}

	existing.UpdatedAt = now
	existing.Items = updatedItems
	return existing, nil
}

func RemoveReadingHistoryItem(ctx context.Context, logtoID string, seriesID string) error {
	if logtoID == "" || seriesID == "" {
		return ErrInvalidRequest
	}

	existing, err := GetReadingHistory(ctx, logtoID)
	if err != nil {
		return err
	}

	if existing.ID == "" || len(existing.Items) == 0 {
		return nil
	}

	var updatedItems []any
	found := false
	for _, it := range existing.Items {
		if it.SeriesID == seriesID {
			found = true
			continue
		}
		updatedItems = append(updatedItems, itemToMap(it))
	}

	if !found {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	patch := map[string]any{
		"patch": map[string]any{
			"id": existing.ID,
			"set": map[string]any{
				"items":     updatedItems,
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patch}); err != nil {
		return fmt.Errorf("gecmis ogesi silinemedi: %w", err)
	}

	return nil
}

func ClearReadingHistory(ctx context.Context, logtoID string) error {
	if logtoID == "" {
		return ErrInvalidRequest
	}

	existing, err := GetReadingHistory(ctx, logtoID)
	if err != nil {
		return err
	}

	if existing.ID == "" {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	patch := map[string]any{
		"patch": map[string]any{
			"id": existing.ID,
			"set": map[string]any{
				"items":     []any{},
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patch}); err != nil {
		return fmt.Errorf("okuma gecmisi temizlenemedi: %w", err)
	}

	return nil
}

func SyncReadingHistory(ctx context.Context, logtoID string, incomingItems []models.ReadingHistoryItem) (*models.ReadingHistoryDocument, error) {
	if logtoID == "" {
		return nil, ErrInvalidRequest
	}

	user, err := GetUser(ctx, logtoID)
	if err != nil || user == nil || user.ID == "" {
		return nil, ErrUserNotFound
	}

	existing, err := GetReadingHistory(ctx, logtoID)
	if err != nil {
		return nil, err
	}

	// Merge by seriesId, keeping the one with newer readAt
	seriesMap := make(map[string]models.ReadingHistoryItem)
	for _, it := range existing.Items {
		if it.SeriesID != "" {
			seriesMap[it.SeriesID] = it
		}
	}

	for _, incoming := range incomingItems {
		if incoming.SeriesID == "" {
			continue
		}
		if current, exists := seriesMap[incoming.SeriesID]; exists {
			// Compare readAt timestamps
			if incoming.ReadAt > current.ReadAt {
				seriesMap[incoming.SeriesID] = incoming
			}
		} else {
			seriesMap[incoming.SeriesID] = incoming
		}
	}

	merged := make([]models.ReadingHistoryItem, 0, len(seriesMap))
	for _, it := range seriesMap {
		merged = append(merged, it)
	}

	// Sort descending by readAt
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].ReadAt > merged[j].ReadAt
	})

	if len(merged) > maxHistoryItems {
		merged = merged[:maxHistoryItems]
	}

	now := time.Now().UTC().Format(time.RFC3339)
	var itemsPayload []any
	for _, it := range merged {
		itemsPayload = append(itemsPayload, itemToMap(it))
	}

	if existing.ID == "" {
		doc := map[string]any{
			"_type":     "readingHistory",
			"updatedAt": now,
			"user": map[string]any{
				"_type": "reference",
				"_ref":  user.ID,
			},
			"items": itemsPayload,
		}
		res, err := MutateSanity(ctx, []map[string]any{{"create": doc}})
		if err != nil {
			return nil, fmt.Errorf("okuma gecmisi olusturulamadi: %w", err)
		}
		docID := ""
		if len(res.Results) > 0 {
			docID = res.Results[0].ID
		}
		return &models.ReadingHistoryDocument{
			ID:        docID,
			Type:      "readingHistory",
			UpdatedAt: now,
			User: &models.ReadingHistoryUser{
				ID:       user.ID,
				LogtoID:  user.LogtoID,
				Name:     user.Name,
				Username: user.Username,
			},
			Items: merged,
		}, nil
	}

	patch := map[string]any{
		"patch": map[string]any{
			"id": existing.ID,
			"set": map[string]any{
				"items":     itemsPayload,
				"updatedAt": now,
			},
		},
	}
	if _, err := MutateSanity(ctx, []map[string]any{patch}); err != nil {
		return nil, fmt.Errorf("okuma gecmisi senkronize edilemedi: %w", err)
	}

	existing.UpdatedAt = now
	existing.Items = merged
	return existing, nil
}
