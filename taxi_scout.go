package main

import (
	"flag"
	"fmt"
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

	dao, err := storage.NewSqlDatastore("mysql", "go_rest_test:go_rest_pwd@/go_rest_test", logger)
	defer dao.Close()

	if err != nil {
		panic(err.Error())
	}

	server := &rest_backend.RestServer{
		Dao:    dao,
		Logger: logger,
	}

	r := mux.NewRouter()
	r.HandleFunc("/accounts", server.Accounts)
	r.HandleFunc("/account/{id:[0-9]+}", server.Account)
	r.HandleFunc("/accounts/authenticate", server.AccountsAuthenticate)
	r.HandleFunc("/account/{id:[0-9]+}/password", server.AccountPassword)
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
