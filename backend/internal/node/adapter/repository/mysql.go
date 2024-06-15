package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	nodemodel "github.com/mjannello/smartcompost-webapp/backend/internal/node"
)

const (
	GetAllNodesQuery = "SELECT id, description, type, last_updated FROM nodes"
	GetNodeByIDQuery = "SELECT id, description, type, last_updated FROM nodes WHERE id = ?"
	UpdateNodeQuery  = "UPDATE nodes SET description = ?, type = ?, last_updated = ? WHERE id = ?"
	DeleteNodeQuery  = "DELETE FROM nodes WHERE id = ?"
)

type mySQL struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) nodemodel.Repository {
	return &mySQL{db: db}
}

func (m *mySQL) GetAllNodes(ctx context.Context) ([]nodemodel.Node, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, GetAllNodesQuery)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("could not get nodes from DB: %w", err)
	}
	defer rows.Close()

	var nodes []nodemodel.Node
	for rows.Next() {
		var n nodemodel.Node
		var lastUpdatedStr string
		err = rows.Scan(&n.ID, &n.Description, &n.Type, &lastUpdatedStr)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("could not scan node: %w", err)
		}

		lastUpdated, err := time.Parse("2006-01-02 15:04:05", lastUpdatedStr)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("could not parse last_updated: %w", err)
		}

		n.LastUpdated = lastUpdated
		nodes = append(nodes, n)
	}

	if err = rows.Err(); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return nodes, nil
}

func (m *mySQL) GetNodeByID(ctx context.Context, nodeID uint64) (nodemodel.Node, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nodemodel.Node{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	var n nodemodel.Node
	var lastUpdatedStr string
	row := tx.QueryRowContext(ctx, GetNodeByIDQuery, nodeID)
	err = row.Scan(&n.ID, &n.Description, &n.Type, &lastUpdatedStr)
	if err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return n, fmt.Errorf("node not found")
		}
		return n, fmt.Errorf("could not scan node: %w", err)
	}

	lastUpdated, err := time.Parse("2006-01-02 15:04:05", lastUpdatedStr)
	if err != nil {
		_ = tx.Rollback()
		return n, fmt.Errorf("could not parse last_updated: %w", err)
	}
	n.LastUpdated = lastUpdated

	if err = tx.Commit(); err != nil {
		return n, fmt.Errorf("could not commit transaction: %w", err)
	}

	return n, nil
}

func (m *mySQL) UpdateNode(ctx context.Context, n nodemodel.Node) (nodemodel.Node, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nodemodel.Node{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, UpdateNodeQuery, n.Description, n.Type, n.LastUpdated.Format("2006-01-02 15:04:05"), n.ID)
	if err != nil {
		_ = tx.Rollback()
		return nodemodel.Node{}, fmt.Errorf("could not update node: %w", err)
	}

	updatedNode, err := m.GetNodeByID(ctx, n.ID)
	if err != nil {
		_ = tx.Rollback()
		return nodemodel.Node{}, fmt.Errorf("could not fetch updated node: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nodemodel.Node{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return updatedNode, nil
}

func (m *mySQL) DeleteNode(ctx context.Context, nodeID uint64) (uint64, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, DeleteNodeQuery, nodeID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("could not delete node: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}

	return nodeID, nil
}
