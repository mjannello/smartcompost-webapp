package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
)

const GetAllNodes = "SELECT id, name FROM nodes"

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
		err = rows.Scan(&n.ID, &n.Description)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("could not scan node: %w", err)
		}
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
