package handlers

import (
	"log"
	"net/http"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/abhi10/ideahamster/internal/domain"
	"github.com/abhi10/ideahamster/internal/models"
	"github.com/abhi10/ideahamster/internal/service"
	"github.com/abhi10/ideahamster/web/templates"
	"github.com/go-chi/chi/v5"
)

// Dependencies holds all handler dependencies for dependency injection.
// This struct allows handlers to access services without global state.
type Dependencies struct {
	IdeaService *service.IdeaService
}

// HandleVote processes a vote on an idea using injected dependencies
func (d *Dependencies) HandleVote(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")
	if ideaID == "" {
		http.Error(w, "Invalid idea ID", http.StatusBadRequest)
		return
	}

	// Get verified email from cookie
	email := getVerifiedEmail(r)
	if email == "" {
		// Trigger email modal via HTMX
		w.Header().Set("HX-Trigger", `{"showModal": "email-modal"}`)
		w.Header().Set("Content-Type", "text/html")
		component := templates.VoteButton(ideaID, false)
		component.Render(r.Context(), w)
		return
	}

	// Record vote using service
	if err := d.IdeaService.Vote(r.Context(), ideaID, email); err != nil {
		appErr := domain.GetAppError(err)

		// Already voted - show voted state
		if appErr.Type == domain.ErrorTypeConflict {
			component := templates.VoteButton(ideaID, true)
			component.Render(r.Context(), w)
			return
		}

		// Log and return error
		log.Printf("Vote error: %v", err)
		http.Error(w, appErr.UserMessage(), appErr.HTTPStatus())
		return
	}

	log.Printf("Vote registered: idea=%s, email=%s", ideaID, email)

	// Return voted button state
	component := templates.VoteButton(ideaID, true)
	component.Render(r.Context(), w)
}

// HandleLeaderboard renders the leaderboard using injected dependencies
func (d *Dependencies) HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	// This can be enabled when leaderboard is ready
	// For now, the redirect in main.go handles this

	// csrfToken := middleware.GetCSRFToken(r)
	// sortBy := r.URL.Query().Get("sort")
	// opts := domain.DefaultListOptions()
	//
	// switch sortBy {
	// case "eligible":
	//     opts.MinVotes = config.VoteThreshold
	// }
	//
	// ideas, err := d.IdeaService.ListIdeas(r.Context(), opts)
	// if err != nil {
	//     log.Printf("Failed to list ideas: %v", err)
	//     http.Error(w, "Failed to load ideas", http.StatusInternalServerError)
	//     return
	// }
	//
	// component := templates.Leaderboard(ideas, csrfToken)
	// component.Render(r.Context(), w)
}

// HandleExpandIdea returns expanded view of an idea using injected dependencies
func (d *Dependencies) HandleExpandIdea(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")
	if ideaID == "" {
		http.Error(w, "Invalid idea ID", http.StatusBadRequest)
		return
	}

	idea, err := d.IdeaService.GetIdea(r.Context(), ideaID)
	if err != nil {
		appErr := domain.GetAppError(err)
		http.Error(w, appErr.UserMessage(), appErr.HTTPStatus())
		return
	}

	// Convert to models.Idea for template compatibility
	modelIdea := convertDomainToModelIdea(idea)
	component := templates.IdeaCard(modelIdea, true)
	component.Render(r.Context(), w)
}

// convertDomainToModelIdea converts domain.Idea to models.Idea
// TODO: Update templates to use domain.Idea directly
func convertDomainToModelIdea(idea *domain.Idea) models.Idea {
	return models.Idea{
		ID:          idea.ID,
		Title:       idea.Title,
		Description: idea.Description,
		Category:    idea.Category,
		Tags:        idea.Tags,
		Submitter:   idea.Submitter,
		VoteCount:   idea.VoteCount,
		Status:      idea.Status,
		CreatedAt:   idea.CreatedAt,
		UpdatedAt:   idea.UpdatedAt,
	}
}

// getVerifiedEmail helper is shared with vote.go
func init() {
	// Validate config constants at startup
	if config.EmailCookieName == "" {
		panic("EmailCookieName must not be empty")
	}
}
