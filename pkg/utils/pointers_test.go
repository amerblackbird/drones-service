package utils

import "testing"

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty string", input: ""},
		{name: "simple string", input: "hello"},
		{name: "string with spaces", input: "hello world"},
		{name: "string with special chars", input: "!@#$%^&*()"},
		{name: "long string", input: "this is a very long string with many characters"},
		{name: "unicode string", input: "こんにちは世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringPtr(tt.input)

			if result == nil {
				t.Error("StringPtr returned nil")
				return
			}

			if *result != tt.input {
				t.Errorf("StringPtr(%q) = %q, want %q", tt.input, *result, tt.input)
			}

			// Verify it's actually a pointer to the value
			if &tt.input == result {
				t.Error("StringPtr returned pointer to input variable instead of new allocation")
			}
		})
	}
}

func TestStringPtr_Mutability(t *testing.T) {
	original := "original"
	ptr := StringPtr(original)

	// Modify the dereferenced value
	*ptr = "modified"

	// Original should not change (pointer to copy, not original)
	if original != "original" {
		t.Errorf("Modifying pointer affected original: original = %q", original)
	}

	if *ptr != "modified" {
		t.Errorf("Pointer value not modified: *ptr = %q", *ptr)
	}
}

func TestIntPtr(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{name: "zero", input: 0},
		{name: "positive", input: 42},
		{name: "negative", input: -42},
		{name: "max int32", input: 2147483647},
		{name: "min int32", input: -2147483648},
		{name: "large positive", input: 999999999},
		{name: "large negative", input: -999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntPtr(tt.input)

			if result == nil {
				t.Error("IntPtr returned nil")
				return
			}

			if *result != tt.input {
				t.Errorf("IntPtr(%d) = %d, want %d", tt.input, *result, tt.input)
			}

			// Verify it's actually a pointer to the value
			if &tt.input == result {
				t.Error("IntPtr returned pointer to input variable instead of new allocation")
			}
		})
	}
}

func TestIntPtr_Mutability(t *testing.T) {
	original := 100
	ptr := IntPtr(original)

	// Modify the dereferenced value
	*ptr = 200

	// Original should not change (pointer to copy, not original)
	if original != 100 {
		t.Errorf("Modifying pointer affected original: original = %d", original)
	}

	if *ptr != 200 {
		t.Errorf("Pointer value not modified: *ptr = %d", *ptr)
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{name: "true", input: true},
		{name: "false", input: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolPtr(tt.input)

			if result == nil {
				t.Error("BoolPtr returned nil")
				return
			}

			if *result != tt.input {
				t.Errorf("BoolPtr(%v) = %v, want %v", tt.input, *result, tt.input)
			}

			// Verify it's actually a pointer to the value
			if &tt.input == result {
				t.Error("BoolPtr returned pointer to input variable instead of new allocation")
			}
		})
	}
}

func TestBoolPtr_Mutability(t *testing.T) {
	original := true
	ptr := BoolPtr(original)

	// Modify the dereferenced value
	*ptr = false

	// Original should not change (pointer to copy, not original)
	if original != true {
		t.Errorf("Modifying pointer affected original: original = %v", original)
	}

	if *ptr != false {
		t.Errorf("Pointer value not modified: *ptr = %v", *ptr)
	}
}

func TestPointers_UniqueAddresses(t *testing.T) {
	// Test that multiple calls return different pointers
	t.Run("StringPtr unique addresses", func(t *testing.T) {
		ptr1 := StringPtr("test")
		ptr2 := StringPtr("test")

		if ptr1 == ptr2 {
			t.Error("StringPtr returned same pointer for different calls")
		}
	})

	t.Run("IntPtr unique addresses", func(t *testing.T) {
		ptr1 := IntPtr(42)
		ptr2 := IntPtr(42)

		if ptr1 == ptr2 {
			t.Error("IntPtr returned same pointer for different calls")
		}
	})

	t.Run("BoolPtr unique addresses", func(t *testing.T) {
		ptr1 := BoolPtr(true)
		ptr2 := BoolPtr(true)

		if ptr1 == ptr2 {
			t.Error("BoolPtr returned same pointer for different calls")
		}
	})
}

func BenchmarkStringPtr(b *testing.B) {
	input := "benchmark test string"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringPtr(input)
	}
}

func BenchmarkIntPtr(b *testing.B) {
	input := 42
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntPtr(input)
	}
}

func BenchmarkBoolPtr(b *testing.B) {
	input := true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BoolPtr(input)
	}
}
