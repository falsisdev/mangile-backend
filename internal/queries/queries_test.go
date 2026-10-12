package queries_test

import (
	"strings"
	"testing"

	"github.com/falsisdev/mangile-backend/internal/queries"
)

func TestEmbeddedListQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "ListQuery",
			query:    queries.ListQuery,
			expected: []string{"_type == \"lists\"", "$id", "items", "likes", "likesCount", "itemsCount"},
		},
		{
			name:     "PublicListsQuery",
			query:    queries.PublicListsQuery,
			expected: []string{"_type == \"lists\"", "$limit", "isPrivate"},
		},
		{
			name:     "UserListsQuery",
			query:    queries.UserListsQuery,
			expected: []string{"_type == \"lists\"", "$logtoId", "$includePrivate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.TrimSpace(tt.query) == "" {
				t.Fatalf("query %s is empty", tt.name)
			}
			for _, exp := range tt.expected {
				if !strings.Contains(tt.query, exp) {
					t.Errorf("query %s expected to contain %q, but did not", tt.name, exp)
				}
			}
		})
	}
}
