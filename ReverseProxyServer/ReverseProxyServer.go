package ReverseProxyServer

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ReverseProxyServer struct {
	appKeyManageServer *AppKeyManageServer
}

func New(akms *AppKeyManageServer) *ReverseProxyServer {
	return &ReverseProxyServer{
		appKeyManageServer: akms,
	}
}

func (rps *ReverseProxyServer) ListenAndServe(serve string) {}

func (rps *ReverseProxyServer) Handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
