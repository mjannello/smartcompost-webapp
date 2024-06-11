package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
	"time"
)

const (
	GetAllNodes = "SELECT id, description, type, last_updated FROM nodes"
)

type mySQL struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) node.Repository {
	return &mySQL{db: db}
}

func (m *mySQL) GetAllNodes(ctx context.Context) ([]node.Node, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, GetAllNodes)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("could not get nodes from DB: %w", err)
	}
	defer rows.Close()

	var nodes []node.Node
	for rows.Next() {
		var n node.Node
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
