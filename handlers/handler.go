package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	mysql "github.com/sunao-uehara/go-restapi-sample/storages/mysql"
)

type Handler struct {
	*HandlerOptions
}
type HandlerOptions struct {
	Wg    *sync.WaitGroup
	Mysql *sql.DB
	// Redis
	Log *zap.SugaredLogger
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

	req := &SamplePostRequest{}
	if err := render.Bind(r, req); err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusBadRequest, "cannot create record")
		return
	}
	// h.Log.Debugf("param foo: %s", req.Foo)
	// h.Log.Debugf("param int_val: %d", req.IntVal)

	// write the data into MySQL
	id, err := mysql.CreateSample(h.Mysql, &mysql.Sample{Foo: req.Foo, IntVal: req.IntVal})
	if err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusInternalServerError, "cannot create record")
	}

	// execute asynchronously
	h.Wg.Add(1)
	go func() {
		h.Log.Debug("sample goroutien start")
		defer h.Wg.Done()

		// wait X seconds
		time.Sleep(3 * time.Second)
		// write the data into Redis

		h.Log.Debug("sample goroutine done")
	}()

	type Res struct {
		ID int64 `json:"id"`
	}
	res := &Res{id}
	successJSONResponse(w, res)
}

func (h *Handler) SampleGetHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("SampleGetHandler")

	sampleId := chi.URLParam(r, "sampleId")
	if sampleId != "" {
		id, _ := strconv.Atoi(sampleId)
		data, err := mysql.GetSample(h.Mysql, id)
		if err != nil {
			h.Log.Debug(err)
			errorJSONResponse(w, http.StatusNotFound, "Not Found")
			return
		}
		h.Log.Debug(data)

		successJSONResponse(w, data)
		return
	}

	data, err := mysql.GetManySample(h.Mysql)
	if err != nil {
		h.Log.Debug(err)
		errorJSONResponse(w, http.StatusNotFound, "Not Found")
		return
	}

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

	id, _ := strconv.Atoi(sampleId)

	d := &mysql.Sample{
		Foo:    req.Foo,
		IntVal: req.IntVal,
	}
	rowsAffected, err := mysql.UpdateSample(h.Mysql, id, d)
	if err != nil {
		h.Log.Info(err.Error())
		errorJSONResponse(w, http.StatusBadRequest, "failed patch request")
		return
	}
	h.Log.Debug(rowsAffected, " rows affected")

	type Res struct {
		Message string `json:"message"`
	}
	res := &Res{Message: fmt.Sprintf("%d rows affected", rowsAffected)}
	successJSONResponse(w, res)
}
