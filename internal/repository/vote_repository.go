package repository

import (
	"context"
	"sync"
	"time"

	"github.com/abhi10/ideahamster/internal/domain"
)

// MockVoteRepository is an in-memory implementation for testing/Phase 1
type MockVoteRepository struct {
	votes []domain.Vote
	mu    sync.RWMutex
}

// NewMockVoteRepository creates a new mock vote repository
func NewMockVoteRepository() *MockVoteRepository {
	return &MockVoteRepository{
		votes: make([]domain.Vote, 0),
	}
}

// GetByIdeaAndEmail checks if a user already voted on an idea
func (r *MockVoteRepository) GetByIdeaAndEmail(ctx context.Context, ideaID, email string) (*domain.Vote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.votes {
		if r.votes[i].IdeaID == ideaID && r.votes[i].Email == email {
			vote := r.votes[i]
			return &vote, nil
		}
	}
	return nil, domain.NewNotFoundError("Vote")
}

// Create records a new vote
func (r *MockVoteRepository) Create(ctx context.Context, vote *domain.Vote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate
	for _, v := range r.votes {
		if v.IdeaID == vote.IdeaID && v.Email == vote.Email {
			return domain.NewConflictError("You have already voted on this idea")
		}
	}

	vote.CreatedAt = time.Now()
	r.votes = append(r.votes, *vote)
	return nil
}

// CountByIdea returns the vote count for an idea
func (r *MockVoteRepository) CountByIdea(ctx context.Context, ideaID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, v := range r.votes {
		if v.IdeaID == ideaID {
			count++
		}
	}
	return count, nil
}

// ListByEmail returns all votes by a user
func (r *MockVoteRepository) ListByEmail(ctx context.Context, email string) ([]domain.Vote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Vote, 0)
	for _, v := range r.votes {
		if v.Email == email {
			result = append(result, v)
		}
	}
	return result, nil
}

// Ensure MockVoteRepository implements VoteRepository
var _ domain.VoteRepository = (*MockVoteRepository)(nil)
