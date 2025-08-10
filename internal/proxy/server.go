package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"

	"github.com/gorilla/mux"
	"github.com/ramigou/go-gate/internal/config"
	"github.com/ramigou/go-gate/internal/middleware"
	"github.com/ramigou/go-gate/internal/router"
)

type Server struct {
	config *config.Config
	router *router.Router
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		router: router.New(cfg),
	}
}

func (s *Server) Handler() http.Handler {
	r := mux.NewRouter()
	
	r.Use(middleware.Logging)
	
	r.PathPrefix("/").HandlerFunc(s.proxyHandler)
	
	return r
}

func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	upstream := s.router.SelectUpstream(r)
	if upstream == nil {
		http.Error(w, "No upstream server available", http.StatusBadGateway)
		return
	}

	target, err := url.Parse(upstream.URL)
	if err != nil {
		log.Printf("Invalid upstream URL %s: %v", upstream.URL, err)
		http.Error(w, "Invalid upstream configuration", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for %s: %v", upstream.URL, err)
		http.Error(w, "Upstream server error", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}