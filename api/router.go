package api

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/NikolaiMarkalainen/Router/utils"
)

type handlerFunc func(*utils.ResponseWriter, *http.Request)

type route struct {
	method string
	pattern *regexp.Regexp
	innerHandler handlerFunc
	paramKeys []string
}

type router struct {
	routes []route
}

func NewRouter() *router {
	return &router{routes: []route{}}
}




func (r *router) AddRoute(method, endpoint string, handler handlerFunc) {
	pathPattern := regexp.MustCompile(":([a-z]+)")
	matches := pathPattern.FindAllStringSubmatch(endpoint, -1)
	paramKeys := []string{}
	if len(matches) > 0 {
		endpoint = pathPattern.ReplaceAllLiteralString(endpoint, "([^/]+)")
		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}	

	route := route{method, regexp.MustCompile("^" + endpoint + "$"), handler, paramKeys}
	r.routes = append(r.routes, route)
}

func (r *router) GET(pattern string, handler handlerFunc) {
	r.AddRoute(http.MethodGet, pattern, handler)
}

func (r *route) handler(w http.ResponseWriter, req *http.Request) {
	requestString := fmt.Sprint(req.Method, " ", req.URL)
	fmt.Println("received ", requestString)
	start := time.Now()
		customWriter := utils.NewResponseWriter(w)
	r.innerHandler(customWriter, req)
	customWriter.Time= time.Since(start).Milliseconds()
	fmt.Printf("%s resolved within %s \n", requestString, w)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string
	for _, route:= range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if(len(matches) > 0) {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
		}
		if(len(allow) > 0) {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, req)
	}
}
type ContextKey string

func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
  ctx := req.Context()
  for i := 0; i < len(paramKeys); i++ {
    ctx = context.WithValue(ctx, ContextKey(paramKeys[i]), paramValues[i])
  }
  return req.WithContext(ctx)
}