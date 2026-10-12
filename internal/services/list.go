package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/falsisdev/mangile-backend/internal/models"
	"github.com/falsisdev/mangile-backend/internal/queries"
)

var (
	ErrNotFound       = errors.New("liste bulunamadi")
	ErrUnauthorized   = errors.New("bu islem icin yetkiniz bulunmuyor")
	ErrListIsPrivate  = errors.New("bu liste gizlidir")
	ErrUserNotFound   = errors.New("kullanici bulunamadi")
	ErrInvalidRequest = errors.New("gecersiz istek parametreleri")
)

func generateKey() string {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

func GetList(ctx context.Context, id string, viewerLogtoID string) (*models.List, error) {
	if id == "" {
		return nil, ErrInvalidRequest
	}

	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var list models.List
	if err := sanity.query(ctx, queries.ListQuery, map[string]any{"id": id}, &list); err != nil {
		return nil, err
	}

	if list.ID == "" {
		return nil, ErrNotFound
	}

	if list.IsPrivate {
		if list.User == nil || list.User.LogtoID != viewerLogtoID {
			return nil, ErrListIsPrivate
		}
	}

	return &list, nil
}

func GetPublicLists(ctx context.Context, limit int) ([]models.List, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	var lists []models.List
	if err := sanity.query(ctx, queries.PublicListsQuery, map[string]any{"limit": limit}, &lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func GetUserLists(ctx context.Context, targetLogtoID string, viewerLogtoID string) ([]models.List, error) {
	if targetLogtoID == "" {
		return nil, ErrInvalidRequest
	}

	sanity, err := newSanityClient()
	if err != nil {
		return nil, err
	}

	includePrivate := (viewerLogtoID != "" && targetLogtoID == viewerLogtoID)

	var lists []models.List
	params := map[string]any{
		"logtoId":        targetLogtoID,
		"includePrivate": includePrivate,
	}
	if err := sanity.query(ctx, queries.UserListsQuery, params, &lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func CreateList(ctx context.Context, req models.CreateListRequest) (*models.List, error) {
	if req.LogtoID == "" || req.Title == "" {
		return nil, fmt.Errorf("%w: logtoId ve title alanlari zorunludur", ErrInvalidRequest)
	}

	user, err := GetUser(ctx, req.LogtoID)
	if err != nil || user == nil || user.ID == "" {
		return nil, ErrUserNotFound
	}

	now := time.Now().UTC().Format(time.RFC3339)
	doc := map[string]any{
		"_type":       "lists",
		"title":       req.Title,
		"description": req.Description,
		"isPrivate":   req.IsPrivate,
		"createdAt":   now,
		"updatedAt":   now,
		"user": map[string]any{
			"_type": "reference",
			"_ref":  user.ID,
		},
		"likes": []any{},
		"items": []any{},
	}

	createResp, err := MutateSanity(ctx, []map[string]any{
		{"create": doc},
	})
	if err != nil {
		return nil, fmt.Errorf("liste olusturulamadi: %w", err)
	}

	if len(createResp.Results) == 0 || createResp.Results[0].ID == "" {
		return nil, fmt.Errorf("sanity liste kimligi dondurmedi")
	}

	createdListID := createResp.Results[0].ID

	patchUser := map[string]any{
		"patch": map[string]any{
			"id":           user.ID,
			"setIfMissing": map[string]any{"lists": []any{}},
			"insert": map[string]any{
				"after": "lists[-1]",
				"items": []any{
					map[string]any{
						"_type": "reference",
						"_ref":  createdListID,
						"_key":  generateKey(),
					},
				},
			},
		},
	}
	_, _ = MutateSanity(ctx, []map[string]any{patchUser})

	return GetList(ctx, createdListID, req.LogtoID)
}

func UpdateList(ctx context.Context, listID string, req models.UpdateListRequest) (*models.List, error) {
	if listID == "" || req.LogtoID == "" {
		return nil, ErrInvalidRequest
	}

	existing, err := GetList(ctx, listID, req.LogtoID)
	if err != nil {
		return nil, err
	}

	if existing.User == nil || existing.User.LogtoID != req.LogtoID {
		return nil, ErrUnauthorized
	}

	setFields := map[string]any{
		"updatedAt": time.Now().UTC().Format(time.RFC3339),
	}
	if req.Title != nil {
		setFields["title"] = *req.Title
	}
	if req.Description != nil {
		setFields["description"] = *req.Description
	}
	if req.IsPrivate != nil {
		setFields["isPrivate"] = *req.IsPrivate
	}

	patchMutation := map[string]any{
		"patch": map[string]any{
			"id":  listID,
			"set": setFields,
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patchMutation}); err != nil {
		return nil, fmt.Errorf("liste guncellenemedi: %w", err)
	}

	return GetList(ctx, listID, req.LogtoID)
}

func DeleteList(ctx context.Context, listID string, logtoID string) error {
	if listID == "" || logtoID == "" {
		return ErrInvalidRequest
	}

	existing, err := GetList(ctx, listID, logtoID)
	if err != nil {
		return err
	}

	if existing.User == nil || existing.User.LogtoID != logtoID {
		return ErrUnauthorized
	}

	mutations := []map[string]any{
		{"delete": map[string]any{"id": listID}},
	}

	if existing.User.ID != "" {
		mutations = append(mutations, map[string]any{
			"patch": map[string]any{
				"id":    existing.User.ID,
				"unset": []string{fmt.Sprintf("lists[_ref==\"%s\"]", listID)},
			},
		})
	}

	if _, err := MutateSanity(ctx, mutations); err != nil {
		return fmt.Errorf("liste silinemedi: %w", err)
	}

	return nil
}

func AddListItem(ctx context.Context, listID string, req models.AddListItemRequest) (*models.List, error) {
	if listID == "" || req.LogtoID == "" {
		return nil, ErrInvalidRequest
	}

	existing, err := GetList(ctx, listID, req.LogtoID)
	if err != nil {
		return nil, err
	}

	if existing.User == nil || existing.User.LogtoID != req.LogtoID {
		return nil, ErrUnauthorized
	}

	index := len(existing.Items) + 1
	if req.Index != nil {
		index = *req.Index
	}

	now := time.Now().UTC().Format(time.RFC3339)
	itemMap := map[string]any{
		"_key":    generateKey(),
		"index":   index,
		"addedAt": now,
	}

	if req.ReferenceID != nil && *req.ReferenceID != "" {
		itemMap["reference"] = map[string]any{
			"_type": "reference",
			"_ref":  *req.ReferenceID,
		}
	}
	if req.ExternalID != nil && *req.ExternalID != 0 {
		itemMap["externalId"] = *req.ExternalID
	}
	if req.ExternalType != nil && *req.ExternalType != "" {
		itemMap["externalType"] = *req.ExternalType
	}
	if req.Note != "" {
		itemMap["note"] = req.Note
	}

	patchMutation := map[string]any{
		"patch": map[string]any{
			"id":           listID,
			"setIfMissing": map[string]any{"items": []any{}},
			"insert": map[string]any{
				"after": "items[-1]",
				"items": []any{itemMap},
			},
			"set": map[string]any{
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patchMutation}); err != nil {
		return nil, fmt.Errorf("oge listeye eklenemedi: %w", err)
	}

	return GetList(ctx, listID, req.LogtoID)
}

func RemoveListItem(ctx context.Context, listID string, itemKey string, logtoID string) (*models.List, error) {
	if listID == "" || itemKey == "" || logtoID == "" {
		return nil, ErrInvalidRequest
	}

	existing, err := GetList(ctx, listID, logtoID)
	if err != nil {
		return nil, err
	}

	if existing.User == nil || existing.User.LogtoID != logtoID {
		return nil, ErrUnauthorized
	}

	now := time.Now().UTC().Format(time.RFC3339)
	patchMutation := map[string]any{
		"patch": map[string]any{
			"id":    listID,
			"unset": []string{fmt.Sprintf("items[_key==\"%s\"]", itemKey)},
			"set": map[string]any{
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patchMutation}); err != nil {
		return nil, fmt.Errorf("oge listeden kaldirilamadi: %w", err)
	}

	return GetList(ctx, listID, logtoID)
}

func UpdateListItems(ctx context.Context, listID string, req models.UpdateListItemsRequest) (*models.List, error) {
	if listID == "" || req.LogtoID == "" {
		return nil, ErrInvalidRequest
	}

	existing, err := GetList(ctx, listID, req.LogtoID)
	if err != nil {
		return nil, err
	}

	if existing.User == nil || existing.User.LogtoID != req.LogtoID {
		return nil, ErrUnauthorized
	}

	formattedItems := make([]map[string]any, 0, len(req.Items))
	for i, item := range req.Items {
		key := item.Key
		if key == "" {
			key = generateKey()
		}
		idx := item.Index
		if idx <= 0 {
			idx = i + 1
		}

		itemMap := map[string]any{
			"_key":  key,
			"index": idx,
		}

		if item.AddedAt != "" {
			itemMap["addedAt"] = item.AddedAt
		} else {
			itemMap["addedAt"] = time.Now().UTC().Format(time.RFC3339)
		}

		if item.Reference != nil && item.Reference.ID != "" {
			itemMap["reference"] = map[string]any{
				"_type": "reference",
				"_ref":  item.Reference.ID,
			}
		}
		if item.ExternalID != nil && *item.ExternalID != 0 {
			itemMap["externalId"] = *item.ExternalID
		}
		if item.ExternalType != "" {
			itemMap["externalType"] = item.ExternalType
		}
		if item.Note != "" {
			itemMap["note"] = item.Note
		}

		formattedItems = append(formattedItems, itemMap)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	patchMutation := map[string]any{
		"patch": map[string]any{
			"id": listID,
			"set": map[string]any{
				"items":     formattedItems,
				"updatedAt": now,
			},
		},
	}

	if _, err := MutateSanity(ctx, []map[string]any{patchMutation}); err != nil {
		return nil, fmt.Errorf("liste ogeleri guncellenemedi: %w", err)
	}

	return GetList(ctx, listID, req.LogtoID)
}

func ToggleLikeList(ctx context.Context, listID string, logtoID string) (*models.List, error) {
	if listID == "" || logtoID == "" {
		return nil, ErrInvalidRequest
	}

	user, err := GetUser(ctx, logtoID)
	if err != nil || user == nil || user.ID == "" {
		return nil, ErrUserNotFound
	}

	list, err := GetList(ctx, listID, logtoID)
	if err != nil {
		return nil, err
	}

	isLiked := false
	for _, likeUser := range list.Likes {
		if likeUser.LogtoID == logtoID || likeUser.ID == user.ID {
			isLiked = true
			break
		}
	}

	var patchMutation map[string]any
	if isLiked {
		patchMutation = map[string]any{
			"patch": map[string]any{
				"id":    listID,
				"unset": []string{fmt.Sprintf("likes[_ref==\"%s\"]", user.ID)},
			},
		}
	} else {
		patchMutation = map[string]any{
			"patch": map[string]any{
				"id":           listID,
				"setIfMissing": map[string]any{"likes": []any{}},
				"insert": map[string]any{
					"after": "likes[-1]",
					"items": []any{
						map[string]any{
							"_type": "reference",
							"_ref":  user.ID,
							"_key":  generateKey(),
						},
					},
				},
			},
		}
	}

	if _, err := MutateSanity(ctx, []map[string]any{patchMutation}); err != nil {
		return nil, fmt.Errorf("begeni islemi tamamlanamadi: %w", err)
	}

	return GetList(ctx, listID, logtoID)
}
