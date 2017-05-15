package main

import (
	"github.com/lijianying10/api-auth/AppKeyManageServer"
	"github.com/lijianying10/api-auth/ReverseProxyServer"
	"github.com/lijianying10/log"
)

func main() {
	log.Info("Start Auth HTTP Server")
	akms := AppKeyManageServer.New("root:123456@tcp(esm-mysql-normal.jianying.svc.pso.elenet.me:3306)/config")
	rps := ReverseProxyServer.New(akms)

	go func() {
		akms.ListenAndServe(":7776")
	}()
	go func() {
		rps.ListenAndServe(":7777")
	}()
}
