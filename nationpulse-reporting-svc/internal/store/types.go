package store

type ExportResponse struct {
	Status     string
	StatusCode int
	Data       any
}

type Filter struct {
	Headers []string   `json:"headers"`
	Records [][]string `json:"records"`
	Query   string     `json:"query"`
}

type ExportRequest struct {
	ExportID           string `json:"exportID"`
	UserID             int    `json:"userID"`
	Filters            Filter `json:"filters"`
	RequestTableString string `json:"requestTableString"`
	RequestCountryCode string `json:"requestCountryCode"`
}
