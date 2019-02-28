package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"go.uber.org/zap"

	"github.com/makeroo/taxi_scout/rest_backend"
	"github.com/makeroo/taxi_scout/storage"
)

func main() {
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	port := flag.Int("http-port", 8008, "HTTP port to listen to")
	secure := flag.Bool("secure", false, "Publish secure cookies (requires https)")

	dao, err := storage.NewSqlDatastore("mysql", "taxi_scout_user:taxi_scout_pwd@/taxi_scout?parseTime=true", logger)
	defer dao.Close()

	if err != nil {
		panic(err.Error())
	}

	// Hash keys should be at least 32 bytes long
	hashKey := []byte("0123456789 123456789 123456789 1") // FIXME: read from protected file
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	blockKey := []byte("0123456789 123456789 123456789 1") // FIXME: read from protected file

	secureCookies := securecookie.New(hashKey, blockKey)

	server := &rest_backend.RestServer{
		Dao:    dao,
		Logger: logger,
		Configuration: rest_backend.Configuration{
			SecureCookies: secureCookies,
			HttpsCookies: *secure,
		},
	}

	r := mux.NewRouter()
//	r.HandleFunc("/invitation/{token:.+}", rest_backend.DisableBrowserCache(server.Invitation))
	r.HandleFunc("/accounts", rest_backend.DisableBrowserCache(server.Accounts))
	r.HandleFunc("/account/{id:[0-9]+|me}", rest_backend.DisableBrowserCache(server.Account))
	r.HandleFunc("/accounts/authenticate", rest_backend.DisableBrowserCache(server.AccountsAuthenticate))
	r.HandleFunc("/account/{id:[0-9]+}/password", rest_backend.DisableBrowserCache(server.AccountPassword))
	r.HandleFunc("/account/{id:[0-9]+}/groups", rest_backend.DisableBrowserCache(server.AccountGroups))
	r.HandleFunc("/account/{id:[0-9]+}/scouts", rest_backend.DisableBrowserCache(server.AccountScouts))
//	r.HandleFunc("/scouts", rest_backend.DisableBrowserCache(server.Scouts))
	r.HandleFunc("/scout/{id:[0-9]+}", rest_backend.DisableBrowserCache(server.Scout))
	http.Handle("/", r)

	//fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	addr := fmt.Sprintf(":%d", *port)
	// TODO: use logging
	logger.Infow("http server setup",
		"port", *port,
	)

	err = http.ListenAndServe(addr, nil)

	if err != nil {
		logger.DPanic(err.Error())
	}
}
