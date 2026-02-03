package service_test

import (
	"context"
	"testing"

	"github.com/abhi10/ideahamster/internal/domain"
	"github.com/abhi10/ideahamster/internal/repository"
	"github.com/abhi10/ideahamster/internal/service"
)

func setupService() *service.IdeaService {
	ideaRepo := repository.NewMockIdeaRepository()
	voteRepo := repository.NewMockVoteRepository()
	return service.NewIdeaService(ideaRepo, voteRepo)
}

func TestIdeaService_GetIdea(t *testing.T) {
	svc := setupService()
	ctx := context.Background()

	t.Run("existing idea", func(t *testing.T) {
		idea, err := svc.GetIdea(ctx, "1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if idea == nil {
			t.Error("expected idea, got nil")
		}
		if idea.ID != "1" {
			t.Errorf("ID = %s, want 1", idea.ID)
		}
	})

	t.Run("non-existent idea", func(t *testing.T) {
		_, err := svc.GetIdea(ctx, "999")
		if err == nil {
			t.Error("expected error for non-existent idea")
		}
	})

	t.Run("empty id", func(t *testing.T) {
		_, err := svc.GetIdea(ctx, "")
		if err == nil {
			t.Error("expected error for empty id")
		}
		appErr := domain.GetAppError(err)
		if appErr.Type != domain.ErrorTypeValidation {
			t.Errorf("error type = %v, want Validation", appErr.Type)
		}
	})
}

func TestIdeaService_ListIdeas(t *testing.T) {
	svc := setupService()
	ctx := context.Background()

	t.Run("default options", func(t *testing.T) {
		ideas, err := svc.ListIdeas(ctx, domain.DefaultListOptions())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(ideas) == 0 {
			t.Error("expected ideas, got none")
		}
	})

	t.Run("with min votes filter", func(t *testing.T) {
		opts := domain.ListOptions{MinVotes: 50}
		ideas, err := svc.ListIdeas(ctx, opts)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for _, idea := range ideas {
			if idea.VoteCount < 50 {
				t.Errorf("idea %s has %d votes, want >= 50", idea.ID, idea.VoteCount)
			}
		}
	})
}

func TestIdeaService_ListEligibleIdeas(t *testing.T) {
	svc := setupService()
	ctx := context.Background()

	ideas, err := svc.ListEligibleIdeas(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// All returned ideas should meet vote threshold
	for _, idea := range ideas {
		if idea.VoteCount < 50 {
			t.Errorf("idea %s has %d votes, should be >= 50 for eligible", idea.ID, idea.VoteCount)
		}
	}
}

func TestIdeaService_Vote(t *testing.T) {
	ctx := context.Background()

	t.Run("successful vote", func(t *testing.T) {
		svc := setupService()
		err := svc.Vote(ctx, "1", "voter@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("duplicate vote", func(t *testing.T) {
		svc := setupService()
		email := "duplicate@example.com"

		// First vote
		err := svc.Vote(ctx, "1", email)
		if err != nil {
			t.Fatalf("first vote failed: %v", err)
		}

		// Second vote should fail
		err = svc.Vote(ctx, "1", email)
		if err == nil {
			t.Error("expected error on duplicate vote")
		}

		appErr := domain.GetAppError(err)
		if appErr.Type != domain.ErrorTypeConflict {
			t.Errorf("error type = %v, want Conflict", appErr.Type)
		}
	})

	t.Run("empty idea id", func(t *testing.T) {
		svc := setupService()
		err := svc.Vote(ctx, "", "voter@example.com")
		if err == nil {
			t.Error("expected error for empty idea id")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		svc := setupService()
		err := svc.Vote(ctx, "1", "")
		if err == nil {
			t.Error("expected error for empty email")
		}

		appErr := domain.GetAppError(err)
		if appErr.Type != domain.ErrorTypeUnauthorized {
			t.Errorf("error type = %v, want Unauthorized", appErr.Type)
		}
	})

	t.Run("non-existent idea", func(t *testing.T) {
		svc := setupService()
		err := svc.Vote(ctx, "999", "voter@example.com")
		if err == nil {
			t.Error("expected error for non-existent idea")
		}

		appErr := domain.GetAppError(err)
		if appErr.Type != domain.ErrorTypeNotFound {
			t.Errorf("error type = %v, want NotFound", appErr.Type)
		}
	})
}

func TestIdeaService_HasVoted(t *testing.T) {
	ctx := context.Background()
	svc := setupService()

	email := "checker@example.com"

	t.Run("before voting", func(t *testing.T) {
		hasVoted, err := svc.HasVoted(ctx, "1", email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if hasVoted {
			t.Error("HasVoted should return false before voting")
		}
	})

	// Cast a vote
	_ = svc.Vote(ctx, "1", email)

	t.Run("after voting", func(t *testing.T) {
		hasVoted, err := svc.HasVoted(ctx, "1", email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !hasVoted {
			t.Error("HasVoted should return true after voting")
		}
	})

	t.Run("different idea", func(t *testing.T) {
		hasVoted, err := svc.HasVoted(ctx, "2", email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if hasVoted {
			t.Error("HasVoted should return false for different idea")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		hasVoted, err := svc.HasVoted(ctx, "1", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if hasVoted {
			t.Error("HasVoted should return false for empty email")
		}
	})
}

func TestIdeaService_GetUserVotes(t *testing.T) {
	ctx := context.Background()
	svc := setupService()

	email := "multivote@example.com"

	// Cast votes on multiple ideas
	_ = svc.Vote(ctx, "1", email)
	_ = svc.Vote(ctx, "2", email)

	t.Run("user with votes", func(t *testing.T) {
		votes, err := svc.GetUserVotes(ctx, email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(votes) != 2 {
			t.Errorf("got %d votes, want 2", len(votes))
		}
	})

	t.Run("user without votes", func(t *testing.T) {
		votes, err := svc.GetUserVotes(ctx, "novotes@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(votes) != 0 {
			t.Errorf("got %d votes, want 0", len(votes))
		}
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := svc.GetUserVotes(ctx, "")
		if err == nil {
			t.Error("expected error for empty email")
		}
	})
}
