package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid UUID v4",
			input:    "550e8400-e29b-41d4-a716-446655440000",
			expected: true,
		},
		{
			name:     "valid UUID v1",
			input:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expected: true,
		},
		{
			name:     "valid UUID v3",
			input:    "6fa459ea-ee8a-3ca4-894e-db77e160355e",
			expected: true,
		},
		{
			name:     "valid UUID v5",
			input:    "886313e1-3b8a-5372-9b90-0c9aee199e5d",
			expected: true,
		},
		{
			name:     "valid UUID uppercase",
			input:    "550E8400-E29B-41D4-A716-446655440000",
			expected: true,
		},
		{
			name:     "valid UUID mixed case",
			input:    "550e8400-E29b-41D4-a716-446655440000",
			expected: true,
		},
		{
			name:     "valid UUID no hyphens",
			input:    "550e8400e29b41d4a716446655440000",
			expected: true,
		},
		{
			name:     "generated UUID",
			input:    uuid.New().String(),
			expected: true,
		},
		{
			name:     "nil UUID",
			input:    "00000000-0000-0000-0000-000000000000",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid format - too short",
			input:    "550e8400-e29b-41d4-a716",
			expected: false,
		},
		{
			name:     "invalid format - too long",
			input:    "550e8400-e29b-41d4-a716-446655440000-extra",
			expected: false,
		},
		{
			name:     "invalid characters",
			input:    "550e8400-e29b-41d4-a716-44665544000g",
			expected: false,
		},
		{
			name:     "missing hyphens",
			input:    "550e8400e29b41d4a716446655440000extra",
			expected: false,
		},
		{
			name:     "wrong hyphen positions",
			input:    "550e84-00e29b-41d4-a716-446655440000",
			expected: false,
		},
		{
			name:     "special characters",
			input:    "550e8400-e29b-41d4-a716-44665544000!",
			expected: false,
		},
		{
			name:     "spaces",
			input:    "550e8400-e29b-41d4-a716-446655440000 ",
			expected: false,
		},
		{
			name:     "leading spaces",
			input:    " 550e8400-e29b-41d4-a716-446655440000",
			expected: false,
		},
		{
			name:     "random string",
			input:    "not-a-uuid-at-all",
			expected: false,
		},
		{
			name:     "numeric string",
			input:    "12345678901234567890123456789012",
			expected: false,
		},
		{
			name:     "SQL injection attempt",
			input:    "'; DROP TABLE users; --",
			expected: false,
		},
		{
			name:     "URN format valid",
			input:    "urn:uuid:550e8400-e29b-41d4-a716-446655440000",
			expected: false, // URN format not supported by uuid.Parse
		},
		{
			name:     "curly braces format",
			input:    "{550e8400-e29b-41d4-a716-446655440000}",
			expected: true, // uuid.Parse accepts this format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUUID(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateUUID_MultipleFormats(t *testing.T) {
	// Test that different valid representations of the same UUID all validate
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	formats := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"550E8400-E29B-41D4-A716-446655440000",
		"550e8400e29b41d4a716446655440000",
		"{550e8400-e29b-41d4-a716-446655440000}",
	}

	for _, format := range formats {
		if !ValidateUUID(format) {
			t.Errorf("ValidateUUID(%q) should be valid for UUID %q", format, validUUID)
		}
	}
}

func TestValidateUUID_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "all zeros",
			input:    "00000000-0000-0000-0000-000000000000",
			expected: true,
		},
		{
			name:     "all F's",
			input:    "ffffffff-ffff-ffff-ffff-ffffffffffff",
			expected: true,
		},
		{
			name:     "all F's uppercase",
			input:    "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF",
			expected: true,
		},
		{
			name:     "one character off",
			input:    "ffffffff-ffff-ffff-ffff-fffffffffffff",
			expected: false,
		},
		{
			name:     "one character short",
			input:    "ffffffff-ffff-ffff-ffff-fffffffffff",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUUID(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkValidateUUID_Valid(b *testing.B) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUUID(validUUID)
	}
}

func BenchmarkValidateUUID_Invalid(b *testing.B) {
	invalidUUID := "not-a-valid-uuid"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUUID(invalidUUID)
	}
}

func BenchmarkValidateUUID_Empty(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUUID("")
	}
}
