package models

import "time"

// Idea represents an app idea submitted by a user
type Idea struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"` // Frontend, Backend, Full Stack
	Tags        []string  `json:"tags" db:"tags"`
	Submitter   string    `json:"submitter" db:"submitter"` // Optional name
	VoteCount   int       `json:"vote_count" db:"vote_count"`
	Status      string    `json:"status" db:"status"` // pending, eligible, building, built
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Vote represents a user's vote on an idea
type Vote struct {
	ID        string    `json:"id" db:"id"`
	IdeaID    string    `json:"idea_id" db:"idea_id"`
	Email     string    `json:"email" db:"email"`
	Verified  bool      `json:"verified" db:"verified"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Category constants
const (
	CategoryFrontend  = "Frontend"
	CategoryBackend   = "Backend"
	CategoryFullStack = "Full Stack"
)

// Status constants
const (
	StatusPending   = "pending"
	StatusEligible  = "eligible"  // 50+ votes
	StatusBuilding  = "building"  // Phase 2
	StatusBuilt     = "built"     // Phase 2
	StatusArchived  = "archived"  // Phase 2
)
