package main

import (
	"fmt"
	"log"
	"net/http"
	"task1/items_manager/pkg/application"
	"task1/items_manager/pkg/configs"
)

func main() {
	port := configs.GetPort()
	addr := ""

	if port == "" {
		port = "4000"
		addr = "localhost:"
	} else {
		addr = "0.0.0.0:"
	}

	app, err := application.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer app.DbClient.Close()
	srv := &http.Server{
		Addr:     addr + port,
		Handler:  app.Router(),
		ErrorLog: app.ErrorLog,
	}
	fmt.Println("server listening on", addr+port)
	log.Fatal(srv.ListenAndServe())
}
