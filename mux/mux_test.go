package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mdigger/log"
)

func TestServeMux(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	testHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test", http.StatusOK)
	}

	Default.Headers = map[string]string{"test": "test"}
	Handle("", "", nil)
	Handle("GET", "/:name", http.HandlerFunc(testHandleFunc))
	Handle("GET", "/:name/test/", http.HandlerFunc(testHandleFunc))
	HandleFunc("POST", "/:name", testHandleFunc)

	for _, req := range []struct {
		method, url string
		status      int
	}{
		{"", "/test/", 301},
		{"", "/test/test", 301},
		{"POST", "/test/", 308},
		{"TEST", "/test/", 404},
		{"", "/test", 200},
		{"POST", "/test", 200},
		{"TEST", "/test", 405},
	} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(req.method, req.url, nil)
		Default.ServeHTTP(w, r)
		// dump, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	t.Error(err)
		// }
		// fmt.Printf("%s", dump)
		resp := w.Result()
		// dump, err = httputil.DumpResponse(resp, true)
		// if err != nil {
		// 	t.Error(err)
		// }
		// fmt.Printf("%s", dump)
		if resp.StatusCode != req.status {
			t.Errorf("bad status: %v - %03d", req, resp.StatusCode)
		}
		// fmt.Println(strings.Repeat("-", 60))
	}
}

func TestServeMuxError(t *testing.T) {
	defer func() {
		p := recover()
		if err, ok := p.(error); !ok ||
			err.Error() != "path parts overflow: 50000" {
			t.Error(p)
		}
	}()
	path := strings.Repeat("/test", 50000)
	Handle("GET", path, http.NotFoundHandler())
}

func TestParams(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	testHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		params := PathParams(r)
		if params == nil || len(params) != 1 || params.Get("name") != "test" {
			t.Error("bad params")
		}
		fmt.Println(params)
		http.Error(w, "test", http.StatusOK)
	}
	Handle("GET", "/:name", http.HandlerFunc(testHandleFunc))

	ts := httptest.NewServer(Default)
	defer ts.Close()

	_, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error(err)
	}

}
