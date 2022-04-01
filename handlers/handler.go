package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	mysql "github.com/sunao-uehara/go-restapi-sample/storages/mysql"
	myRedis "github.com/sunao-uehara/go-restapi-sample/storages/redis"
)

type Handler struct {
	*HandlerOptions
}
type HandlerOptions struct {
	Wg    *sync.WaitGroup
	Mysql *sql.DB
	Redis *redis.Client
	Log   *zap.SugaredLogger
}

func NewHandler(handlerOptions *HandlerOptions) *Handler {
	return &Handler{
		handlerOptions,
	}
}

// IndexHandler returns output 'hello handler'
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	successResponse(w, "hello world")
}

type SamplePostRequest struct {
	Foo    string `json:"foo"`
	IntVal int64  `json:"int_val"`
}

func (spr *SamplePostRequest) Bind(r *http.Request) error {
	if spr.Foo == "" {
		return errors.New("missin required field: Foo")
	}

	return nil
}

func (h *Handler) SamplePostHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("SamplePostHandler")
	// ctx := r.Context()

	req := &SamplePostRequest{}
	if err := render.Bind(r, req); err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusBadRequest, "cannot create record")
		return
	}
	// h.Log.Debugf("param foo: %s", req.Foo)
	// h.Log.Debugf("param int_val: %d", req.IntVal)

	// write the data into MySQL
	sc := mysql.NewSample(h.Mysql)
	id, err := sc.CreateSample(&mysql.SampleData{Foo: req.Foo, IntVal: req.IntVal})
	if err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusInternalServerError, "cannot create record")
	}

	h.Wg.Add(1)
	go func() {
		defer h.Wg.Done()

		// purge cache
		h.purgeCache(context.Background(), []string{"/sample", "/sample/"})
	}()

	type Res struct {
		ID int64 `json:"id"`
	}
	res := &Res{id}
	successJSONResponse(w, res)
}

func (h *Handler) SampleGetHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("SampleGetHandler")
	// ctx := r.Context()

	sampleId := chi.URLParam(r, "sampleId")
	if sampleId != "" {
		id, _ := strconv.ParseInt(sampleId, 10, 64)

		// get the data from mysql
		sc := mysql.NewSample(h.Mysql)
		data, err := sc.GetSample(id)
		if err != nil {
			h.Log.Debug(err)
			errorJSONResponse(w, http.StatusNotFound, "Not Found")
			return
		}
		h.Log.Debug(data)

		// execute asynchronously
		h.Wg.Add(1)
		go func() {
			h.Log.Debug("sample goroutine start")
			defer h.Wg.Done()

			// wait X seconds for testing graceful shutdown for goroutine
			// time.Sleep(3 * time.Second)

			// write the data into Redis
			h.setCache(context.Background(), r.URL.Path, data)
			h.Log.Debug("sample goroutine done")
		}()

		successJSONResponse(w, data)
		return
	}

	sc := mysql.NewSample(h.Mysql)
	data, err := sc.GetManySample()
	if err != nil {
		h.Log.Debug(err)
		errorJSONResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	// execute asynchronously
	h.Wg.Add(1)
	go func() {
		defer h.Wg.Done()
		h.setCache(context.Background(), r.URL.Path, data)
	}()

	successJSONResponse(w, data)
}

type SamplePatchRequest struct {
	Foo    string `json:"foo,omitempty"`
	IntVal int64  `json:"int_val,omitempty"`
}

func (spr *SamplePatchRequest) Bind(r *http.Request) error {
	if spr.Foo == "" && spr.IntVal == 0 {
		return errors.New("missing patch fields")
	}

	return nil
}

func (h *Handler) SamplePatchHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("SamplePatchHandler")

	req := &SamplePatchRequest{}
	if err := render.Bind(r, req); err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusBadRequest, "cannot apply patch")
		return
	}

	sampleId := chi.URLParam(r, "sampleId")
	if sampleId == "" {
		errorJSONResponse(w, http.StatusBadRequest, "Not Found")
		return
	}

	id, _ := strconv.ParseInt(sampleId, 10, 64)
	d := &mysql.SampleData{
		Foo:    req.Foo,
		IntVal: req.IntVal,
	}

	sc := mysql.NewSample(h.Mysql)
	rowsAffected, err := sc.UpdateSample(id, d)
	if err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusBadRequest, "failed patch request")
		return
	}
	h.Log.Debug(rowsAffected, " rows affected")

	h.Wg.Add(1)
	go func() {
		defer h.Wg.Done()

		// purge cache
		endpoints := []string{
			"/sample",
			"/sample/",
			r.URL.Path,
		}
		h.purgeCache(context.Background(), endpoints)
	}()

	type Res struct {
		Message string `json:"message"`
	}
	res := &Res{Message: fmt.Sprintf("%d rows affected", rowsAffected)}
	successJSONResponse(w, res)
}

func (h *Handler) setCache(ctx context.Context, endpoint string, data interface{}) {
	// write the data into Redis
	d, err := json.Marshal(data)
	if err != nil {
		h.Log.Error(err.Error())
	} else {
		if err := myRedis.SetCache(context.Background(), h.Redis, endpoint, string(d)); err != nil {
			h.Log.Error(err.Error())
		}
	}
}

func (h *Handler) purgeCache(ctx context.Context, endpoints []string) {
	for _, e := range endpoints {
		if err := myRedis.DelCache(context.Background(), h.Redis, e); err != nil {
			h.Log.Error(err.Error())
		}
	}
}
