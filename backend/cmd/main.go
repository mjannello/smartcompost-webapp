package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mjannello/smartcompost-webapp/backend/config"
	"github.com/mjannello/smartcompost-webapp/backend/internal/http"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node/adapter/repository"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node/port"
)

func main() {
	cfg, err := config.LoadConfig("./cfg/cfg.yaml")
	if err != nil {
		log.Fatalf("Error loading cfg: %v", err)
	}

	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbUser := cfg.Database.User
	dbPassword := cfg.Database.Password
	dbName := cfg.Database.Name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	database, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.Close()

	nodeRepo := repository.NewMySQLNodeRepository(database)
	nodeApp := app.NewNodeService(nodeRepo)
	router := mux.NewRouter()

	port.RegisterRoutes(router)
	http.RegisterRoutes(router)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
