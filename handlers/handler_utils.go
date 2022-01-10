package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

func successResponse(w http.ResponseWriter, res interface{}) {
	var resStr string

	// support int and string only for now
	t := reflect.TypeOf(res)
	switch t.Kind() {
	case reflect.Int:
		resStr = strconv.Itoa(res.(int))
	case reflect.String:
		resStr = res.(string)
	}

	w.WriteHeader(200)
	w.Write([]byte(resStr))
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// successJSONResponse respond success http status(200) with given data in json
func successJSONResponse(w http.ResponseWriter, data interface{}) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		log.Printf("cannot convert struct %v to json", data)
		errorJSONResponse(w, 500, "Unknown Error Occurred")
		return
	}

	// log.Println(string(jsonString))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write([]byte(jsonString))
}

// errorJSONResponse respond given failure/error http status(4xx~5xx) with message in json
func errorJSONResponse(w http.ResponseWriter, code int, msg string) {
	jsonString, _ := json.Marshal(DefaultResponseBody(msg))

	// log.Println(string(jsonString))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write([]byte(jsonString))
}

func DefaultResponseBody(msg string) interface{} {
	type Data struct {
		Msg string `json:"message"`
	}

	return Data{Msg: msg}
}
