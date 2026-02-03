package repository_test

import (
	"context"
	"testing"

	"github.com/abhi10/ideahamster/internal/domain"
	"github.com/abhi10/ideahamster/internal/repository"
)

func TestMockVoteRepository_Create(t *testing.T) {
	repo := repository.NewMockVoteRepository()
	ctx := context.Background()

	vote := &domain.Vote{
		IdeaID: "1",
		Email:  "test@example.com",
	}

	// First vote should succeed
	err := repo.Create(ctx, vote)
	if err != nil {
		t.Errorf("unexpected error on first vote: %v", err)
	}

	// Duplicate vote should fail
	err = repo.Create(ctx, vote)
	if err == nil {
		t.Error("expected error on duplicate vote")
	}

	// Verify it's a conflict error
	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Error("expected AppError type")
		return
	}
	if appErr.Type != domain.ErrorTypeConflict {
		t.Errorf("got error type %v, want Conflict", appErr.Type)
	}
}

func TestMockVoteRepository_GetByIdeaAndEmail(t *testing.T) {
	repo := repository.NewMockVoteRepository()
	ctx := context.Background()

	// Create a vote first
	vote := &domain.Vote{
		IdeaID: "1",
		Email:  "test@example.com",
	}
	_ = repo.Create(ctx, vote)

	tests := []struct {
		name    string
		ideaID  string
		email   string
		wantErr bool
	}{
		{
			name:    "existing vote",
			ideaID:  "1",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "non-existent vote",
			ideaID:  "1",
			email:   "other@example.com",
			wantErr: true,
		},
		{
			name:    "different idea",
			ideaID:  "2",
			email:   "test@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.GetByIdeaAndEmail(ctx, tt.ideaID, tt.email)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if found.IdeaID != tt.ideaID || found.Email != tt.email {
				t.Errorf("got vote for idea=%s email=%s, want idea=%s email=%s",
					found.IdeaID, found.Email, tt.ideaID, tt.email)
			}
		})
	}
}

func TestMockVoteRepository_CountByIdea(t *testing.T) {
	repo := repository.NewMockVoteRepository()
	ctx := context.Background()

	// Initially no votes
	count, err := repo.CountByIdea(ctx, "1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("initial count = %d, want 0", count)
	}

	// Add votes
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "1", Email: "a@test.com"})
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "1", Email: "b@test.com"})
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "2", Email: "a@test.com"}) // Different idea

	count, _ = repo.CountByIdea(ctx, "1")
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestMockVoteRepository_ListByEmail(t *testing.T) {
	repo := repository.NewMockVoteRepository()
	ctx := context.Background()

	email := "voter@test.com"

	// Add votes for same email on different ideas
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "1", Email: email})
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "2", Email: email})
	_ = repo.Create(ctx, &domain.Vote{IdeaID: "3", Email: "other@test.com"})

	votes, err := repo.ListByEmail(ctx, email)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(votes) != 2 {
		t.Errorf("got %d votes, want 2", len(votes))
	}

	for _, v := range votes {
		if v.Email != email {
			t.Errorf("got vote with email %s, want %s", v.Email, email)
		}
	}
}
