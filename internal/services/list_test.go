package services

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	k1 := generateKey()
	k2 := generateKey()
	if k1 == "" || k2 == "" {
		t.Fatalf("expected non-empty key")
	}
	if k1 == k2 {
		t.Fatalf("expected unique keys, got %s and %s", k1, k2)
	}
}
