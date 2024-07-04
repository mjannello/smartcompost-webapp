package node

import (
	"context"
	"time"
)

type Repository interface {
	GetAllNodes(ctx context.Context) ([]Node, error)
	GetNodeByID(ctx context.Context, nodeID uint64) (Node, error)
	CreateNode(ctx context.Context, fabricCode, description, nodeType string, lastUpdated time.Time) (Node, error)
	UpdateNode(ctx context.Context, node Node) (Node, error)
	DeleteNode(ctx context.Context, nodeID uint64) (uint64, error)
	GetNodeIDByFabricCode(ctx context.Context, fabricCode string) (uint64, error)
}
