// Package service contains business logic that orchestrates domain operations.
// Services depend on repository interfaces, not concrete implementations.
package service

import (
	"context"
	"log"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/abhi10/ideahamster/internal/domain"
)

// IdeaService handles business logic for ideas
type IdeaService struct {
	ideaRepo domain.IdeaRepository
	voteRepo domain.VoteRepository
}

// NewIdeaService creates a new IdeaService with injected dependencies
func NewIdeaService(ideaRepo domain.IdeaRepository, voteRepo domain.VoteRepository) *IdeaService {
	return &IdeaService{
		ideaRepo: ideaRepo,
		voteRepo: voteRepo,
	}
}

// GetIdea retrieves an idea by ID
func (s *IdeaService) GetIdea(ctx context.Context, id string) (*domain.Idea, error) {
	if id == "" {
		return nil, domain.NewValidationError("Idea ID is required")
	}
	return s.ideaRepo.GetByID(ctx, id)
}

// ListIdeas retrieves ideas with optional filtering
func (s *IdeaService) ListIdeas(ctx context.Context, opts domain.ListOptions) ([]domain.Idea, error) {
	return s.ideaRepo.List(ctx, opts)
}

// ListEligibleIdeas retrieves ideas that have reached the vote threshold
func (s *IdeaService) ListEligibleIdeas(ctx context.Context) ([]domain.Idea, error) {
	opts := domain.DefaultListOptions()
	opts.MinVotes = config.VoteThreshold
	return s.ideaRepo.List(ctx, opts)
}

// Vote records a vote for an idea
func (s *IdeaService) Vote(ctx context.Context, ideaID, email string) error {
	if ideaID == "" {
		return domain.NewValidationError("Idea ID is required")
	}
	if email == "" {
		return domain.NewUnauthorizedError("Email verification required")
	}

	// Check if idea exists
	idea, err := s.ideaRepo.GetByID(ctx, ideaID)
	if err != nil {
		return err
	}

	// Check if user already voted
	existingVote, err := s.voteRepo.GetByIdeaAndEmail(ctx, ideaID, email)
	if err == nil && existingVote != nil {
		return domain.NewConflictError("You have already voted on this idea")
	}

	// Record the vote
	vote := &domain.Vote{
		IdeaID: ideaID,
		Email:  email,
	}
	if err := s.voteRepo.Create(ctx, vote); err != nil {
		return err
	}

	// Increment vote count
	if err := s.ideaRepo.IncrementVoteCount(ctx, ideaID); err != nil {
		log.Printf("Failed to increment vote count for idea %s: %v", ideaID, err)
		// Don't return error - vote was recorded
	}

	// Check if idea should be marked as eligible
	if idea.VoteCount+1 >= config.VoteThreshold && idea.Status == domain.StatusPending {
		idea.Status = domain.StatusEligible
		if err := s.ideaRepo.Update(ctx, idea); err != nil {
			log.Printf("Failed to update idea status to eligible: %v", err)
		}
	}

	return nil
}

// HasVoted checks if a user has already voted on an idea
func (s *IdeaService) HasVoted(ctx context.Context, ideaID, email string) (bool, error) {
	if email == "" {
		return false, nil
	}

	vote, err := s.voteRepo.GetByIdeaAndEmail(ctx, ideaID, email)
	if err != nil {
		// Not found is not an error for this check
		if appErr, ok := err.(*domain.AppError); ok && appErr.Type == domain.ErrorTypeNotFound {
			return false, nil
		}
		return false, err
	}
	return vote != nil, nil
}

// GetUserVotes returns all ideas a user has voted on
func (s *IdeaService) GetUserVotes(ctx context.Context, email string) ([]domain.Vote, error) {
	if email == "" {
		return nil, domain.NewUnauthorizedError("Email required")
	}
	return s.voteRepo.ListByEmail(ctx, email)
}
