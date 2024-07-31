package node

import (
	"context"
	"time"
)

type Repository interface {
	GetAllNodes(ctx context.Context) ([]Node, error)
	GetNodeByID(ctx context.Context, nodeID uint64) (Node, error)
	CreateNode(ctx context.Context, serialNumber, description, model string, dateCreated time.Time) (Node, error)
	UpdateNode(ctx context.Context, n Node, lastUpdated time.Time) (Node, error)
	DeleteNode(ctx context.Context, nodeID uint64) (uint64, error)
	GetNodeIDBySerialNumber(ctx context.Context, serialNumber string) (uint64, error)
}
