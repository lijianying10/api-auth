package ReverseProxyServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/lijianying10/api-auth/AppKeyManageServer"
)

type ReverseProxyServer struct {
	appKeyManageServer *AppKeyManageServer.AppKeyManageServer
}

func New(akms *AppKeyManageServer.AppKeyManageServer) *ReverseProxyServer {
	return &ReverseProxyServer{
		appKeyManageServer: akms,
	}
}

func (rps *ReverseProxyServer) ListenAndServe(serve string) {
	http.Handle("/", rps.logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = "127.0.0.1:7778"
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})))
	log.Fatal(http.ListenAndServe(serve, nil))

}

func (rps *ReverseProxyServer) logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check key
		if _, ok := r.Header["Date"]; !ok {
			w.WriteHeader(400)
			return
		}
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(400)
			return
		}

		// Checking keys
		sss, err := json.Marshal(r.Header["Authorization"])
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(sss))
		h.ServeHTTP(w, r)
	})
}
