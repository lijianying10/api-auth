package AppKeyManageServer

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/lijianying10/api-auth/AuthKey"
	"github.com/lijianying10/log"
)

type AppKeyManageServer struct {
	connStr string
	keys    map[string]AuthKey.AuthKey
	db      *sqlx.DB
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func New(dbConn string) AppKeyManageServer {
	rand.Seed(time.Now().UnixNano())
	return &AppKeyManageServer{
		connStr: dbConn,
		db:      sqlx.MustConnect("mysql", dbConn),
	}
}

func (akms *AppKeyManageServer) ListenAndServe() {
	router := httprouter.New()
	router.POST("/auth/key", akms.HandlerNewAppKey)
	router.GET("/auth/key", akms.HandlerGetAppKeys)
	router.Option("/auth/key", akms.HandlerRefreashAppKeys)
}

func (akms *AppKeyManageServer) HandlerNewAppKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	headers := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(headers); err != nil {
		w.WriteHeader(400)
		return
	}

	// Make sure headers with no indent when insert to database
	jsonHeader, err := json.Marshal(headers)
	if err != nil {
		log.Error("error json marshal", err.Error())
		w.WriteHeader(500)
		return
	}

	Appkey := randStringRunes(32)
	AppSecret := randStringRunes(32)

	_, err := akms.db.NamedExec("INSERT INTO `ApiAuthKey`(`AppKey`,`AppSecret`,`Headers`) VALUES (:AppKey, :AppSecret, :headers)",
		map[string]interface{}{
			"AppKey":    Appkey,
			"AppSecret": AppSecret,
			"headers":   string(headers),
		})
	if err != nil {
		log.Error("error database storage : ", err.Error())
		w.WriteHeader(500)
		return
	}

	akms.Refreash()

	resData, err := json.Marshal(map[string]interface{}{
		"message": "",
		"data": map[string]string{
			"AppKey":    AppKey,
			"AppSecret": AppSecret,
		},
	})
	w.Write(resData)
	w.WriteHeader(200)
	return
}

func (akms *AppKeyManageServer) HandlerGetAppKeys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
func (akms *AppKeyManageServer) HandlerRefreashAppKeys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	akms.Refreash()
	w.WriteHeader(200)
	return
}

func (akms *AppKeyManageServer) Refreash() error {
	type ApiAuthMapping struct {
		AppKey    string `db:"AppKey"`
		AppSecret string `db:"AppSecret"`
		headers   string `db:"headers"`
	}
	keys := []ApiAuthMapping{}
	err := akms.db.Get(&keys, "select AppKey,AppSecret,Headers from ApiAuthKey")
	if err != nil {
		log.Error("error Refreash data", err.Error())
		w.WriteHeader(500)
		return err
	}
	akms.keys = make(map[string]AuthKey.AuthKey)
	for _, key := range keys {
		headers := make(map[string]string)
		err := json.Unmarshal(key.headers, &headers)
		if err != nil {
			log.Error("header format error, "+key.AppKey+" :", err.Error())
			continue
		}
		akms.keys[key.AppKey] = AuthKey.AuthKey{
			Appkey:    key.AppKey,
			AppSecret: key.AppSecret,
			headers:   headers,
		}
	}

	return nil
}

func (akms *AppKeyManageServer) Get(AppKey string) AuthKey.AuthKey {
	return akms.keys[Appkey]
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
