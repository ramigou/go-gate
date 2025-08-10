package router

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ramigou/go-gate/internal/config"
)

type Router struct {
	config    *config.Config
	upstreams map[string]*config.UpstreamConfig
	rnd       *rand.Rand
}

func New(cfg *config.Config) *Router {
	upstreams := make(map[string]*config.UpstreamConfig)
	for i := range cfg.Upstreams {
		upstreams[cfg.Upstreams[i].Name] = &cfg.Upstreams[i]
	}

	return &Router{
		config:    cfg,
		upstreams: upstreams,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Router) SelectUpstream(req *http.Request) *config.UpstreamConfig {
	for _, route := range r.config.Routes {
		if r.matchRoute(req, &route) {
			return r.selectUpstreamForRoute(&route)
		}
	}

	if len(r.config.Upstreams) > 0 {
		return &r.config.Upstreams[r.rnd.Intn(len(r.config.Upstreams))]
	}

	return nil
}

func (r *Router) matchRoute(req *http.Request, route *config.RouteConfig) bool {
	if route.Match.Host != "" {
		if !r.matchHost(req.Host, route.Match.Host) {
			return false
		}
	}

	if route.Match.Path != "" {
		if !r.matchPath(req.URL.Path, route.Match.Path) {
			return false
		}
	}

	return true
}

func (r *Router) matchHost(reqHost, routeHost string) bool {
	reqHost = strings.Split(reqHost, ":")[0]
	
	if routeHost == reqHost {
		return true
	}

	if strings.HasPrefix(routeHost, "*.") {
		domain := routeHost[2:]
		return strings.HasSuffix(reqHost, "."+domain) || reqHost == domain
	}

	return false
}

func (r *Router) matchPath(reqPath, routePath string) bool {
	if strings.HasSuffix(routePath, "*") {
		prefix := routePath[:len(routePath)-1]
		return strings.HasPrefix(reqPath, prefix)
	}

	return reqPath == routePath
}

func (r *Router) selectUpstreamForRoute(route *config.RouteConfig) *config.UpstreamConfig {
	if route.Upstream != "" {
		return r.upstreams[route.Upstream]
	}

	if len(route.Upstreams) == 0 {
		return nil
	}

	if len(route.Upstreams) == 1 {
		return r.upstreams[route.Upstreams[0]]
	}

	var candidates []*config.UpstreamConfig
	var totalWeight int

	for _, upstreamName := range route.Upstreams {
		if upstream, exists := r.upstreams[upstreamName]; exists {
			weight := upstream.Weight
			if weight <= 0 {
				weight = 1
			}
			for i := 0; i < weight; i++ {
				candidates = append(candidates, upstream)
			}
			totalWeight += weight
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	return candidates[r.rnd.Intn(len(candidates))]
}