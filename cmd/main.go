package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"task1/items_manager/pkg/application"
)

func main() {
	addr := flag.String("addr", ":4000", "Server address")
	flag.Parse()

	app, err := application.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer app.DbClient.Close()
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.Router(),
		ErrorLog: app.ErrorLog,
	}
	fmt.Println("server listening on", *addr)
	log.Fatal(srv.ListenAndServe())
}
