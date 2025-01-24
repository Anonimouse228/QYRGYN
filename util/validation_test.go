package util

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"example@example.com", true},
		{"user.name+tag+sorting@example.com", true},
		{"user@sub.example.com", true},
		{"user@example", false},
		{"user@.com", false},
		{"@example.com", false},
		{"user@exam_ple.com", false},
		{"user@exam+ple.com", false},
		{"user@exam!ple.com", false},
		{"user@com.", false},
		{"user@.com.", false},
		{"user@domain.c", false},
		{"user@domain.corporate", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v; want %v", tt.email, result, tt.expected)
			}
		})
	}
}
