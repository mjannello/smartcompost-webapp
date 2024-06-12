package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	http2 "github.com/mjannello/smartcompost-webapp/backend/internal/http"
	measurementrepo "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/adapter/repository"
	measurementapp "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/app"
	measurementport "github.com/mjannello/smartcompost-webapp/backend/internal/measurement/port"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mjannello/smartcompost-webapp/backend/config"
	noderepository "github.com/mjannello/smartcompost-webapp/backend/internal/node/adapter/repository"
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	nodeport "github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
)

func main() {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading cfg: %v", err)
	}

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPassword := cfg.Database.Password
	dbName := cfg.Database.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.Close()

	nodeRepo := noderepository.NewMySQL(database)
	nodeService := nodeapp.NewNodeService(nodeRepo)
	nodeHandler := nodeport.NewHTTPHandler(nodeService)

	measurementRepo := measurementrepo.NewMySQL(database)
	measurementService := measurementapp.NewService(measurementRepo)
	measurementHandler := measurementport.NewHTTPHandler(measurementService)

	router := mux.NewRouter()
	routerHandler := http2.NewRouterHandler(nodeHandler, measurementHandler)
	routerHandler.RouteURLs(router)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
