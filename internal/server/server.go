package server

import (
	"net/http"

	"github.com/tera-language/teralang/internal/logger"
	"github.com/tera-language/teralang/internal/parser"
)

func Server(routes []parser.Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == route.Method || (r.Method == "" && route.Method == "GET") {
				logger.Infoln(route.Path, route.Method, route.Status)
				w.WriteHeader(route.Status)
				for k, v := range route.Headers {
					w.Header().Set(k, v)
				}
				_, err := w.Write([]byte(route.Body))
				if err != nil {
					logger.Warningln(err)
				}
			}
		})
	}

	return mux
}
