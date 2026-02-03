package handlers

import (
	"net/http"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/abhi10/ideahamster/internal/middleware"
	"github.com/abhi10/ideahamster/internal/models"
	"github.com/abhi10/ideahamster/web/templates"
)

// validSortParams is a whitelist of allowed sort parameters
var validSortParams = map[string]bool{
	"":         true, // default (votes)
	"votes":    true,
	"recent":   true,
	"eligible": true,
}

// HandleLeaderboard renders the leaderboard page with ideas
func HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	csrfToken := middleware.GetCSRFToken(r)

	// TODO: Fetch from database when connected
	ideas := getMockIdeas()

	// Validate and sanitize sort parameter
	sortBy := r.URL.Query().Get("sort")
	if !validSortParams[sortBy] {
		sortBy = "votes"
	}

	// Apply sorting/filtering
	switch sortBy {
	case "recent":
		// Sort by created_at desc (already in order for mock)
	case "eligible":
		// Filter only ideas meeting vote threshold
		eligible := make([]models.Idea, 0)
		for _, idea := range ideas {
			if idea.VoteCount >= config.VoteThreshold {
				eligible = append(eligible, idea)
			}
		}
		ideas = eligible
	default:
		// Default: sort by votes desc (already in order for mock)
	}

	component := templates.Leaderboard(ideas, csrfToken)
	component.Render(r.Context(), w)
}

// getMockIdeas returns sample data for testing
// TODO: Replace with database queries
func getMockIdeas() []models.Idea {
	return []models.Idea{
		{
			ID:          "1",
			Title:       "AI-Powered Recipe Generator",
			Description: "Generate custom recipes based on ingredients you have at home. Uses AI to suggest creative combinations and cooking instructions. Perfect for reducing food waste and discovering new meals!",
			Category:    "Full Stack",
			Tags:        []string{"AI", "Food", "Mobile"},
			Submitter:   "Chef Mike",
			VoteCount:   67,
			Status:      models.StatusEligible,
		},
		{
			ID:          "2",
			Title:       "Retro Pomodoro Timer",
			Description: "A productivity timer with 90s aesthetics. Track your work sessions with pixel art animations and chiptune sounds. Includes stats tracking and customizable work/break intervals.",
			Category:    "Frontend",
			Tags:        []string{"Productivity", "Retro", "Timer"},
			Submitter:   "ProductivityGuru",
			VoteCount:   52,
			Status:      models.StatusEligible,
		},
		{
			ID:          "3",
			Title:       "Local Event Discovery API",
			Description: "API that aggregates events from multiple sources and provides a unified interface for discovering local happenings. Supports filtering by date, category, and location.",
			Category:    "Backend",
			Tags:        []string{"API", "Events", "Location"},
			Submitter:   "DevCommunity",
			VoteCount:   45,
			Status:      models.StatusPending,
		},
		{
			ID:          "4",
			Title:       "Habit Tracker with Friends",
			Description: "Track daily habits and compete with friends. See who can maintain the longest streak. Gamified habit building with rewards and achievements.",
			Category:    "Full Stack",
			Tags:        []string{"Habits", "Social", "Gamification"},
			Submitter:   "LifeHacker",
			VoteCount:   38,
			Status:      models.StatusPending,
		},
		{
			ID:          "5",
			Title:       "Markdown-Based Wiki",
			Description: "Simple wiki system powered by Markdown files. Perfect for personal knowledge bases or team documentation. Supports full-text search and automatic linking.",
			Category:    "Full Stack",
			Tags:        []string{"Documentation", "Markdown", "Knowledge"},
			Submitter:   "",
			VoteCount:   31,
			Status:      models.StatusPending,
		},
	}
}
