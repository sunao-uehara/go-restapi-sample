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
			s := &mysql.SampleData{}
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

func (h *Handler) StatsMiddleware(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		h.Log.Debug("before func")

		// // do something before `func`
		// now := time.Now()

		nextFunc(w, r)

		// // do something after `func`

		// elapsed := now - time.Now()
		// endpoint := r.URL.Path

		// uniq_id := unique_

		// // send it to ES.
		// json := `
		// 	{
		// 		"trans_id": uniq_id
		// 		"endpoint":  endpoint,
		// 		"elapsed_time" : elapsed
		// 		"container_id": XXX
		// 		"user_id": id,
		// 		"http_response": res.Number
		// 		"time": current_unixtimestamp
		// 	}
		// `
		// go func() {

		// }()

		h.Log.Debug("after func")
	}
}
