package dto

// TimeWindowsRequestDTO is the request body for POST /schedule/getTimeWindows.
type TimeWindowsRequestDTO struct {
	Date           string `json:"date"` // dd.MM.yyyy
	StationIdsList []int  `json:"stationIdsList"`
}

// TimeWindowsOutputDTO lists one station's booked time windows for the
// requested date, each formatted as "dd.MM.yyyy HH:MM-HH:MM".
type TimeWindowsOutputDTO struct {
	StationId       int      `json:"stationId"`
	TimeWindowsList []string `json:"timeWindowsList"`
}

// TimeWindowsSaveRequestDTO is one station's booking in the request body for
// POST /schedule/saveTimeWindows.
type TimeWindowsSaveRequestDTO struct {
	StationId       int      `json:"stationId"`
	Code            string   `json:"code"`
	TimeWindowsList []string `json:"timeWindowsList"` // "dd.MM.yyyy HH:MM-HH:MM"
}

// TimeWindowsSaveOutputDTO reports the outcome of a saveTimeWindows call.
type TimeWindowsSaveOutputDTO struct {
	Status int `json:"status"`
}
