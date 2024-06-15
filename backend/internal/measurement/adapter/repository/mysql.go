package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	measurementmodel "github.com/mjannello/smartcompost-webapp/backend/internal/measurement"
)

const (
	DBDateTimeFormat                = "2006-01-02 15:04:05"
	GetAllMeasurementsByNodeIDQuery = "SELECT id, node_id, value, timestamp, type FROM measurements WHERE node_id = ?"
	GetMeasurementByIDQuery         = "SELECT id, node_id, value, timestamp, type FROM measurements WHERE id = ?"
	UpdateMeasurementQuery          = "UPDATE measurements SET value = ?, timestamp = ?, type = ? WHERE id = ?"
	DeleteMeasurementQuery          = "DELETE FROM measurements WHERE id = ?"
	AddMeasurementQuery             = "INSERT INTO measurements (node_id, value, timestamp, type) VALUES (?, ?, ?, ?)"
)

type mySQL struct {
	db *sql.DB
}

func NewMeasurementRepository(db *sql.DB) measurementmodel.Repository {
	return &mySQL{db: db}
}

func (m *mySQL) GetAllMeasurementsByNodeID(ctx context.Context, nodeID uint64) ([]measurementmodel.Measurement, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, GetAllMeasurementsByNodeIDQuery, nodeID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("could not get measurements from DB: %w", err)
	}
	defer rows.Close()

	var measurements []measurementmodel.Measurement
	for rows.Next() {
		var m measurementmodel.Measurement
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

		m.Timestamp = timestamp.UTC()
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

func (m *mySQL) GetMeasurementByID(ctx context.Context, measurementID uint64) (measurementmodel.Measurement, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return measurementmodel.Measurement{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	var measurement measurementmodel.Measurement
	var timestampStr string
	row := tx.QueryRowContext(ctx, GetMeasurementByIDQuery, measurementID)
	err = row.Scan(&measurement.ID, &measurement.NodeID, &measurement.Value, &timestampStr, &measurement.Type)
	if err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return measurement, fmt.Errorf("measurement not found")
		}
		return measurement, fmt.Errorf("could not scan measurement: %w", err)
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		_ = tx.Rollback()
		return measurement, fmt.Errorf("could not parse timestamp: %w", err)
	}
	measurement.Timestamp = timestamp.UTC()

	if err = tx.Commit(); err != nil {
		return measurement, fmt.Errorf("could not commit transaction: %w", err)
	}

	return measurement, nil
}

func (m *mySQL) UpdateMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return measurementmodel.Measurement{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, UpdateMeasurementQuery, measurement.Value, measurement.Timestamp.Format(time.RFC3339), measurement.Type, measurement.ID)
	if err != nil {
		_ = tx.Rollback()
		return measurementmodel.Measurement{}, fmt.Errorf("could not update measurement: %w", err)
	}

	updatedMeasurement, err := m.GetMeasurementByID(ctx, measurement.ID)
	if err != nil {
		_ = tx.Rollback()
		return measurementmodel.Measurement{}, fmt.Errorf("could not fetch updated measurement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return measurementmodel.Measurement{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return updatedMeasurement, nil
}

func (m *mySQL) DeleteMeasurement(ctx context.Context, measurementID uint64) (uint64, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, DeleteMeasurementQuery, measurementID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("could not delete measurement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}

	return measurementID, nil
}

func (m *mySQL) AddMeasurement(ctx context.Context, measurement measurementmodel.Measurement) (measurementmodel.Measurement, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return measurementmodel.Measurement{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	timestampStr := measurement.Timestamp.Format("2006-01-02 15:04:05")
	result, err := tx.ExecContext(ctx, AddMeasurementQuery, measurement.NodeID, measurement.Value, timestampStr, measurement.Type)
	if err != nil {
		_ = tx.Rollback()
		return measurementmodel.Measurement{}, fmt.Errorf("could not add measurement: %w", err)
	}

	measurementID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return measurementmodel.Measurement{}, fmt.Errorf("could not get last insert ID: %w", err)
	}

	measurement.ID = uint64(measurementID)

	if err = tx.Commit(); err != nil {
		return measurementmodel.Measurement{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return measurement, nil
}
