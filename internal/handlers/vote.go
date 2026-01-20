package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abhi10/idea-hamster/internal/models"
	"github.com/abhi10/idea-hamster/web/templates"
	"github.com/go-chi/chi/v5"
)

// Session cookie name for storing verified email
const emailCookieName = "voter_email"

// HandleVote processes a vote on an idea
func HandleVote(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")

	// Check if user has verified email in session
	cookie, err := r.Cookie(emailCookieName)
	if err != nil || cookie.Value == "" {
		// No email verified - show verification modal
		w.Header().Set("HX-Trigger", "showEmailModal")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<script>document.getElementById('email-modal').style.display='flex'</script>`))
		return
	}

	email := cookie.Value

	// TODO: Save vote to database
	// For now, simulate voting
	fmt.Printf("Vote registered: idea=%s, email=%s\n", ideaID, email)

	// Return updated vote button (showing voted state)
	component := templates.VoteButton(ideaID, true)
	component.Render(r.Context(), w)
}

// HandleVerifyEmail verifies an email and stores it in session
func HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	// TODO: Send verification code via email
	// For Phase 1, we'll just accept the email and set cookie

	// Set cookie with email (expires in 1 year)
	http.SetCookie(w, &http.Cookie{
		Name:     emailCookieName,
		Value:    email,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Close modal and show success
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<script>
			document.getElementById('email-modal').style.display='none';
			alert('✅ Email verified! You can now vote.');
			window.location.reload();
		</script>
	`))
}

// HandleExpandIdea returns expanded view of an idea
func HandleExpandIdea(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")

	// TODO: Fetch from database
	// For now, find in mock data
	ideas := getMockIdeas()
	var idea *models.Idea
	for _, i := range ideas {
		if i.ID == ideaID {
			idea = &i
			break
		}
	}

	if idea == nil {
		http.Error(w, "Idea not found", http.StatusNotFound)
		return
	}

	// Return expanded card
	component := templates.IdeaCard(*idea, true)
	component.Render(r.Context(), w)
}

// API response structure
type VoteResponse struct {
	Success   bool   `json:"success"`
	NewCount  int    `json:"new_count"`
	Message   string `json:"message,omitempty"`
	NeedsAuth bool   `json:"needs_auth,omitempty"`
}

// HandleVoteAPI is a JSON API endpoint for voting (alternative to HTMX)
func HandleVoteAPI(w http.ResponseWriter, r *http.Request) {
	ideaID := chi.URLParam(r, "ideaID")

	// Check if user has verified email
	cookie, err := r.Cookie(emailCookieName)
	if err != nil || cookie.Value == "" {
		json.NewEncoder(w).Encode(VoteResponse{
			Success:   false,
			NeedsAuth: true,
			Message:   "Email verification required",
		})
		return
	}

	email := cookie.Value

	// TODO: Save vote to database and get new count
	// For now, simulate
	fmt.Printf("API Vote: idea=%s, email=%s\n", ideaID, email)
	newCount := 42

	json.NewEncoder(w).Encode(VoteResponse{
		Success:  true,
		NewCount: newCount,
		Message:  "Vote recorded!",
	})
}
