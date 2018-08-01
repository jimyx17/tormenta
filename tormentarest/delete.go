package tormentarest

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jpincas/gouuidv6"
)

func deleteByID(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new one of that type to hold the results of the query
	idString := chi.URLParam(r, "id")
	entityName := strings.Split(r.URL.Path, "/")[1]

	id := gouuidv6.UUID{}
	if err := id.UnmarshalText([]byte(idString)); err != nil {
		renderError(w, errBadIDFormat, idString)
		return
	}

	// Delete the record
	n, err := App.DB.Delete(entityName, id)
	if err != nil {
		renderError(w, errDBConnection)
		return
	}

	if n == 0 {
		renderError(w, errRecordNotFound, idString)
		return
	}

	// Render
	renderSuccess(w, n)
	return
}