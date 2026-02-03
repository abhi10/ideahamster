// Package repository provides data access implementations.
// Repositories implement domain interfaces and handle database operations.
package repository

import (
	"context"
	"sync"

	"github.com/abhi10/ideahamster/internal/domain"
)

// ===========================================
// Mock Implementation (for Phase 1)
// Replace with PostgreSQL implementation when DB is connected
// ===========================================

// MockIdeaRepository is an in-memory implementation for testing/Phase 1
type MockIdeaRepository struct {
	ideas []domain.Idea
	mu    sync.RWMutex
}

// NewMockIdeaRepository creates a new mock repository with sample data
func NewMockIdeaRepository() *MockIdeaRepository {
	return &MockIdeaRepository{
		ideas: []domain.Idea{
			{
				ID:          "1",
				Title:       "AI-Powered Recipe Generator",
				Description: "Generate custom recipes based on ingredients you have at home. Uses AI to suggest creative combinations and cooking instructions. Perfect for reducing food waste and discovering new meals!",
				Category:    domain.CategoryFullStack,
				Tags:        []string{"AI", "Food", "Mobile"},
				Submitter:   "Chef Mike",
				VoteCount:   67,
				Status:      domain.StatusEligible,
			},
			{
				ID:          "2",
				Title:       "Retro Pomodoro Timer",
				Description: "A productivity timer with 90s aesthetics. Track your work sessions with pixel art animations and chiptune sounds. Includes stats tracking and customizable work/break intervals.",
				Category:    domain.CategoryFrontend,
				Tags:        []string{"Productivity", "Retro", "Timer"},
				Submitter:   "ProductivityGuru",
				VoteCount:   52,
				Status:      domain.StatusEligible,
			},
			{
				ID:          "3",
				Title:       "Local Event Discovery API",
				Description: "API that aggregates events from multiple sources and provides a unified interface for discovering local happenings. Supports filtering by date, category, and location.",
				Category:    domain.CategoryBackend,
				Tags:        []string{"API", "Events", "Location"},
				Submitter:   "DevCommunity",
				VoteCount:   45,
				Status:      domain.StatusPending,
			},
			{
				ID:          "4",
				Title:       "Habit Tracker with Friends",
				Description: "Track daily habits and compete with friends. See who can maintain the longest streak. Gamified habit building with rewards and achievements.",
				Category:    domain.CategoryFullStack,
				Tags:        []string{"Habits", "Social", "Gamification"},
				Submitter:   "LifeHacker",
				VoteCount:   38,
				Status:      domain.StatusPending,
			},
			{
				ID:          "5",
				Title:       "Markdown-Based Wiki",
				Description: "Simple wiki system powered by Markdown files. Perfect for personal knowledge bases or team documentation. Supports full-text search and automatic linking.",
				Category:    domain.CategoryFullStack,
				Tags:        []string{"Documentation", "Markdown", "Knowledge"},
				Submitter:   "",
				VoteCount:   31,
				Status:      domain.StatusPending,
			},
		},
	}
}

// GetByID retrieves an idea by ID
func (r *MockIdeaRepository) GetByID(ctx context.Context, id string) (*domain.Idea, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.ideas {
		if r.ideas[i].ID == id {
			idea := r.ideas[i] // Copy to avoid returning pointer to slice element
			return &idea, nil
		}
	}
	return nil, domain.NewNotFoundError("Idea")
}

// List retrieves ideas with optional filtering
func (r *MockIdeaRepository) List(ctx context.Context, opts domain.ListOptions) ([]domain.Idea, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Idea, 0, len(r.ideas))

	for _, idea := range r.ideas {
		// Apply filters
		if opts.Status != "" && idea.Status != opts.Status {
			continue
		}
		if opts.Category != "" && idea.Category != opts.Category {
			continue
		}
		if opts.MinVotes > 0 && idea.VoteCount < opts.MinVotes {
			continue
		}
		result = append(result, idea)
	}

	// Note: In a real implementation, sorting and pagination would be handled here
	// For mock, we return all matching results (already sorted by vote count in initial data)

	return result, nil
}

// Create creates a new idea
func (r *MockIdeaRepository) Create(ctx context.Context, idea *domain.Idea) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ideas = append(r.ideas, *idea)
	return nil
}

// Update updates an existing idea
func (r *MockIdeaRepository) Update(ctx context.Context, idea *domain.Idea) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.ideas {
		if r.ideas[i].ID == idea.ID {
			r.ideas[i] = *idea
			return nil
		}
	}
	return domain.NewNotFoundError("Idea")
}

// Delete removes an idea
func (r *MockIdeaRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.ideas {
		if r.ideas[i].ID == id {
			r.ideas = append(r.ideas[:i], r.ideas[i+1:]...)
			return nil
		}
	}
	return domain.NewNotFoundError("Idea")
}

// IncrementVoteCount atomically increments the vote count
func (r *MockIdeaRepository) IncrementVoteCount(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.ideas {
		if r.ideas[i].ID == id {
			r.ideas[i].VoteCount++
			return nil
		}
	}
	return domain.NewNotFoundError("Idea")
}

// Ensure MockIdeaRepository implements IdeaRepository
var _ domain.IdeaRepository = (*MockIdeaRepository)(nil)
