package main

import (
	"fmt"
	"local/neo-api/api"
	"local/neo-api/db"
	neoclient "local/neo-api/neo-client"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const port = "8090"

func getmyenv(s string) string {
	return os.Getenv(s)
}
func main() {

	pgConfig := db.DbConfig{DbName: "neo", UserName: "crisie", UserPassword: "Lus", Host: "localhost", Port: "5432"}
	neoConfig := neoclient.NeoClientConfig{
		Url:    "https://api.nasa.gov/neo/rest/v1/feed",
		ApiKey: os.Getenv("NASA_KEY"),
	}
	pgClient := db.CreateConnectClient(pgConfig)
	neoclient := neoclient.NewNeoClient(&pgClient, neoConfig)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	err := pgClient.CreateSchemaIfNotExists()
	if err != nil {
		panic("Create schema failed: " + err.Error())
	}
	apiServer := api.NewApiServer(pgClient)
	from := time.Now()
	to := time.Now().AddDate(0, 0, 1)
	neoclient.UpsertEntries(from, to)
	defer pgClient.Close()
	fmt.Println("Starting application")

	e.Router().Add("GET", "/liveness", apiServer.Liveness)
	e.Router().Add("GET", "/status", apiServer.Status)
	e.Router().Add("GET", "/neo/week", apiServer.Week)
	e.Router().Add("GET", "/neo/next", apiServer.Next)
	e.File("/favicon.ico", "/favicon.ico")

	go func() {

		fmt.Println("Starting to listen on Port:" + port)
		err := e.Start(fmt.Sprintf(":%v", port))
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	select {}

}
