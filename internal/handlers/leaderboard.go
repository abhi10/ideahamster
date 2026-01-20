package handlers

import (
	"net/http"

	"github.com/abhi10/idea-hamster/internal/models"
	"github.com/abhi10/idea-hamster/web/templates"
)

func HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch from database when connected
	// For now, use mock data
	ideas := getMockIdeas()

	// Sort based on query parameter
	sortBy := r.URL.Query().Get("sort")
	switch sortBy {
	case "recent":
		// Sort by created_at desc (already in order for mock)
	case "eligible":
		// Filter only ideas with 50+ votes
		eligible := []models.Idea{}
		for _, idea := range ideas {
			if idea.VoteCount >= 50 {
				eligible = append(eligible, idea)
			}
		}
		ideas = eligible
	default:
		// Default: sort by votes desc (already in order for mock)
	}

	component := templates.Leaderboard(ideas)
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
