CREATE TABLE IF NOT EXISTS schedule (
    id SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL,
    code VARCHAR(255) NOT NULL,
    datetime_from TIMESTAMP NOT NULL,
    datetime_to TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_schedule_station_date ON schedule (station_id, datetime_from);
