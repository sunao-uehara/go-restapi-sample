package router

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	handler "github.com/sunao-uehara/go-restapi-sample/handlers"
)

// NewRouter registers and associates the endpoint with Handler function
// returns registered handlers
func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", h.IndexHandler)

	// /sample
	r.Route("/sample", func(r chi.Router) {
		r.Post("/", h.SamplePostHandler)
		r.Get("/", h.SampleGetHandler)
		r.Get("/{sampleId}", h.SampleGetHandler)
		r.Patch("/{sampleId}", h.SamplePatchHandler)
		// r.Put("/{sampleId}", h.SamplePostHandler)
		// r.Delete("/{sampleId}", h.SamplePostHandler)
	})

	// route not exits

	return r
}
