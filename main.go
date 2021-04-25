package main

import (
	"fmt"
	"local/neo-api/api"
	"local/neo-api/db"
	"local/neo-api/neoClient"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// Default webserver port
const port = "8090"

func main() {

	// Set DB port to ENV value or default. TODO do this for all arguments where applicable
	dbport := os.Getenv("DB_PORT")
	if len(dbport) == 0 {
		dbport = "5432"
	}

	// Generate configuration object for use in clients (load values from ENV)
	pgConfig := db.DbConfig{
		DbName:       os.Getenv("DB_NAME"),
		UserName:     os.Getenv("DB_USER"),
		UserPassword: os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		Port:         dbport,
	}
	neoConfig := neoClient.NeoClientConfig{
		Url:    "https://api.nasa.gov/neo/rest/v1/feed",
		ApiKey: os.Getenv("NASA_KEY"),
	}

	// Create clients for db connection and NASA JPL NEO API connection
	pgClient := db.CreateConnectClient(pgConfig)
	neoClient := neoClient.NewNeoClient(&pgClient, neoConfig)

	// Create new web server instance
	echoWebServer := echo.New()

	//Create DB Schema in selected DB
	err := pgClient.CreateSchemaIfNotExists()
	if err != nil {
		panic("Create schema failed: " + err.Error())
	}

	//Setup handler object for web server functions
	apiServer := api.NewApiServer(pgClient)

	//Seed the database
	from := time.Now()
	to := time.Now().AddDate(0, 0, 7)
	neoClient.UpsertEntries(from, to)

	//Remember to close the DB connection on program end
	defer pgClient.Close()

	fmt.Println("Starting application")

	//Add routes, actions are handled by the apiServer object created above
	echoWebServer.Router().Add("GET", "/liveness", apiServer.Liveness)
	echoWebServer.Router().Add("GET", "/status", apiServer.Status)
	echoWebServer.Router().Add("GET", "/neo/week", apiServer.Week)
	echoWebServer.Router().Add("GET", "/neo/next", apiServer.Next)

	// Serve file request for favicon
	echoWebServer.File("/favicon.ico", "/favicon.ico")

	fmt.Println("Starting to listen on Port:" + port)
	errStart := echoWebServer.Start(fmt.Sprintf(":%v", port))
	if errStart != nil {
		panic("Error during echo Start(): " + err.Error())
	}

}
