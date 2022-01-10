package handler

/*
import (
	"log"
	"net/http"
	"time"

	s "github.com/sunao-uehara/go-restapi-sample/storages"
)

// StatsMiddleware is a middleware that wraps another handler function
// to calculate and save elapsed time of the endpoint/handler.
func (h *Handler) StatsMiddleware(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		nextFunc(w, r)

		method := r.Method
		// TODO: this path contains REST resource id. e.g. /hash/1, it should be excluded to save proper stats by endpoint
		endpoint := r.URL.Path
		elapsed := float64(time.Since(start) / time.Microsecond)

		key := method + ":" + endpoint
		log.Println(key, elapsed)

		// save elapsed time to storaege
		s.SaveStats(key, elapsed)
	}
}
*/
