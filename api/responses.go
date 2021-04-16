package api

//Response type for the /neo/week query
type WeekResponse struct {
	Count int    `json:"count"`
	From  string `json:"from"`
	To    string `json:"to"`
}

//General json status response
type StatusResponse struct {
	Status int `json:"status"`
}
