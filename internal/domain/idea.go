package domain

import (
	"context"
	"time"
)

// ===========================================
// Domain Entities
// ===========================================

// Idea represents an app idea submitted by a user
type Idea struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Submitter   string    `json:"submitter"`
	VoteCount   int       `json:"vote_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Vote represents a user's vote on an idea
type Vote struct {
	ID        string    `json:"id"`
	IdeaID    string    `json:"idea_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ===========================================
// Constants
// ===========================================

// Category constants
const (
	CategoryFrontend  = "Frontend"
	CategoryBackend   = "Backend"
	CategoryFullStack = "Full Stack"
)

// Status constants
const (
	StatusPending  = "pending"
	StatusEligible = "eligible"
	StatusBuilding = "building"
	StatusBuilt    = "built"
	StatusArchived = "archived"
)

// ===========================================
// Repository Interfaces
// These allow for easy mocking in unit tests
// ===========================================

// IdeaRepository defines the interface for idea data access
type IdeaRepository interface {
	// GetByID retrieves an idea by its ID
	GetByID(ctx context.Context, id string) (*Idea, error)

	// List retrieves ideas with optional filtering and sorting
	List(ctx context.Context, opts ListOptions) ([]Idea, error)

	// Create creates a new idea
	Create(ctx context.Context, idea *Idea) error

	// Update updates an existing idea
	Update(ctx context.Context, idea *Idea) error

	// Delete removes an idea
	Delete(ctx context.Context, id string) error

	// IncrementVoteCount atomically increments the vote count
	IncrementVoteCount(ctx context.Context, id string) error
}

// VoteRepository defines the interface for vote data access
type VoteRepository interface {
	// GetByIdeaAndEmail checks if a user already voted on an idea
	GetByIdeaAndEmail(ctx context.Context, ideaID, email string) (*Vote, error)

	// Create records a new vote
	Create(ctx context.Context, vote *Vote) error

	// CountByIdea returns the vote count for an idea
	CountByIdea(ctx context.Context, ideaID string) (int, error)

	// ListByEmail returns all votes by a user
	ListByEmail(ctx context.Context, email string) ([]Vote, error)
}

// ===========================================
// Query Options
// ===========================================

// ListOptions defines filtering and sorting options for listing ideas
type ListOptions struct {
	// Sort field: "votes", "recent", "title"
	SortBy string
	// Sort direction: "asc" or "desc"
	SortDir string
	// Filter by status
	Status string
	// Filter by category
	Category string
	// Minimum vote count
	MinVotes int
	// Pagination
	Limit  int
	Offset int
}

// DefaultListOptions returns sensible defaults
func DefaultListOptions() ListOptions {
	return ListOptions{
		SortBy:  "votes",
		SortDir: "desc",
		Limit:   50,
		Offset:  0,
	}
}
