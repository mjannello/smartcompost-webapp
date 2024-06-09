package node

import "context"

type Repository interface {
	GetAllNodes(ctx context.Context) ([]Node, error)
}
