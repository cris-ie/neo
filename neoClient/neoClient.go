package neoClient

import (
	"encoding/json"
	"fmt"
	"io"
	"local/neo-api/db"
	"net/http"
	"time"
)

type NeoClient struct {
	DbCon  *db.Pgconnector
	Url    string
	ApiKey string
}
type NeoClientConfig struct {
	Url    string
	ApiKey string
}

type neoResponse struct {
	Count int                      `json:"element_count,omitempty"`
	Dates map[string][]interface{} `json:"near_earth_objects,omitempty"`
}
type neoResponseInner struct {
	id string
}

const NASA_API_FEED = "https://api.nasa.gov/neo/rest/v1/feed"
const millisInSecond = 1000
const nsInSecond = 1000000
const dateLayout = "2006-01-02"

func NewNeoClient(connection *db.Pgconnector, config NeoClientConfig) NeoClient {
	return NeoClient{
		DbCon:  connection,
		Url:    config.Url,
		ApiKey: config.ApiKey,
	}
}

//Upsertentries - Inserts entries from -> to the db of the client
//returns  : number of inserted entries or an error
func (n NeoClient) UpsertEntries(from time.Time, to time.Time) (int, error) {

	fromStr := from.Format(dateLayout)
	toStr := to.Format(dateLayout)
	response, err := http.Get(fmt.Sprintf("%s?start_date=%s&end_date=%s&api_key=%s", n.Url, fromStr, toStr, n.ApiKey))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		fmt.Println([]byte("{\"status\":500,\"error\":" + err.Error() + "}"))
		return -1, err
	}

	data, _ := io.ReadAll(response.Body)
	myVar := neoResponse{}
	json.Unmarshal(data, &myVar)
	for _, date := range myVar.Dates {

		for _, neo := range date {

			foo := db.Neo{}
			foo.Id = neo.(map[string]interface{})["id"].(string)
			foo.Name = neo.(map[string]interface{})["name"].(string)
			foo.NasaJplUrl = neo.(map[string]interface{})["nasa_jpl_url"].(string)
			foo.IsPotentiallyHazardousAsteroid = neo.(map[string]interface{})["is_potentially_hazardous_asteroid"].(bool)
			foo.Date = fromUnixMilli(int64(neo.(map[string]interface{})["close_approach_data"].([]interface{})[0].(map[string]interface{})["epoch_date_close_approach"].(float64)))

			upsert(n.DbCon, foo)

		}

	}

	return 0, nil
}

func upsert(connection *db.Pgconnector, entry db.Neo) {
	result, error := connection.Db.Model(&entry).OnConflict("(id) DO UPDATE").
		Insert()

	fmt.Println(result)
	fmt.Println(error)
}

func fromUnixMilli(ms int64) time.Time {
	return time.Unix(ms/int64(millisInSecond), (ms%int64(millisInSecond))*int64(nsInSecond))
}
