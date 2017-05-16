package ReverseProxyServer

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/lijianying10/api-auth/AppKeyManageServer"
	"github.com/lijianying10/log"
)

type ReverseProxyServer struct {
	appKeyManageServer *AppKeyManageServer.AppKeyManageServer
}

func New(akms *AppKeyManageServer.AppKeyManageServer) *ReverseProxyServer {
	return &ReverseProxyServer{
		appKeyManageServer: akms,
	}
}

func (rps *ReverseProxyServer) ListenAndServe(serve string, backendServer string) {
	http.Handle("/", rps.logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = backendServer
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
		keysSequence := strings.Split(r.Header["Authorization"][0], ":")
		exist, authKey := rps.appKeyManageServer.Get(keysSequence[0])
		if !exist {
			w.WriteHeader(403)
			w.Write([]byte("Bad AppKey"))
			return
		}

		valid, err := authKey.CheckSignature(keysSequence[1], r.Method, r.Header["Date"][0], strings.Split(r.RequestURI, "?")[0])
		if err != nil {
			log.Error("Check Signature error: ", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("Auth Server ERROR"))
			return
		}

		if !valid {
			w.WriteHeader(403)
			w.Write([]byte("Bad Signature"))
			return
		}

		h.ServeHTTP(w, r)
	})
}
