package domain_test

import (
	"errors"
	"testing"

	"github.com/abhi10/ideahamster/internal/domain"
)

func TestAppError_HTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		errType  domain.ErrorType
		wantCode int
	}{
		{"internal", domain.ErrorTypeInternal, 500},
		{"validation", domain.ErrorTypeValidation, 400},
		{"not found", domain.ErrorTypeNotFound, 404},
		{"unauthorized", domain.ErrorTypeUnauthorized, 401},
		{"forbidden", domain.ErrorTypeForbidden, 403},
		{"conflict", domain.ErrorTypeConflict, 409},
		{"rate limit", domain.ErrorTypeRateLimit, 429},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := &domain.AppError{Type: tc.errType}
			if got := err.HTTPStatus(); got != tc.wantCode {
				t.Errorf("HTTPStatus() = %d, want %d", got, tc.wantCode)
			}
		})
	}
}

func TestAppError_UserMessage(t *testing.T) {
	err := &domain.AppError{
		Type:     domain.ErrorTypeInternal,
		Message:  "Something went wrong",
		Internal: errors.New("database connection failed"),
	}

	// UserMessage should return the safe message
	if got := err.UserMessage(); got != "Something went wrong" {
		t.Errorf("UserMessage() = %q, want %q", got, "Something went wrong")
	}

	// Error() should include internal details for logging
	fullErr := err.Error()
	if fullErr == "Something went wrong" {
		t.Error("Error() should include internal details")
	}
}

func TestAppError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &domain.AppError{
		Type:     domain.ErrorTypeInternal,
		Message:  "outer error",
		Internal: innerErr,
	}

	// Should be able to unwrap to inner error
	if !errors.Is(err, innerErr) {
		t.Error("errors.Is should find inner error")
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		err      *domain.AppError
		wantType domain.ErrorType
	}{
		{
			name:     "NewInternalError",
			err:      domain.NewInternalError("internal", errors.New("db error")),
			wantType: domain.ErrorTypeInternal,
		},
		{
			name:     "NewValidationError",
			err:      domain.NewValidationError("invalid input"),
			wantType: domain.ErrorTypeValidation,
		},
		{
			name:     "NewNotFoundError",
			err:      domain.NewNotFoundError("User"),
			wantType: domain.ErrorTypeNotFound,
		},
		{
			name:     "NewUnauthorizedError",
			err:      domain.NewUnauthorizedError("login required"),
			wantType: domain.ErrorTypeUnauthorized,
		},
		{
			name:     "NewForbiddenError",
			err:      domain.NewForbiddenError("access denied"),
			wantType: domain.ErrorTypeForbidden,
		},
		{
			name:     "NewConflictError",
			err:      domain.NewConflictError("already exists"),
			wantType: domain.ErrorTypeConflict,
		},
		{
			name:     "NewRateLimitError",
			err:      domain.NewRateLimitError(),
			wantType: domain.ErrorTypeRateLimit,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Type != tc.wantType {
				t.Errorf("Type = %v, want %v", tc.err.Type, tc.wantType)
			}
		})
	}
}

func TestGetAppError(t *testing.T) {
	t.Run("with AppError", func(t *testing.T) {
		original := domain.NewValidationError("test")
		result := domain.GetAppError(original)
		if result.Type != domain.ErrorTypeValidation {
			t.Errorf("Type = %v, want Validation", result.Type)
		}
	})

	t.Run("with regular error", func(t *testing.T) {
		original := errors.New("regular error")
		result := domain.GetAppError(original)
		if result.Type != domain.ErrorTypeInternal {
			t.Errorf("Type = %v, want Internal", result.Type)
		}
	})
}

func TestIsAppError(t *testing.T) {
	t.Run("with AppError", func(t *testing.T) {
		err := domain.NewValidationError("test")
		if !domain.IsAppError(err) {
			t.Error("IsAppError should return true for AppError")
		}
	})

	t.Run("with regular error", func(t *testing.T) {
		err := errors.New("regular error")
		if domain.IsAppError(err) {
			t.Error("IsAppError should return false for regular error")
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("wrap nil", func(t *testing.T) {
		result := domain.Wrap(nil, "message")
		if result != nil {
			t.Error("Wrap(nil) should return nil")
		}
	})

	t.Run("wrap AppError", func(t *testing.T) {
		original := domain.NewValidationError("original")
		wrapped := domain.Wrap(original, "wrapped")

		appErr, ok := wrapped.(*domain.AppError)
		if !ok {
			t.Fatal("wrapped error should be AppError")
		}

		// Should preserve type
		if appErr.Type != domain.ErrorTypeValidation {
			t.Errorf("Type = %v, want Validation", appErr.Type)
		}

		// Should have new message
		if appErr.Message != "wrapped" {
			t.Errorf("Message = %q, want %q", appErr.Message, "wrapped")
		}
	})

	t.Run("wrap regular error", func(t *testing.T) {
		original := errors.New("original")
		wrapped := domain.Wrap(original, "wrapped")

		appErr, ok := wrapped.(*domain.AppError)
		if !ok {
			t.Fatal("wrapped error should be AppError")
		}

		// Should default to internal
		if appErr.Type != domain.ErrorTypeInternal {
			t.Errorf("Type = %v, want Internal", appErr.Type)
		}
	})
}
