package application

import (
	"log"
	"os"
	"task1/items_manager/pkg/configs"
	"task1/items_manager/pkg/db_client"
	"time"
)

type Application struct {
	infoLog  *log.Logger
	ErrorLog *log.Logger
	DbClient *db_client.Client
	*sessionManager
}

func Get() (*Application, error) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db_client, err := db_client.Get(15*time.Second, configs.GetMongoURI())

	if err != nil {
		return nil, err
	}
	var sess sessionManager
	sess.keyToUser = make(map[string]session)
	sess.userToKey = make(map[string]session)
	return &Application{
		infoLog:        infoLog,
		ErrorLog:       errorLog,
		DbClient:       db_client,
		sessionManager: &sess,
	}, nil
}
