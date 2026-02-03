package handlers

import (
	"net/http"

	"github.com/abhi10/ideahamster/internal/middleware"
	"github.com/abhi10/ideahamster/web/templates"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	// Get CSRF token for form protection
	csrfToken := middleware.GetCSRFToken(r)

	// Render the home page template
	component := templates.Home(csrfToken)
	component.Render(r.Context(), w)
}
