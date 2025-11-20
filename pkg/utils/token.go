package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const lettersAndDigits = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateSalt(length int) string {
	if length <= 0 {
		return ""
	}

	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		// fallback: all zeros if crypto random fails
		for i := range salt {
			salt[i] = 0
		}
	}

	// return as hex string
	return hex.EncodeToString(salt)
}

// GeneratePassword generates a random alphanumeric password of given length
func GenerateToken(length int) string {
	if length <= 0 {
		return ""
	}

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersAndDigits))))
		if err != nil {
			// fallback: use '0' if crypto random fails
			password[i] = '0'
			continue
		}
		password[i] = lettersAndDigits[index.Int64()]
	}

	return string(password)
}

// GenerateOTP generates a secure numeric OTP of given length (always returns string).
func GenerateOTP(length int) string {
	if length <= 0 {
		return ""
	}

	max := big.NewInt(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10)) // 10^length
	}

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// fallback: return all zeros if random fails
		return fmt.Sprintf("%0*d", length, 0)
	}

	// pad with leading zeros
	return fmt.Sprintf("%0*d", length, n.Int64())
}

func HashText(text, salt string) string {
	password := []byte(text + salt)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func VerifyOtpHash(otp, salt, hash string) bool {
	hashToCompare := []byte(hash)
	password := []byte(otp + salt)
	// Verify password
	err := bcrypt.CompareHashAndPassword(hashToCompare, password)
	return err == nil
}

func GenerateVerfiyCred() (string, string, string, string) {
	otpCode := GenerateOTP(4)
	salt := GenerateSalt(16)
	hashedOtp := HashText(otpCode, salt)
	token := GenerateToken(32)
	return otpCode, hashedOtp, salt, token
}


func RandomPassword() string {
	return GenerateToken(16)
}

func RandomSalt() string {
	return GenerateSalt(16)
}