package handlers

import (
	"net/http"

	"github.com/abhishekrajuchamarthi/idea-hamster/web/templates"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	// Render the home page template
	component := templates.Home()
	component.Render(r.Context(), w)
}
