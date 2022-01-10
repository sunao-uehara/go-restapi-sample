package handler

import (
	"encoding/json"
	"net/http"

	mysql "github.com/sunao-uehara/go-restapi-sample/storages/mysql"
	myRedis "github.com/sunao-uehara/go-restapi-sample/storages/redis"
)

func (h *Handler) CacheMiddleware(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// do something before `func`
		h.Log.Debug("before func")
		// get the data from redis/cache first
		endpoint := r.URL.Path
		val, err := myRedis.GetCache(ctx, h.Redis, endpoint)
		if err == nil && val != "" {
			s := &mysql.Sample{}
			if err := json.Unmarshal([]byte(val), s); err == nil {
				h.Log.Debugf("get sample data from redis: %s", val)
				successJSONResponse(w, s)
				return
			}
		}

		nextFunc(w, r)

		// do something after `func`
		h.Log.Debug("after func")
	}
}
