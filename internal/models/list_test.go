package models_test

import (
	"encoding/json"
	"testing"

	"github.com/falsisdev/mangile-backend/internal/models"
)

func TestListItemSerialization(t *testing.T) {
	extID := 12345
	item := models.ListItem{
		Key: "k1",
		Reference: &models.ListItemReference{
			ID:    "manga-1",
			Type:  "manga",
			Title: "Test Manga",
			Slug:  "test-manga",
		},
		ExternalID:   &extID,
		ExternalType: "manga",
		Index:        1,
		Note:         "Harika bir seri",
		AddedAt:      "2026-10-12T10:00:00Z",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("failed to marshal ListItem: %v", err)
	}

	var parsed models.ListItem
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal ListItem: %v", err)
	}

	if parsed.Key != "k1" {
		t.Errorf("expected Key 'k1', got '%s'", parsed.Key)
	}
	if parsed.Index != 1 {
		t.Errorf("expected Index 1, got %d", parsed.Index)
	}
	if parsed.Reference == nil || parsed.Reference.Title != "Test Manga" {
		t.Errorf("expected Reference.Title 'Test Manga', got %+v", parsed.Reference)
	}
	if parsed.ExternalID == nil || *parsed.ExternalID != 12345 {
		t.Errorf("expected ExternalID 12345, got %v", parsed.ExternalID)
	}
	if parsed.Note != "Harika bir seri" {
		t.Errorf("expected Note 'Harika bir seri', got '%s'", parsed.Note)
	}
}

func TestListSerialization(t *testing.T) {
	list := models.List{
		ID:          "list-123",
		Type:        "lists",
		Title:       "Okuma Listem",
		Description: "Favori serilerim",
		IsPrivate:   false,
		CreatedAt:   "2026-10-12T10:00:00Z",
		UpdatedAt:   "2026-10-12T10:05:00Z",
		User: &models.ListUser{
			ID:       "auth-1",
			LogtoID:  "logto-abc",
			Name:     "Test User",
			Username: "testuser",
		},
		Likes: []models.ListUser{
			{ID: "auth-2", LogtoID: "logto-xyz", Name: "Fan", Username: "fan"},
		},
		Items: []models.ListItem{
			{Key: "item-1", Index: 1, Note: "Not 1"},
			{Key: "item-2", Index: 2, Note: "Not 2"},
		},
		LikesCount: 1,
		ItemsCount: 2,
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal List: %v", err)
	}

	var parsed models.List
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal List: %v", err)
	}

	if parsed.ID != "list-123" {
		t.Errorf("expected ID 'list-123', got '%s'", parsed.ID)
	}
	if parsed.ItemsCount != 2 || len(parsed.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(parsed.Items))
	}
	if parsed.LikesCount != 1 || len(parsed.Likes) != 1 {
		t.Errorf("expected 1 like, got %d", len(parsed.Likes))
	}
}
