// Package repository holds the project's database access layer.
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// ScheduleRecord is one booked time window row.
type ScheduleRecord struct {
	StationId    int
	DatetimeFrom time.Time
	DatetimeTo   time.Time
}

// ScheduleRepository persists and queries booked charging-station time windows.
type ScheduleRepository struct {
	db *sql.DB
}

// NewScheduleRepository builds a ScheduleRepository backed by db.
func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// GetTimeWindowsByDateAndStationIDs returns every booked window on date for any of stationIds.
func (r *ScheduleRepository) GetTimeWindowsByDateAndStationIDs(
	date time.Time,
	stationIds []int,
) (records []ScheduleRecord, err error) {
	rows, err := r.db.Query(
		`SELECT station_id, datetime_from, datetime_to
		 FROM schedule
		 WHERE station_id = ANY($1) AND datetime_from::DATE = $2::DATE`,
		pq.Array(stationIds), date,
	)
	if err != nil {
		return nil, fmt.Errorf("query time windows: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close time window rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		var record ScheduleRecord
		if err := rows.Scan(&record.StationId, &record.DatetimeFrom, &record.DatetimeTo); err != nil {
			return nil, fmt.Errorf("scan time window row: %w", err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate time window rows: %w", err)
	}

	return records, nil
}

// BeginTx starts a transaction for a batch of SaveTimeWindow calls.
func (r *ScheduleRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

// SaveTimeWindow inserts one booked time window within tx.
func (r *ScheduleRepository) SaveTimeWindow(
	tx *sql.Tx,
	code string,
	stationId int,
	datetimeFrom time.Time,
	datetimeTo time.Time,
) error {
	_, err := tx.Exec(
		`INSERT INTO schedule (station_id, code, datetime_from, datetime_to)
		 VALUES ($1, $2, $3, $4)`,
		stationId, code, datetimeFrom, datetimeTo,
	)
	if err != nil {
		return fmt.Errorf("insert time window: %w", err)
	}

	return nil
}
