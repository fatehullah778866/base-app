package auth

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// Password requirements
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

var (
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasNumber  = regexp.MustCompile(`[0-9]`)
	hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

type PasswordStrength int

const (
	PasswordWeak PasswordStrength = iota
	PasswordMedium
	PasswordStrong
)

type PasswordValidationResult struct {
	Valid   bool
	Strength PasswordStrength
	Errors  []string
}

func ValidatePassword(password string) PasswordValidationResult {
	result := PasswordValidationResult{
		Valid:  true,
		Errors: []string{},
	}

	// Length check
	if len(password) < MinPasswordLength {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must be at least 8 characters long")
	}
	if len(password) > MaxPasswordLength {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must be at most 128 characters long")
	}

	// Complexity checks
	hasUpperChar := hasUpper.MatchString(password)
	hasLowerChar := hasLower.MatchString(password)
	hasNumberChar := hasNumber.MatchString(password)
	hasSpecialChar := hasSpecial.MatchString(password)

	if !hasUpperChar {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one uppercase letter")
	}
	if !hasLowerChar {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one lowercase letter")
	}
	if !hasNumberChar {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one number")
	}
	if !hasSpecialChar {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one special character")
	}

	// Check for common weak passwords
	if isCommonPassword(password) {
		result.Valid = false
		result.Errors = append(result.Errors, "Password is too common. Please choose a stronger password")
	}

	// Determine strength
	if result.Valid {
		score := 0
		if len(password) >= 12 {
			score++
		}
		if hasUpperChar {
			score++
		}
		if hasLowerChar {
			score++
		}
		if hasNumberChar {
			score++
		}
		if hasSpecialChar {
			score++
		}
		// Check for mixed case and numbers
		if hasUpperChar && hasLowerChar && hasNumberChar {
			score++
		}
		// Check for special characters
		if hasSpecialChar {
			score++
		}

		if score >= 6 {
			result.Strength = PasswordStrong
		} else if score >= 4 {
			result.Strength = PasswordMedium
		} else {
			result.Strength = PasswordWeak
		}
	}

	return result
}

func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "12345678", "qwerty", "abc123", "password123",
		"admin", "letmein", "welcome", "monkey", "1234567890",
		"password1", "123456", "123456789", "1234567", "sunshine",
		"princess", "football", "iloveyou", "welcome123",
	}

	lowerPassword := password
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	// Check if password is all same character
	if len(password) > 0 {
		firstChar := password[0]
		allSame := true
		for _, char := range password {
			if char != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	// Check if password is sequential (12345678, abcdefgh)
	if isSequential(password) {
		return true
	}

	return false
}

func isSequential(s string) bool {
	if len(s) < 3 {
		return false
	}

	// Check numeric sequence
	isNumericSeq := true
	for i := 1; i < len(s); i++ {
		if !unicode.IsDigit(rune(s[i])) || !unicode.IsDigit(rune(s[i-1])) {
			isNumericSeq = false
			break
		}
		if int(s[i])-int(s[i-1]) != 1 {
			isNumericSeq = false
			break
		}
	}
	if isNumericSeq {
		return true
	}

	// Check reverse numeric sequence
	isReverseNumericSeq := true
	for i := 1; i < len(s); i++ {
		if !unicode.IsDigit(rune(s[i])) || !unicode.IsDigit(rune(s[i-1])) {
			isReverseNumericSeq = false
			break
		}
		if int(s[i-1])-int(s[i]) != 1 {
			isReverseNumericSeq = false
			break
		}
	}
	if isReverseNumericSeq {
		return true
	}

	return false
}

func HashPassword(password string) (string, error) {
	// Validate password before hashing
	validation := ValidatePassword(password)
	if !validation.Valid {
		return "", errors.New("password does not meet requirements: " + validation.Errors[0])
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
