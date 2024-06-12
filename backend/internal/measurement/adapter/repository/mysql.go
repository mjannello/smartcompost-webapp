package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

const (
	GetAllMeasurementsByNodeID = "SELECT id, node_id, value, timestamp, type FROM measurements WHERE node_id = ?"
	InsertNewMeasurement       = "INSERT INTO measurements (node_id, value, timestamp, type) VALUES (?, ?, ?, ?)"
)

type mySQL struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) measurement.Repository {
	return &mySQL{db: db}
}

func (m *mySQL) GetAllMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurement.Measurement, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, GetAllMeasurementsByNodeID, nodeID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("could not get measurements from DB: %w", err)
	}
	defer rows.Close()

	var measurements []measurement.Measurement
	for rows.Next() {
		var m measurement.Measurement
		var timestampStr string
		err = rows.Scan(&m.ID, &m.NodeID, &m.Value, &timestampStr, &m.Type)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("could not scan measurement: %w", err)
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}

		m.Timestamp = timestamp

		measurements = append(measurements, m)
	}

	if err = rows.Err(); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return measurements, nil
}

func (m *mySQL) AddMeasurement(ctx context.Context, measurement measurement.Measurement) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, InsertNewMeasurement, measurement.NodeID, measurement.Value, measurement.Timestamp, measurement.Type)
	if err != nil {
		return fmt.Errorf("could not insert measurement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
