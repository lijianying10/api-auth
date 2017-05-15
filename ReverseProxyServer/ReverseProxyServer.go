package ReverseProxyServer

import (
	"net/http"
	"net/http/httputil"

	"github.com/lijianying10/api-auth/AppKeyManageServer"
	"github.com/lijianying10/log"
)

type ReverseProxyServer struct {
	appKeyManageServer *AppKeyManageServer.AppKeyManageServer
}

func New(akms *AppKeyManageServer) *ReverseProxyServer {
	return &ReverseProxyServer{
		appKeyManageServer: akms,
	}
}

func (rps *ReverseProxyServer) ListenAndServe(serve string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// check key
		if val, ok := r.Header["Date"]; !ok {
			log.Info("code here")
			w.WriteHeader(400)
			return
		}
		if val, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(400)
			return
		}

		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = r.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(serve, nil))
}
