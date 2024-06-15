package node

import (
	"context"
)

type Repository interface {
	GetAllNodes(ctx context.Context) ([]Node, error)
	GetNodeByID(ctx context.Context, nodeID uint64) (Node, error)
	UpdateNode(ctx context.Context, node Node) (Node, error)
	DeleteNode(ctx context.Context, nodeID uint64) (uint64, error)
}
