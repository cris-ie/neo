package api

type WeekResponse struct {
	Count int    `json:"count"`
	From  string `json:"from"`
	To    string `json:"to"`
}
