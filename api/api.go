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

const NASA_API_FEED = "https://api.nasa.gov/neo/rest/v1/feed"
const internalErrorMessage = "internal server error"

type ApiServer struct {
	DbCon db.Pgconnector
}

func NewApiServer(connection db.Pgconnector) *ApiServer {
	return &ApiServer{
		DbCon: connection,
	}
}

//Liveness - reports the liveness
func (server *ApiServer) Liveness(ctx echo.Context) error {
	status := 200
	return ctx.JSON(http.StatusOK, status)
}

//Week - gets neos for next week
func (server *ApiServer) Week(ctx echo.Context) error {
	layout := "2006-01-02"
	today := time.Now().UTC().Format(layout)
	nextWeek := time.Now().UTC().AddDate(0, 0, 7).Format(layout)
	var neos []db.Neo

	error := server.DbCon.Db.Model(&neos).Where(
		fmt.Sprintf("date BETWEEN '%s' AND '%s'", today, nextWeek),
	).Select()
	if error != nil {
		return ctx.JSON(http.StatusInternalServerError, error.Error())

	}
	ans := WeekResponse{
		Count: len(neos),
		From:  today,
		To:    nextWeek,
	}

	return ctx.JSON(http.StatusOK, ans)

}

//Next - gets the next (hazardous) objects
func (server *ApiServer) Next(ctx echo.Context) error {

	hazardous := ctx.QueryParam("hazardous")
	whereClause := "date > now()"
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
func checkBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

//Status - reports status of the server
func (server *ApiServer) Status(ctx echo.Context) error {
	ctx2 := context.Background()
	if err := server.DbCon.Db.Ping(ctx2); err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, "not ready")
	}

	return ctx.JSON(http.StatusOK, "ready")

}
