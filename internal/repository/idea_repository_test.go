package repository_test

import (
	"context"
	"testing"

	"github.com/abhi10/ideahamster/internal/domain"
	"github.com/abhi10/ideahamster/internal/repository"
)

func TestMockIdeaRepository_GetByID(t *testing.T) {
	repo := repository.NewMockIdeaRepository()
	ctx := context.Background()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing idea",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "non-existent idea",
			id:      "999",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idea, err := repo.GetByID(ctx, tt.id)

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

			if idea.ID != tt.id {
				t.Errorf("got ID %s, want %s", idea.ID, tt.id)
			}
		})
	}
}

func TestMockIdeaRepository_List(t *testing.T) {
	repo := repository.NewMockIdeaRepository()
	ctx := context.Background()

	tests := []struct {
		name     string
		opts     domain.ListOptions
		minCount int
		maxCount int
	}{
		{
			name:     "list all",
			opts:     domain.DefaultListOptions(),
			minCount: 5,
			maxCount: 10,
		},
		{
			name: "filter eligible (50+ votes)",
			opts: domain.ListOptions{
				MinVotes: 50,
			},
			minCount: 2, // Ideas with 67 and 52 votes
			maxCount: 2,
		},
		{
			name: "filter by category",
			opts: domain.ListOptions{
				Category: domain.CategoryFrontend,
			},
			minCount: 1,
			maxCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ideas, err := repo.List(ctx, tt.opts)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			count := len(ideas)
			if count < tt.minCount || count > tt.maxCount {
				t.Errorf("got %d ideas, want between %d and %d", count, tt.minCount, tt.maxCount)
			}
		})
	}
}

func TestMockIdeaRepository_IncrementVoteCount(t *testing.T) {
	repo := repository.NewMockIdeaRepository()
	ctx := context.Background()

	// Get initial vote count
	idea, err := repo.GetByID(ctx, "1")
	if err != nil {
		t.Fatalf("failed to get idea: %v", err)
	}
	initialCount := idea.VoteCount

	// Increment vote
	err = repo.IncrementVoteCount(ctx, "1")
	if err != nil {
		t.Errorf("failed to increment vote: %v", err)
	}

	// Verify increment
	idea, _ = repo.GetByID(ctx, "1")
	if idea.VoteCount != initialCount+1 {
		t.Errorf("vote count = %d, want %d", idea.VoteCount, initialCount+1)
	}
}

func TestMockIdeaRepository_IncrementVoteCount_NotFound(t *testing.T) {
	repo := repository.NewMockIdeaRepository()
	ctx := context.Background()

	err := repo.IncrementVoteCount(ctx, "999")
	if err == nil {
		t.Error("expected error for non-existent idea")
	}

	// Verify it's a NotFound error
	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Error("expected AppError type")
		return
	}
	if appErr.Type != domain.ErrorTypeNotFound {
		t.Errorf("got error type %v, want NotFound", appErr.Type)
	}
}
