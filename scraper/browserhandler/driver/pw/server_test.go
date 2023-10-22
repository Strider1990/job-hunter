package pw_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

func newTestServer(tls ...bool) *testServer {
	ts := &testServer{
		routes:              make(map[string]http.HandlerFunc),
		requestSubscriberes: make(map[string][]chan *http.Request),
	}
	if len(tls) > 0 && tls[0] {
		ts.testServer = httptest.NewTLSServer(http.HandlerFunc(ts.serveHTTP))
	} else {
		ts.testServer = httptest.NewServer(http.HandlerFunc(ts.serveHTTP))
	}
	ts.PREFIX = ts.testServer.URL
	ts.EMPTY_PAGE = ts.testServer.URL + "/empty.html"
	ts.CROSS_PROCESS_PREFIX = strings.Replace(ts.testServer.URL, "127.0.0.1", "localhost", 1)
	return ts
}

type testServer struct {
	sync.Mutex
	testServer           *httptest.Server
	routes               map[string]http.HandlerFunc
	requestSubscriberes  map[string][]chan *http.Request
	PREFIX               string
	EMPTY_PAGE           string
	CROSS_PROCESS_PREFIX string
}

func (t *testServer) AfterEach() {
	t.Lock()
	defer t.Unlock()
	t.routes = make(map[string]http.HandlerFunc)
	t.requestSubscriberes = make(map[string][]chan *http.Request)
}

func (t *testServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	t.Lock()
	defer t.Unlock()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if handlers, ok := t.requestSubscriberes[r.URL.Path]; ok {
		for _, handler := range handlers {
			handler <- r
		}
	}
	if route, ok := t.routes[r.URL.Path]; ok {
		route(w, r)
		return
	}
	w.Header().Add("Cache-Control", "no-cache, no-store")
	http.FileServer(http.Dir("assets")).ServeHTTP(w, r)
}

func (s *testServer) SetBasicAuth(path, username, password string) {
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Add("WWW-Authenticate", "Basic") // needed or playwright will do not send auth header
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}
	})
}

func (s *testServer) SetRoute(path string, f http.HandlerFunc) {
	s.Lock()
	defer s.Unlock()
	s.routes[path] = f
}

func (s *testServer) SetRedirect(from, to string) {
	s.SetRoute(from, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusFound)
	})
}

func (s *testServer) WaitForRequestChan(path string) <-chan *http.Request {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.requestSubscriberes[path]; !ok {
		s.requestSubscriberes[path] = make([]chan *http.Request, 0)
	}
	channel := make(chan *http.Request, 1)
	s.requestSubscriberes[path] = append(s.requestSubscriberes[path], channel)
	return channel
}
