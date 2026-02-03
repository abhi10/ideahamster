package sanitizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/microcosm-cc/bluemonday"
)

var (
	// strictPolicy strips ALL HTML
	strictPolicy *bluemonday.Policy

	// emailRegex for validation
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// spaceNormalizer collapses multiple spaces
	spaceNormalizer = regexp.MustCompile(`\s+`)

	// disposableDomains to block
	disposableDomains = map[string]bool{
		"tempmail.com":      true,
		"throwaway.email":   true,
		"guerrillamail.com": true,
		"mailinator.com":    true,
		"10minutemail.com":  true,
		"temp-mail.org":     true,
		"fakeinbox.com":     true,
		"trashmail.com":     true,
		"yopmail.com":       true,
		"getnada.com":       true,
		"maildrop.cc":       true,
		"dispostable.com":   true,
		"mailnesia.com":     true,
		"tempail.com":       true,
		"mohmal.com":        true,
		"emailondeck.com":   true,
		"tempr.email":       true,
		"discard.email":     true,
		"tmpmail.org":       true,
		"tmpmail.net":       true,
	}

	// profanityWords list (extend as needed)
	profanityWords = []string{
		// Add profanity words here
	}
)

func init() {
	strictPolicy = bluemonday.StrictPolicy()
}

// SanitizeText strips all HTML and normalizes whitespace
func SanitizeText(input string) string {
	clean := strictPolicy.Sanitize(input)
	clean = strings.TrimSpace(clean)
	clean = spaceNormalizer.ReplaceAllString(clean, " ")
	return clean
}

// SanitizeTitle sanitizes and validates idea titles
func SanitizeTitle(title string) (string, error) {
	clean := SanitizeText(title)

	if len(clean) < config.TitleMinLength {
		return "", ErrTitleTooShort
	}

	if len(clean) > config.TitleMaxLength {
		return "", ErrTitleTooLong
	}

	if ContainsProfanity(clean) {
		return "", ErrContainsProfanity
	}

	return clean, nil
}

// SanitizeDescription sanitizes and validates idea descriptions
func SanitizeDescription(desc string) (string, error) {
	clean := SanitizeText(desc)

	if len(clean) < config.DescriptionMinLength {
		return "", ErrDescriptionTooShort
	}

	if len(clean) > config.DescriptionMaxLength {
		return "", ErrDescriptionTooLong
	}

	if ContainsProfanity(clean) {
		return "", ErrContainsProfanity
	}

	return clean, nil
}

// ValidateEmail checks if an email is valid and not disposable
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return ErrEmailRequired
	}

	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}

	// Check for disposable email domains
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrEmailInvalid
	}

	if disposableDomains[parts[1]] {
		return ErrDisposableEmail
	}

	return nil
}

// SanitizeEmail cleans and validates an email address
func SanitizeEmail(email string) (string, error) {
	clean := strings.TrimSpace(strings.ToLower(email))

	if err := ValidateEmail(clean); err != nil {
		return "", err
	}

	return clean, nil
}

// ContainsProfanity checks if text contains profanity
func ContainsProfanity(text string) bool {
	lower := strings.ToLower(text)
	for _, word := range profanityWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

// SanitizeCategory validates category selection
func SanitizeCategory(category string) (string, error) {
	category = strings.TrimSpace(category)

	validCategories := map[string]bool{
		"Frontend":   true,
		"Backend":    true,
		"Full Stack": true,
	}

	if !validCategories[category] {
		return "", ErrInvalidCategory
	}

	return category, nil
}

// SanitizeTags cleans and validates tags
func SanitizeTags(tags []string) ([]string, error) {
	if len(tags) > config.MaxTags {
		return nil, ErrTooManyTags
	}

	cleaned := make([]string, 0, len(tags))
	seen := make(map[string]bool)

	for _, tag := range tags {
		tag = SanitizeText(tag)

		if tag == "" {
			continue
		}

		// Limit tag length
		if len(tag) > config.TagMaxLength {
			tag = tag[:config.TagMaxLength]
		}

		// Remove duplicates
		lower := strings.ToLower(tag)
		if seen[lower] {
			continue
		}
		seen[lower] = true

		// Skip profane tags silently
		if ContainsProfanity(tag) {
			continue
		}

		cleaned = append(cleaned, tag)
	}

	return cleaned, nil
}

// SanitizeError represents a validation error
type SanitizeError struct {
	Message string
}

func (e *SanitizeError) Error() string {
	return e.Message
}

// Validation errors
var (
	ErrTitleTooShort = &SanitizeError{
		Message: fmt.Sprintf("Title must be at least %d characters", config.TitleMinLength),
	}
	ErrTitleTooLong = &SanitizeError{
		Message: fmt.Sprintf("Title must be less than %d characters", config.TitleMaxLength),
	}
	ErrDescriptionTooShort = &SanitizeError{
		Message: fmt.Sprintf("Description must be at least %d characters", config.DescriptionMinLength),
	}
	ErrDescriptionTooLong = &SanitizeError{
		Message: fmt.Sprintf("Description must be less than %d characters", config.DescriptionMaxLength),
	}
	ErrContainsProfanity = &SanitizeError{Message: "Content contains inappropriate language"}
	ErrEmailRequired     = &SanitizeError{Message: "Email is required"}
	ErrEmailInvalid      = &SanitizeError{Message: "Invalid email format"}
	ErrDisposableEmail   = &SanitizeError{Message: "Disposable email addresses are not allowed"}
	ErrInvalidCategory   = &SanitizeError{Message: "Invalid category"}
	ErrTooManyTags       = &SanitizeError{
		Message: fmt.Sprintf("Maximum %d tags allowed", config.MaxTags),
	}
)
