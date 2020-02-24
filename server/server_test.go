package server_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/marrbor/go-http-server/server"
	"github.com/marrbor/gohttp"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

func nowStr(t time.Time) string {
	return t.Format("2006/01/02T15:04:05")
}

type testA struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var a = testA{ID: 1, Name: "taro"}

type testB struct {
	Now time.Time `json:"now"`
}

var b = testB{Now: now}

func A(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		gohttp.MethodNotAllowed(w, nil)
	}
	if err := gohttp.JSONResponse(w, &a); err != nil {
		gohttp.InternalServerError(w, err)
	}
}

func B(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		gohttp.MethodNotAllowed(w, nil)
	}

	if err := gohttp.JSONResponse(w, &b); err != nil {
		gohttp.InternalServerError(w, err)
	}
}

var eps = []server.EntryPoint{
	{Resource: "a", Function: A},
	{Resource: "b", Function: B},
}

func TestNewServer1(t *testing.T) {
	hs := server.NewServer(8888, eps, nil)
	ch := make(chan error)
	hs.Start(ch)

	// A test
	res, err := http.Get("http://localhost:8888/a")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	var ra testA
	err = gohttp.ResponseJSONToParams(res, &ra)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ID, ra.ID)
	assert.EqualValues(t, a.Name, ra.Name)

	res, err = http.Head("http://localhost:8888/a")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	res, err = http.Get("http://localhost:8888/a/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	err = gohttp.ResponseJSONToParams(res, &ra)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ID, ra.ID)
	assert.EqualValues(t, a.Name, ra.Name)

	res, err = http.Head("http://localhost:8888/a/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	// not found test
	res, err = http.Get("http://localhost:8888/ab")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusNotFound, res.StatusCode)

	res, err = http.Head("http://localhost:8888/ab")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusNotFound, res.StatusCode)

	// B test
	res, err = http.Get("http://localhost:8888/b")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	var rb testB
	err = gohttp.ResponseJSONToParams(res, &rb)
	assert.NoError(t, err)
	assert.EqualValues(t, nowStr(now), nowStr(rb.Now))

	res, err = http.Head("http://localhost:8888/b")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	res, err = http.Get("http://localhost:8888/b/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	err = gohttp.ResponseJSONToParams(res, &rb)
	assert.NoError(t, err)
	assert.EqualValues(t, nowStr(now), nowStr(rb.Now))

	res, err = http.Head("http://localhost:8888/b/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	hs.Stop()
}

func defFunc(w http.ResponseWriter, r *http.Request) {
	gohttp.InternalServerError(w, nil)
}

func TestNewServer2(t *testing.T) {
	hs := server.NewServer(8888, eps, defFunc)
	ch := make(chan error)
	hs.Start(ch)

	// A test
	res, err := http.Get("http://localhost:8888/a")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	var ra testA
	err = gohttp.ResponseJSONToParams(res, &ra)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ID, ra.ID)
	assert.EqualValues(t, a.Name, ra.Name)

	res, err = http.Head("http://localhost:8888/a")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	res, err = http.Get("http://localhost:8888/a/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	err = gohttp.ResponseJSONToParams(res, &ra)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ID, ra.ID)
	assert.EqualValues(t, a.Name, ra.Name)

	res, err = http.Head("http://localhost:8888/a/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	res, err = http.Get("http://localhost:8888/ab")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)

	res, err = http.Head("http://localhost:8888/ab")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)

	// B test
	res, err = http.Get("http://localhost:8888/b")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	var rb testB
	err = gohttp.ResponseJSONToParams(res, &rb)
	assert.NoError(t, err)
	assert.EqualValues(t, nowStr(now), nowStr(rb.Now))

	res, err = http.Head("http://localhost:8888/b")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	res, err = http.Get("http://localhost:8888/b/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	err = gohttp.ResponseJSONToParams(res, &rb)
	assert.NoError(t, err)
	assert.EqualValues(t, nowStr(now), nowStr(rb.Now))

	res, err = http.Head("http://localhost:8888/b/")
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, res.StatusCode)

	hs.Stop()
}
