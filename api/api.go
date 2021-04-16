package api

import (
	"context"
	"fmt"
	"local/neo-api/db"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

/*
Type for the API Server
DbCon - Connection to the DB
*/
type ApiServer struct {
	DbCon db.Pgconnector
}

/*
Create a new ApiServer instance
connection - The DB connection to use
*/
func NewApiServer(connection db.Pgconnector) *ApiServer {
	return &ApiServer{
		DbCon: connection,
	}
}

// Handler for /Liveness route - reports the liveness
// Always returns HTTP 200 and 200 in the body
func (server *ApiServer) Liveness(ctx echo.Context) error {
	status := StatusResponse{
		Status: 200,
	}
	return ctx.JSON(http.StatusOK, status)
}

//Handler for /neo/week route  - gets neos for next week
func (server *ApiServer) Week(ctx echo.Context) error {
	layout := "2006-01-02"
	today := time.Now().UTC().Format(layout)
	nextWeek := time.Now().UTC().AddDate(0, 0, 7).Format(layout)
	var neos []db.Neo

	//Get neos for the next 7 days from the DB
	error := server.DbCon.Db.Model(&neos).Where(
		fmt.Sprintf("date BETWEEN '%s' AND '%s'", today, nextWeek),
	).Select()
	if error != nil {
		return ctx.JSON(http.StatusInternalServerError, error.Error())
	}

	//Return the number of neos in the next 7 days and the period queried
	ans := WeekResponse{
		Count: len(neos),
		From:  today,
		To:    nextWeek,
	}

	return ctx.JSON(http.StatusOK, ans)

}

//Handler for /neo/next route  - gets the next (hazardous) objects
func (server *ApiServer) Next(ctx echo.Context) error {

	//Get query paramater hazardous
	hazardous := ctx.QueryParam("hazardous")
	whereClause := "date > now()"

	// Amend where clause if hazardous was passed and looks like a bool
	if len(hazardous) > 0 && checkBool(hazardous) {
		whereClause += fmt.Sprintf(" AND hazardous = %s", hazardous)
	}

	var neos []db.Neo
	error := server.DbCon.Db.Model(&neos).
		Where(whereClause).
		Order("date").
		Limit(1).
		Select()
	if error != nil {
		return ctx.JSON(http.StatusInternalServerError, error.Error())
	}
	if len(neos) > 0 {
		return ctx.JSON(http.StatusOK, neos[0])
	} else {
		return ctx.JSON(http.StatusOK, nil)
	}

}

// Helperfunction to check if the string looks like a bool
func checkBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

//Handler for /Liveness status  - reports status of the server (e.g is the db connected)
func (server *ApiServer) Status(ctx echo.Context) error {

	//dummy context for Ping()
	ctx2 := context.Background()

	if err := server.DbCon.Db.Ping(ctx2); err != nil {
		status := StatusResponse{
			Status: http.StatusServiceUnavailable,
		}
		return ctx.JSON(http.StatusServiceUnavailable, status)
	}

	status := StatusResponse{
		Status: 200,
	}
	return ctx.JSON(http.StatusOK, status)

}
