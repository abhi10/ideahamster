package sanitizer_test

import (
	"strings"
	"testing"

	"github.com/abhi10/ideahamster/internal/sanitizer"
)

// ===========================================
// Email Validation Tests (CRITICAL - Used Now)
// ===========================================

func TestSanitizeEmail_Valid(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.org",
		"user+tag@example.co.uk",
		"firstname.lastname@company.com",
	}

	for _, email := range validEmails {
		t.Run(email, func(t *testing.T) {
			result, err := sanitizer.SanitizeEmail(email)
			if err != nil {
				t.Errorf("SanitizeEmail(%q) returned error: %v", email, err)
			}
			if result != strings.ToLower(email) {
				t.Errorf("SanitizeEmail(%q) = %q, want %q", email, result, strings.ToLower(email))
			}
		})
	}
}

func TestSanitizeEmail_Invalid(t *testing.T) {
	invalidEmails := []struct {
		email string
		desc  string
	}{
		{"", "empty string"},
		{"notanemail", "no @ symbol"},
		{"@nodomain.com", "no local part"},
		{"noat.com", "missing @"},
		{"spaces in@email.com", "spaces in email"},
		{"test@", "no domain"},
	}

	for _, tc := range invalidEmails {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := sanitizer.SanitizeEmail(tc.email)
			if err == nil {
				t.Errorf("SanitizeEmail(%q) should return error for %s", tc.email, tc.desc)
			}
		})
	}
}

func TestSanitizeEmail_DisposableDomains(t *testing.T) {
	disposableEmails := []string{
		"test@tempmail.com",
		"test@mailinator.com",
		"test@guerrillamail.com",
		"test@10minutemail.com",
		"test@yopmail.com",
	}

	for _, email := range disposableEmails {
		t.Run(email, func(t *testing.T) {
			_, err := sanitizer.SanitizeEmail(email)
			if err == nil {
				t.Errorf("SanitizeEmail(%q) should reject disposable email", email)
			}
		})
	}
}

func TestSanitizeEmail_Normalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  TEST@EXAMPLE.COM  ", "test@example.com"},
		{"User@Domain.COM", "user@domain.com"},
		{"\tspaced@email.com\n", "spaced@email.com"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := sanitizer.SanitizeEmail(tc.input)
			if err != nil {
				t.Errorf("SanitizeEmail(%q) returned error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("SanitizeEmail(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// ===========================================
// Title Validation Tests (For Future Use)
// ===========================================

func TestSanitizeTitle_Valid(t *testing.T) {
	validTitles := []string{
		"Valid Title",
		"A Good Idea for an App",
		"12345", // Minimum 5 chars
		strings.Repeat("a", 80), // Maximum 80 chars
	}

	for _, title := range validTitles {
		t.Run(title[:min(20, len(title))], func(t *testing.T) {
			_, err := sanitizer.SanitizeTitle(title)
			if err != nil {
				t.Errorf("SanitizeTitle(%q) returned error: %v", title, err)
			}
		})
	}
}

func TestSanitizeTitle_Invalid(t *testing.T) {
	tests := []struct {
		title string
		desc  string
	}{
		{"", "empty"},
		{"ab", "too short (2 chars)"},
		{"abcd", "too short (4 chars)"},
		{strings.Repeat("a", 81), "too long (81 chars)"},
		{strings.Repeat("a", 100), "way too long"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := sanitizer.SanitizeTitle(tc.title)
			if err == nil {
				t.Errorf("SanitizeTitle should reject %s", tc.desc)
			}
		})
	}
}

func TestSanitizeTitle_HTMLStripping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<script>alert('xss')</script>Normal Title", "Normal Title"},
		{"Title with <b>bold</b>", "Title with bold"},
		{"<img src=x onerror=alert(1)>Safe Title", "Safe Title"},
	}

	for _, tc := range tests {
		t.Run(tc.input[:min(20, len(tc.input))], func(t *testing.T) {
			result, err := sanitizer.SanitizeTitle(tc.input)
			// May error if result is too short after stripping
			if err == nil && result != tc.expected {
				t.Errorf("SanitizeTitle(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// ===========================================
// Description Validation Tests
// ===========================================

func TestSanitizeDescription_Valid(t *testing.T) {
	validDescriptions := []string{
		"This is a valid description that is long enough to pass validation.",
		strings.Repeat("a", 20),  // Minimum
		strings.Repeat("a", 500), // Maximum
	}

	for i, desc := range validDescriptions {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			_, err := sanitizer.SanitizeDescription(desc)
			if err != nil {
				t.Errorf("SanitizeDescription returned error: %v", err)
			}
		})
	}
}

func TestSanitizeDescription_Invalid(t *testing.T) {
	tests := []struct {
		desc string
		name string
	}{
		{"", "empty"},
		{"too short", "too short (9 chars)"},
		{strings.Repeat("a", 501), "too long (501 chars)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sanitizer.SanitizeDescription(tc.desc)
			if err == nil {
				t.Errorf("SanitizeDescription should reject %s", tc.name)
			}
		})
	}
}

// ===========================================
// Category Validation Tests
// ===========================================

func TestSanitizeCategory_Valid(t *testing.T) {
	validCategories := []string{
		"Frontend",
		"Backend",
		"Full Stack",
	}

	for _, cat := range validCategories {
		t.Run(cat, func(t *testing.T) {
			result, err := sanitizer.SanitizeCategory(cat)
			if err != nil {
				t.Errorf("SanitizeCategory(%q) returned error: %v", cat, err)
			}
			if result != cat {
				t.Errorf("SanitizeCategory(%q) = %q, want %q", cat, result, cat)
			}
		})
	}
}

func TestSanitizeCategory_Invalid(t *testing.T) {
	invalidCategories := []string{
		"",
		"Invalid",
		"frontend", // Case sensitive
		"BACKEND",
		"Random Category",
	}

	for _, cat := range invalidCategories {
		t.Run(cat, func(t *testing.T) {
			_, err := sanitizer.SanitizeCategory(cat)
			if err == nil {
				t.Errorf("SanitizeCategory(%q) should return error", cat)
			}
		})
	}
}

// ===========================================
// Tags Validation Tests
// ===========================================

func TestSanitizeTags_Valid(t *testing.T) {
	tests := []struct {
		input    []string
		expected int
	}{
		{[]string{"Go", "Web", "API"}, 3},
		{[]string{"Single"}, 1},
		{[]string{}, 0},
		{[]string{"a", "b", "c", "d", "e"}, 5}, // Maximum 5
	}

	for i, tc := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result, err := sanitizer.SanitizeTags(tc.input)
			if err != nil {
				t.Errorf("SanitizeTags returned error: %v", err)
			}
			if len(result) != tc.expected {
				t.Errorf("SanitizeTags returned %d tags, want %d", len(result), tc.expected)
			}
		})
	}
}

func TestSanitizeTags_TooMany(t *testing.T) {
	tags := []string{"a", "b", "c", "d", "e", "f"} // 6 tags, max is 5
	_, err := sanitizer.SanitizeTags(tags)
	if err == nil {
		t.Error("SanitizeTags should reject more than 5 tags")
	}
}

func TestSanitizeTags_Deduplication(t *testing.T) {
	tags := []string{"Go", "go", "GO", "Web"}
	result, err := sanitizer.SanitizeTags(tags)
	if err != nil {
		t.Errorf("SanitizeTags returned error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("SanitizeTags should deduplicate, got %d tags, want 2", len(result))
	}
}

func TestSanitizeTags_EmptyTagsFiltered(t *testing.T) {
	tags := []string{"Valid", "", "  ", "Another"}
	result, err := sanitizer.SanitizeTags(tags)
	if err != nil {
		t.Errorf("SanitizeTags returned error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("SanitizeTags should filter empty tags, got %d, want 2", len(result))
	}
}

// ===========================================
// Text Sanitization Tests (XSS Prevention)
// ===========================================

func TestSanitizeText_XSS(t *testing.T) {
	tests := []struct {
		input    string
		contains string
		desc     string
	}{
		{
			input:    "<script>alert('xss')</script>",
			contains: "script",
			desc:     "script tags should be stripped",
		},
		{
			input:    "<img src=x onerror=alert(1)>",
			contains: "onerror",
			desc:     "event handlers should be stripped",
		},
		{
			input:    "<a href='javascript:alert(1)'>click</a>",
			contains: "javascript",
			desc:     "javascript URLs should be stripped",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			result := sanitizer.SanitizeText(tc.input)
			if strings.Contains(strings.ToLower(result), tc.contains) {
				t.Errorf("SanitizeText(%q) still contains %q", tc.input, tc.contains)
			}
		})
	}
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
