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
	serialnumbergeneratorapp "github.com/mjannello/smartcompost-webapp/backend/internal/serial_number_generator/app"
	"github.com/mjannello/smartcompost-webapp/backend/internal/serial_number_generator/generators"
	"github.com/mjannello/smartcompost-webapp/backend/pkg/clock"
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

	// [from here] if we want to change the log file
	//logFilePath := "/var/log/web/app.log"
	//
	//logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("Error abriendo el archivo de log: %v", err)
	//}
	//defer logFile.Close()
	//
	//log.SetOutput(logFile)
	// to here
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

	realClock := clock.NewClock()
	uuidGenerator := generators.NewUUIDGenerator()
	serialNumberGeneratorService := serialnumbergeneratorapp.NewSerialNumberGeneratorService(uuidGenerator)

	nodeRepo := noderepository.NewNodeRepository(database)
	nodeService := nodeapp.NewNodeService(nodeRepo, serialNumberGeneratorService, realClock)
	nodeHandler := nodeport.NewNodeHandler(nodeService)

	measurementRepo := measurementrepo.NewMeasurementRepository(database)
	measurementService := measurementapp.NewMeasurementService(measurementRepo, nodeService)
	measurementHandler := measurementport.NewMeasurementHandler(measurementService)

	router := mux.NewRouter()
	routerHandler := http2.NewRouterHandler(nodeHandler, measurementHandler)
	routerHandler.RouteURLs(router)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
