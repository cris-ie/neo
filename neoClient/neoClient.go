package neoClient

import (
	"encoding/json"
	"fmt"
	"io"
	"local/neo-api/db"
	"net/http"
	"time"
)

/*
Type for the NEO API Client
  DbCon - The Database connection
  Url - URL for the API endpoint
  ApiKey - The API key to use
*/
type NeoClient struct {
	DbCon  *db.Pgconnector
	Url    string
	ApiKey string
}

/*
Type for configuration of a NeoClient
  Url - URL for the API endpoint
  ApiKey - The API key to use
*/
type NeoClientConfig struct {
	Url    string
	ApiKey string
}

//Internal Type for json parsing
type neoResponse struct {
	Count int                      `json:"element_count,omitempty"`
	Dates map[string][]interface{} `json:"near_earth_objects,omitempty"`
}

// Constants for timestamp calculation
const (
	millisInSecond = 1000
	nsInSecond     = 1000000
	dateLayout     = "2006-01-02"
)

/*
 Create a new instance of a NeoClient from config
  connection - A DB connection to use
  config - configuration for the client
*/
func NewNeoClient(connection *db.Pgconnector, config NeoClientConfig) NeoClient {
	return NeoClient{
		DbCon:  connection,
		Url:    config.Url,
		ApiKey: config.ApiKey,
	}
}

/*
Upsertentries - Inserts entries in the range of from -> to, to the db of the client
from - start time
to - end time
returns the number of inserted entries or an error
*/
func (n NeoClient) UpsertEntries(from time.Time, to time.Time) (int, error) {

	//Create GET Request
	fromStr := from.Format(dateLayout)
	toStr := to.Format(dateLayout)
	response, err := http.Get(fmt.Sprintf("%s?start_date=%s&end_date=%s&api_key=%s", n.Url, fromStr, toStr, n.ApiKey))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		fmt.Println([]byte("{\"status\":500,\"error\":" + err.Error() + "}"))
		return -1, err
	}

	//Parse the JSON response
	data, _ := io.ReadAll(response.Body)
	responseData := neoResponse{}
	json.Unmarshal(data, &responseData)

	//Iterate over response dates
	for _, date := range responseData.Dates {

		//Iterate over NEOS for a given day
		for _, neo := range date {
			//Cast neos to dictionary and access values by name.
			value := db.Neo{}
			value.Id = neo.(map[string]interface{})["id"].(string)
			value.Name = neo.(map[string]interface{})["name"].(string)
			value.NasaJplUrl = neo.(map[string]interface{})["nasa_jpl_url"].(string)
			value.IsPotentiallyHazardousAsteroid = neo.(map[string]interface{})["is_potentially_hazardous_asteroid"].(bool)

			//Actual date is nested in close_approach_data.epoch_date_close_approach as timestamp in millis
			value.Date = fromUnixMilli(int64(neo.(map[string]interface{})["close_approach_data"].([]interface{})[0].(map[string]interface{})["epoch_date_close_approach"].(float64)))

			//Upsert value to DB
			upsert(n.DbCon, value)

		}

	}

	return 0, nil
}

//Upserts a value into the DB
func upsert(connection *db.Pgconnector, entry db.Neo) {

	//Use go-pg/orm to do the upsert
	result, error := connection.Db.Model(&entry).OnConflict("(id) DO UPDATE").
		Insert()

	fmt.Println(result)
	fmt.Println(error)
}

// Convert unix timestamp to go time
func fromUnixMilli(ms int64) time.Time {
	return time.Unix(ms/int64(millisInSecond), (ms%int64(millisInSecond))*int64(nsInSecond))
}
