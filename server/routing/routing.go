package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)


type handlerFunc func(*ResponseWriter, *http.Request)

type route struct {
	method string
	pattern *regexp.Regexp
	innerHandler handlerFunc
	paramKeys []string
	verbose bool
}

type router struct {
	routes []route
}

func NewRouter() *router {
	return &router{ routes: []route{} }
}

func (r *route) handler(w http.ResponseWriter, req *http.Request) {
	requestString := fmt.Sprint(req.Method, " ", req.URL)
	if r.verbose {
		fmt.Println("Recieved ", requestString)
	}
	start := time.Now()
	nw := NewResponseWriter(w)
	r.innerHandler(nw, req)
	nw.Time = time.Since(start).Milliseconds()
	if r.verbose {
		fmt.Printf("%s resolved with %s\n", requestString, nw)
	}
}

func (r *router) addRoute(method, endpoint string, handler handlerFunc, verbose bool) {
	pathParamPattern := regexp.MustCompile(":([a-z]+)")
	matches := pathParamPattern.FindAllStringSubmatch(endpoint, -1)
	paramKeys := []string{}

	if len(matches) > 0 {
		endpoint = pathParamPattern.ReplaceAllLiteralString(endpoint, "([^/]+)")

		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}

	route := route { method, regexp.MustCompile("^" + endpoint + "$"), handler, paramKeys, verbose }
	r.routes = append(r.routes, route)
}

func (r *router) GET(pattern string, handler handlerFunc, verbose bool) {
	r.addRoute(http.MethodGet, pattern, handler, verbose)
}

func (r *router) POST(pattern string, handler handlerFunc, verbose bool) {
	r.addRoute(http.MethodPost, pattern, handler, verbose)
}

func (r *router) DELETE(pattern string, handler handlerFunc, verbose bool) {
	r.addRoute(http.MethodDelete, pattern, handler, verbose)
}

func (r *router) PUT(pattern string, handler handlerFunc, verbose bool) {
	r.addRoute(http.MethodPut, pattern, handler, verbose)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string
	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)

		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			route.handler(w, buildContext(req, route.paramKeys, matches[1:]))

			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, req)
}

type ContextKey string

func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
	ctx := req.Context()

	for i := 0; i < len(paramKeys); i++ {
		ctx = context.WithValue(ctx, ContextKey(paramKeys[i]), paramValues[i])
	}

	return req.WithContext(ctx)
}

type ResponseWriter struct {
	Status int
	Body string
	Time int64
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter { ResponseWriter: w }
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	w.Body = string(body)
	return w.ResponseWriter.Write(body)
}

func (w *ResponseWriter) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)

	if w.Body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}

	return out
}

func (w *ResponseWriter) StringResponse(code int, response string) {
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func (w *ResponseWriter) JSONResponse(code int, responseObject any) {
	w.WriteHeader(code)
	response, err := json.Marshal(responseObject)

	if err != nil {
		w.StringResponse(http.StatusBadRequest, err.Error())
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
