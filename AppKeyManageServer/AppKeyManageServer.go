package AppKeyManageServer

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lijianying10/api-auth/AuthKey"
)

type AppKeyManageServer struct {
	connStr string
	keys    map[string]AuthKey.AuthKey
}

func New(dbConn string) AppKeyManageServer {
	return &AppKeyManageServer{
		connStr: dbConn,
	}
}

func (akms *AppKeyManageServer) ListenAndServe() error {
	return nil
}

func (akms *AppKeyManageServer) HandlerNewAppKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
func (akms *AppKeyManageServer) HandlerGetAppKeys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
func (akms *AppKeyManageServer) HandlerRefreashAppKeys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (akms *AppKeyManageServer) Refreash() error {
	return nil
}

func (akms *AppKeyManageServer) Get(AppKey string) (AppSecret string, err error) {
	return "", nil
}
