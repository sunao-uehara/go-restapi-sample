package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
)

type Server struct {
	Header map[string]string
}

func NewServer() *Server {
	return &Server{
		// default headers
		Header: map[string]string{
			"Content-Type":    "application/json",
			"X-Forwarded-For": "123.234.345.456",
		},
	}
}

func (s *Server) GET(endpoint string, handler http.HandlerFunc) (*http.Response, error) {
	return s.exec(http.MethodGet, endpoint, handler, nil)
}

func (s *Server) POST(endpoint string, handler http.HandlerFunc, formData url.Values) (*http.Response, error) {
	return s.exec(http.MethodPost, endpoint, handler, formData)
}

func (s *Server) exec(httpMethod, endpoint string, handler http.HandlerFunc, formData url.Values) (*http.Response, error) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	var buf io.Reader
	if formData != nil {
		buf = strings.NewReader(formData.Encode())
	}

	req, err := http.NewRequest(httpMethod, ts.URL+endpoint, buf)
	if err != nil {
		return nil, err
	}

	for k, v := range s.Header {
		req.Header.Set(k, v)
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) AddHeader(headers map[string]string) {
	for k, v := range headers {
		s.Header[k] = v
	}
}

func ParseResponseBody(res *http.Response) (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("read from HTTP Response Body failed: %v", err)
	}

	return string(body), nil
}

func EqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error

	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	x1 := bytes.TrimLeft([]byte(s1), " \t\r\n")
	isS1Array := len(x1) > 0 && x1[0] == '['

	x2 := bytes.TrimLeft([]byte(s2), " \t\r\n")
	isS2Array := len(x2) > 0 && x2[0] == '['

	if isS1Array && isS2Array {
		// if it is json array, just check count for now
		// TODO: sort the array and check the value as well.
		v1, ok := o1.([]interface{})
		if !ok {
			return false, nil
		}
		v2, ok := o2.([]interface{})
		if !ok {
			return false, nil
		}

		return len(v1) == len(v2), nil
	}
	return reflect.DeepEqual(o1, o2), nil
}

type HandlerInput struct {
	Endpoint    string
	RequestBody string
}

type HandlerOutput struct {
	ResponseStatusCode int
	ResponseBody       string
	ResponseBodyStruct interface{}
}

type HandlerTestCase struct {
	Scenario string
	In       HandlerInput
	Out      HandlerOutput
}
