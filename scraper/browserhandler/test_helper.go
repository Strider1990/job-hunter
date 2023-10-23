package browserhandler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

func NewTestServer(testAssetsPath string, tls ...bool) *TestServer {
	ts := &TestServer{
		Routes:             make(map[string]http.HandlerFunc),
		RequestSubscribers: make(map[string][]chan *http.Request),
		testAssetsPath:     testAssetsPath,
	}
	if len(tls) > 0 && tls[0] {
		ts.TestServer = httptest.NewTLSServer(http.HandlerFunc(ts.serveHTTP))
	} else {
		ts.TestServer = httptest.NewServer(http.HandlerFunc(ts.serveHTTP))
	}
	ts.PREFIX = ts.TestServer.URL
	ts.EMPTY_PAGE = ts.TestServer.URL + "/empty.html"
	ts.FORM_PAGE = ts.PREFIX + "/form.html"
	ts.CROSS_PROCESS_PREFIX = strings.Replace(ts.TestServer.URL, "127.0.0.1", "localhost", 1)
	return ts
}

type TestServer struct {
	sync.Mutex
	TestServer           *httptest.Server
	Routes               map[string]http.HandlerFunc
	RequestSubscribers   map[string][]chan *http.Request
	testAssetsPath       string
	PREFIX               string
	EMPTY_PAGE           string
	FORM_PAGE            string
	CROSS_PROCESS_PREFIX string
}

func (t *TestServer) AfterEach() {
	t.Lock()
	defer t.Unlock()
	t.Routes = make(map[string]http.HandlerFunc)
	t.RequestSubscribers = make(map[string][]chan *http.Request)
}

func (t *TestServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	t.Lock()
	defer t.Unlock()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if handlers, ok := t.RequestSubscribers[r.URL.Path]; ok {
		for _, handler := range handlers {
			handler <- r
		}
	}
	if route, ok := t.Routes[r.URL.Path]; ok {
		route(w, r)
		return
	}
	w.Header().Add("Cache-Control", "no-cache, no-store")
	http.FileServer(http.Dir(t.testAssetsPath)).ServeHTTP(w, r)
}

func (s *TestServer) SetBasicAuth(path, username, password string) {
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Add("WWW-Authenticate", "Basic") // needed or playwright will do not send auth header
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}
	})
}

func (s *TestServer) SetRoute(path string, f http.HandlerFunc) {
	s.Lock()
	defer s.Unlock()
	s.Routes[path] = f
}

func (s *TestServer) SetRedirect(from, to string) {
	s.SetRoute(from, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusFound)
	})
}

func (s *TestServer) WaitForRequestChan(path string) <-chan *http.Request {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.RequestSubscribers[path]; !ok {
		s.RequestSubscribers[path] = make([]chan *http.Request, 0)
	}
	channel := make(chan *http.Request, 1)
	s.RequestSubscribers[path] = append(s.RequestSubscribers[path], channel)
	return channel
}
