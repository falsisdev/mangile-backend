package services

import (
	"testing"
)

func TestStatusNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"reading", "reading"},
		{"current", "reading"},
		{"completed", "completed"},
		{"on_hold", "on_hold"},
		{"paused", "on_hold"},
		{"dropped", "dropped"},
		{"plan_to_read", "plan_to_read"},
		{"planning", "plan_to_read"},
		{"repeating", "rereading"},
		{"re_reading", "rereading"},
	}

	for _, tt := range tests {
		got := normalizeStatus(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeStatus(%q) = %q; expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestToMALStatus(t *testing.T) {
	if got := toMALStatus("reading"); got != "reading" {
		t.Errorf("expected reading, got %s", got)
	}
	if got := toMALStatus("completed"); got != "completed" {
		t.Errorf("expected completed, got %s", got)
	}
	if got := toMALStatus("plan_to_read"); got != "plan_to_read" {
		t.Errorf("expected plan_to_read, got %s", got)
	}
}

func TestToAniListStatus(t *testing.T) {
	if got := toAniListStatus("reading"); got != "CURRENT" {
		t.Errorf("expected CURRENT, got %s", got)
	}
	if got := toAniListStatus("completed"); got != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", got)
	}
	if got := toAniListStatus("plan_to_read"); got != "PLANNING" {
		t.Errorf("expected PLANNING, got %s", got)
	}
}

func TestToMangaDexStatus(t *testing.T) {
	if got := toMangaDexStatus("reading"); got != "reading" {
		t.Errorf("expected reading, got %s", got)
	}
	if got := toMangaDexStatus("rereading"); got != "re_reading" {
		t.Errorf("expected re_reading, got %s", got)
	}
}

func TestGeneratePKCE(t *testing.T) {
	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE failed: %v", err)
	}
	if len(verifier) < 32 {
		t.Errorf("verifier too short: %s", verifier)
	}
	if challenge == "" {
		t.Errorf("challenge is empty")
	}
}

func TestExtractMDTitle(t *testing.T) {
	titles := map[string]string{
		"en": "Attack on Titan",
	}
	alt := []map[string]string{
		{"ja": "進撃の巨人"},
	}
	if got := extractMDTitle(titles, alt); got != "Attack on Titan" {
		t.Errorf("expected Attack on Titan, got %s", got)
	}

	// Fallback to ja-ro if no en
	titles2 := map[string]string{
		"ja-ro": "Shingeki no Kyojin",
	}
	if got := extractMDTitle(titles2, alt); got != "Shingeki no Kyojin" {
		t.Errorf("expected Shingeki no Kyojin, got %s", got)
	}
}
