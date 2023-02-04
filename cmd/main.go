package main

import (
	"Forum"
	"Forum/pkg/handlers"
	"Forum/pkg/repository"
	"Forum/pkg/service"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := repository.NewSQLiteDB()
	config := repository.ReadConfig()
	if err != nil {
		log.Fatal("error with read config")
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handler := handlers.NewHandler(services)

	log.Printf("Starting server...\nhttp://localhost%v/\n", ":"+config.Port)

	srv := new(Forum.Server)
	if err := srv.Run(config.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while runnig http server: %s", err.Error())
	}
}
