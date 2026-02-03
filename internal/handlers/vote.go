package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/abhi10/ideahamster/internal/models"
	"github.com/abhi10/ideahamster/internal/sanitizer"
	"github.com/abhi10/ideahamster/web/templates"
	"github.com/go-chi/chi/v5"
)

// getVerifiedEmail extracts and validates the voter email from cookie.
// Returns empty string if not verified.
func getVerifiedEmail(r *http.Request) string {
	cookie, err := r.Cookie(config.EmailCookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}
	return cookie.Value
}

// HandleVote processes a vote on an idea
func HandleVote(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")
	if ideaID == "" {
		http.Error(w, "Invalid idea ID", http.StatusBadRequest)
		return
	}

	// Check if user has verified email in session
	email := getVerifiedEmail(r)
	if email == "" {
		// No email verified - trigger modal via HTMX event
		w.Header().Set("HX-Trigger", `{"showModal": "email-modal"}`)
		w.Header().Set("Content-Type", "text/html")
		component := templates.VoteButton(ideaID, false)
		component.Render(r.Context(), w)
		return
	}

	// TODO: Save vote to database
	log.Printf("Vote registered: idea=%s, email=%s", ideaID, email)

	// Return updated vote button (showing voted state)
	component := templates.VoteButton(ideaID, true)
	component.Render(r.Context(), w)
}

// HandleVerifyEmail verifies an email and stores it in session
func HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	rawEmail := r.FormValue("email")

	// Sanitize and validate email
	email, err := sanitizer.SanitizeEmail(rawEmail)
	if err != nil {
		log.Printf("Email validation failed: %v", err)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`
			<div class="retro-card max-w-md w-full mx-4">
				<h2 class="text-xl font-pixel theme-primary mb-4">ERROR</h2>
				<p class="text-red-400 mb-4">%s</p>
				<a href="/" class="retro-button inline-block text-center">
					TRY AGAIN
				</a>
			</div>
		`, err.Error())))
		return
	}

	// Set cookie with email
	http.SetCookie(w, &http.Cookie{
		Name:     config.EmailCookieName,
		Value:    email,
		Path:     "/",
		MaxAge:   config.CookieMaxAge,
		HttpOnly: true,
		Secure:   config.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	})

	// Close modal and trigger page refresh via HTMX
	w.Header().Set("HX-Trigger", `{"closeModal": "email-modal", "showNotification": "Email verified! You can now vote."}`)
	w.Header().Set("HX-Refresh", "true")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="text-center p-4 theme-success">Email verified!</div>`))
}

// HandleExpandIdea returns expanded view of an idea
func HandleExpandIdea(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")
	if ideaID == "" {
		http.Error(w, "Invalid idea ID", http.StatusBadRequest)
		return
	}

	// TODO: Fetch from database
	ideas := getMockIdeas()
	var idea *models.Idea
	for i := range ideas {
		if ideas[i].ID == ideaID {
			idea = &ideas[i]
			break
		}
	}

	if idea == nil {
		http.Error(w, "Idea not found", http.StatusNotFound)
		return
	}

	component := templates.IdeaCard(*idea, true)
	component.Render(r.Context(), w)
}

// VoteResponse is the API response structure for voting
type VoteResponse struct {
	Success   bool   `json:"success"`
	NewCount  int    `json:"new_count"`
	Message   string `json:"message,omitempty"`
	NeedsAuth bool   `json:"needs_auth,omitempty"`
}

// HandleVoteAPI is a JSON API endpoint for voting (alternative to HTMX)
func HandleVoteAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ideaID := chi.URLParam(r, "ideaID")
	if ideaID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VoteResponse{
			Success: false,
			Message: "Invalid idea ID",
		})
		return
	}

	// Check if user has verified email
	email := getVerifiedEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VoteResponse{
			Success:   false,
			NeedsAuth: true,
			Message:   "Email verification required",
		})
		return
	}

	// TODO: Save vote to database and get new count
	log.Printf("API Vote: idea=%s, email=%s", ideaID, email)
	newCount := 42 // Mock count

	json.NewEncoder(w).Encode(VoteResponse{
		Success:  true,
		NewCount: newCount,
		Message:  "Vote recorded!",
	})
}
