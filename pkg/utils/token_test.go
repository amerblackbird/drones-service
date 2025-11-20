package utils

import (
	"regexp"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateSalt(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "negative length", length: -1},
		{name: "small salt", length: 8},
		{name: "standard salt", length: 16},
		{name: "large salt", length: 32},
		{name: "very large salt", length: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSalt(tt.length)

			if tt.length <= 0 {
				if result != "" {
					t.Errorf("GenerateSalt(%d) should return empty string, got %q", tt.length, result)
				}
				return
			}

			// Salt should be hex encoded, so length should be 2x the input length
			expectedLen := tt.length * 2
			if len(result) != expectedLen {
				t.Errorf("GenerateSalt(%d) length = %d, want %d", tt.length, len(result), expectedLen)
			}

			// Verify it's valid hex
			hexRegex := regexp.MustCompile("^[0-9a-f]+$")
			if !hexRegex.MatchString(result) {
				t.Errorf("GenerateSalt(%d) = %q, not valid hex", tt.length, result)
			}
		})
	}
}

func TestGenerateSalt_Uniqueness(t *testing.T) {
	length := 16
	iterations := 100
	salts := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		salt := GenerateSalt(length)
		if salts[salt] {
			t.Errorf("GenerateSalt(%d) generated duplicate salt: %q", length, salt)
		}
		salts[salt] = true
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "negative length", length: -1},
		{name: "small token", length: 8},
		{name: "standard token", length: 32},
		{name: "large token", length: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateToken(tt.length)

			if tt.length <= 0 {
				if result != "" {
					t.Errorf("GenerateToken(%d) should return empty string, got %q", tt.length, result)
				}
				return
			}

			if len(result) != tt.length {
				t.Errorf("GenerateToken(%d) length = %d, want %d", tt.length, len(result), tt.length)
			}

			// Verify it only contains alphanumeric characters
			alphanumericRegex := regexp.MustCompile("^[A-Za-z0-9]+$")
			if !alphanumericRegex.MatchString(result) {
				t.Errorf("GenerateToken(%d) = %q, contains invalid characters", tt.length, result)
			}
		})
	}
}

func TestGenerateToken_Uniqueness(t *testing.T) {
	length := 32
	iterations := 100
	tokens := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		token := GenerateToken(length)
		if tokens[token] {
			t.Errorf("GenerateToken(%d) generated duplicate token: %q", length, token)
		}
		tokens[token] = true
	}
}

func TestGenerateOTP(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "negative length", length: -1},
		{name: "4 digit OTP", length: 4},
		{name: "6 digit OTP", length: 6},
		{name: "8 digit OTP", length: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateOTP(tt.length)

			if tt.length <= 0 {
				if result != "" {
					t.Errorf("GenerateOTP(%d) should return empty string, got %q", tt.length, result)
				}
				return
			}

			if len(result) != tt.length {
				t.Errorf("GenerateOTP(%d) length = %d, want %d", tt.length, len(result), tt.length)
			}

			// Verify it only contains digits
			digitRegex := regexp.MustCompile("^[0-9]+$")
			if !digitRegex.MatchString(result) {
				t.Errorf("GenerateOTP(%d) = %q, contains non-digit characters", tt.length, result)
			}
		})
	}
}

func TestGenerateOTP_LeadingZeros(t *testing.T) {
	length := 6
	iterations := 1000

	for i := 0; i < iterations; i++ {
		otp := GenerateOTP(length)
		if len(otp) != length {
			t.Errorf("GenerateOTP(%d) = %q, length mismatch (should pad with leading zeros)", length, otp)
		}
	}
}

func TestHashText(t *testing.T) {
	tests := []struct {
		name string
		text string
		salt string
	}{
		{name: "simple text", text: "password123", salt: "salt123"},
		{name: "empty text", text: "", salt: "salt123"},
		{name: "empty salt", text: "password123", salt: ""},
		{name: "special characters", text: "p@ssw0rd!", salt: "s@lt#123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashText(tt.text, tt.salt)

			// Verify hash is not empty
			if hash == "" {
				t.Error("HashText returned empty hash")
			}

			// Verify it's a valid bcrypt hash (starts with $2a$, $2b$, or $2y$)
			if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") && !strings.HasPrefix(hash, "$2y$") {
				t.Errorf("HashText returned invalid bcrypt hash: %q", hash)
			}

			// Verify hash length is reasonable (bcrypt hashes are typically 60 chars)
			if len(hash) < 50 {
				t.Errorf("HashText returned suspiciously short hash: %q", hash)
			}
		})
	}
}

func TestHashText_Consistency(t *testing.T) {
	text := "password123"
	salt := "salt456"

	hash1 := HashText(text, salt)
	hash2 := HashText(text, salt)

	// Bcrypt includes a random salt in the hash, so hashes should differ
	if hash1 == hash2 {
		t.Error("HashText should generate different hashes each time due to bcrypt's internal salt")
	}

	// But both should verify correctly
	if err := bcrypt.CompareHashAndPassword([]byte(hash1), []byte(text+salt)); err != nil {
		t.Error("First hash should verify correctly")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash2), []byte(text+salt)); err != nil {
		t.Error("Second hash should verify correctly")
	}
}

func TestVerifyOtpHash(t *testing.T) {
	otp := "1234"
	salt := "testsalt"
	hash := HashText(otp, salt)

	tests := []struct {
		name     string
		otp      string
		salt     string
		hash     string
		expected bool
	}{
		{name: "correct otp", otp: otp, salt: salt, hash: hash, expected: true},
		{name: "wrong otp", otp: "5678", salt: salt, hash: hash, expected: false},
		{name: "wrong salt", otp: otp, salt: "wrongsalt", hash: hash, expected: false},
		{name: "empty otp", otp: "", salt: salt, hash: hash, expected: false},
		{name: "empty salt", otp: otp, salt: "", hash: hash, expected: false},
		{name: "empty hash", otp: otp, salt: salt, hash: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyOtpHash(tt.otp, tt.salt, tt.hash)
			if result != tt.expected {
				t.Errorf("VerifyOtpHash(%q, %q, hash) = %v, want %v", tt.otp, tt.salt, result, tt.expected)
			}
		})
	}
}

func TestGenerateVerfiyCred(t *testing.T) {
	otpCode, hashedOtp, salt, token := GenerateVerfiyCred()

	// Verify OTP is 4 digits
	if len(otpCode) != 4 {
		t.Errorf("GenerateVerfiyCred() otp length = %d, want 4", len(otpCode))
	}
	digitRegex := regexp.MustCompile("^[0-9]+$")
	if !digitRegex.MatchString(otpCode) {
		t.Errorf("GenerateVerfiyCred() otp = %q, should be numeric", otpCode)
	}

	// Verify hashed OTP is not empty and looks like bcrypt
	if hashedOtp == "" {
		t.Error("GenerateVerfiyCred() hashedOtp is empty")
	}
	if !strings.HasPrefix(hashedOtp, "$2") {
		t.Errorf("GenerateVerfiyCred() hashedOtp = %q, doesn't look like bcrypt hash", hashedOtp)
	}

	// Verify salt is 32 chars (16 bytes hex encoded)
	if len(salt) != 32 {
		t.Errorf("GenerateVerfiyCred() salt length = %d, want 32", len(salt))
	}
	hexRegex := regexp.MustCompile("^[0-9a-f]+$")
	if !hexRegex.MatchString(salt) {
		t.Errorf("GenerateVerfiyCred() salt = %q, not valid hex", salt)
	}

	// Verify token is 32 chars alphanumeric
	if len(token) != 32 {
		t.Errorf("GenerateVerfiyCred() token length = %d, want 32", len(token))
	}
	alphanumericRegex := regexp.MustCompile("^[A-Za-z0-9]+$")
	if !alphanumericRegex.MatchString(token) {
		t.Errorf("GenerateVerfiyCred() token = %q, should be alphanumeric", token)
	}

	// Verify the hash can be verified with the OTP and salt
	if !VerifyOtpHash(otpCode, salt, hashedOtp) {
		t.Error("GenerateVerfiyCred() generated hash doesn't verify with otp and salt")
	}
}

func TestRandomPassword(t *testing.T) {
	password := RandomPassword()

	if len(password) != 16 {
		t.Errorf("RandomPassword() length = %d, want 16", len(password))
	}

	alphanumericRegex := regexp.MustCompile("^[A-Za-z0-9]+$")
	if !alphanumericRegex.MatchString(password) {
		t.Errorf("RandomPassword() = %q, should be alphanumeric", password)
	}
}

func TestRandomSalt(t *testing.T) {
	salt := RandomSalt()

	if len(salt) != 32 {
		t.Errorf("RandomSalt() length = %d, want 32", len(salt))
	}

	hexRegex := regexp.MustCompile("^[0-9a-f]+$")
	if !hexRegex.MatchString(salt) {
		t.Errorf("RandomSalt() = %q, not valid hex", salt)
	}
}

func BenchmarkGenerateSalt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSalt(16)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken(32)
	}
}

func BenchmarkGenerateOTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateOTP(6)
	}
}

func BenchmarkHashText(b *testing.B) {
	text := "password123"
	salt := "salt456"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashText(text, salt)
	}
}

func BenchmarkVerifyOtpHash(b *testing.B) {
	otp := "1234"
	salt := "testsalt"
	hash := HashText(otp, salt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyOtpHash(otp, salt, hash)
	}
}

func BenchmarkGenerateVerfiyCred(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateVerfiyCred()
	}
}
