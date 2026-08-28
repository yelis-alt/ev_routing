package controller

import (
	"encoding/json"
	"net/http"

	"ev_routing/internal/dto"
	"ev_routing/internal/service"
)

// ScheduleController exposes HTTP endpoints for charging-station booking
// time windows.
type ScheduleController struct {
	schedule *service.ScheduleService
}

// NewScheduleController builds a ScheduleController backed by schedule.
func NewScheduleController(schedule *service.ScheduleService) *ScheduleController {
	return &ScheduleController{schedule: schedule}
}

// RegisterRoutes wires the controller's endpoints onto mux.
func (sc *ScheduleController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /schedule/getTimeWindows", sc.handleGetTimeWindows)
	mux.HandleFunc("POST /schedule/saveTimeWindows", sc.handleSaveTimeWindows)
}

// handleGetTimeWindows returns each requested station's booked time windows.
func (sc *ScheduleController) handleGetTimeWindows(w http.ResponseWriter, r *http.Request) {
	var timeWindowsRequest dto.TimeWindowsRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&timeWindowsRequest); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	timeWindowsList, err := sc.schedule.GetTimeWindows(timeWindowsRequest)
	if err != nil {
		http.Error(w, "failed to get time windows: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, timeWindowsList)
}

// handleSaveTimeWindows books the given stations' time windows.
func (sc *ScheduleController) handleSaveTimeWindows(w http.ResponseWriter, r *http.Request) {
	var timeWindowsSaveRequestsList []dto.TimeWindowsSaveRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&timeWindowsSaveRequestsList); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.schedule.SaveTimeWindows(timeWindowsSaveRequestsList); err != nil {
		http.Error(w, "failed to save time windows: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, dto.TimeWindowsSaveOutputDTO{Status: http.StatusOK})
}
