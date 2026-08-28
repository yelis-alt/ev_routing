package schedule

import (
	"fmt"
	"strings"
	"time"

	"ev_routing/internal/dto"
	"ev_routing/internal/repository"
)

const (
	dateLayout     = "02.01.2006"
	clockLayout    = "15:04"
	dateTimeLayout = dateLayout + " " + clockLayout
)

// ScheduleService manages charging-station booking time windows.
type ScheduleService struct {
	repo *repository.ScheduleRepository
}

// NewScheduleService builds a ScheduleService backed by repo.
func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

// GetTimeWindows returns, for every requested station, its booked time
// windows on the requested date (empty when none are booked).
func (s *ScheduleService) GetTimeWindows(request dto.TimeWindowsRequestDTO) ([]dto.TimeWindowsOutputDTO, error) {
	date, err := time.Parse(dateLayout, request.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date %q: %w", request.Date, err)
	}

	records, err := s.repo.GetTimeWindowsByDateAndStationIDs(date, request.StationIdsList)
	if err != nil {
		return nil, err
	}

	timeWindowsByStation := make(map[int][]string, len(request.StationIdsList))
	for _, stationId := range request.StationIdsList {
		timeWindowsByStation[stationId] = []string{}
	}
	for _, record := range records {
		timeWindowsByStation[record.StationId] = append(
			timeWindowsByStation[record.StationId],
			formatTimeWindow(record.DatetimeFrom, record.DatetimeTo),
		)
	}

	timeWindowsOutputList := make([]dto.TimeWindowsOutputDTO, 0, len(request.StationIdsList))
	for _, stationId := range request.StationIdsList {
		timeWindowsOutputList = append(timeWindowsOutputList, dto.TimeWindowsOutputDTO{
			StationId:       stationId,
			TimeWindowsList: timeWindowsByStation[stationId],
		})
	}

	return timeWindowsOutputList, nil
}

// SaveTimeWindows books every station/window pair in requestsList as a
// single transaction: either all windows are saved, or none are.
func (s *ScheduleService) SaveTimeWindows(requestsList []dto.TimeWindowsSaveRequestDTO) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	for _, request := range requestsList {
		for _, timeWindow := range request.TimeWindowsList {
			datetimeFrom, datetimeTo, err := parseTimeWindow(timeWindow)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("parse time window %q: %w", timeWindow, err)
			}

			if err := s.repo.SaveTimeWindow(tx, request.Code, request.StationId, datetimeFrom, datetimeTo); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// formatTimeWindow renders "dd.MM.yyyy HH:MM-HH:MM".
func formatTimeWindow(datetimeFrom, datetimeTo time.Time) string {
	return datetimeFrom.Format(dateTimeLayout) + "-" + datetimeTo.Format(clockLayout)
}

// parseTimeWindow parses "dd.MM.yyyy HH:MM-HH:MM".
func parseTimeWindow(timeWindow string) (time.Time, time.Time, error) {
	dateAndPeriod := strings.SplitN(timeWindow, " ", 2)
	if len(dateAndPeriod) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf(`expected "<date> <from>-<to>", got %q`, timeWindow)
	}
	date := dateAndPeriod[0]

	period := strings.SplitN(dateAndPeriod[1], "-", 2)
	if len(period) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf(`expected "<from>-<to>", got %q`, dateAndPeriod[1])
	}

	datetimeFrom, err := time.Parse(dateTimeLayout, date+" "+period[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse start time: %w", err)
	}

	datetimeTo, err := time.Parse(dateTimeLayout, date+" "+period[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse finish time: %w", err)
	}

	return datetimeFrom, datetimeTo, nil
}
