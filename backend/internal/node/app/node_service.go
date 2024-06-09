package app

import (
	"context"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
)

type NodeService interface {
	GetNodes(ctx context.Context) ([]node.Node, error)
}

type nodeService struct {
	nodeRepository node.Repository
}

func NewNodeService(repository node.Repository) NodeService {
	return &nodeService{nodeRepository: repository}
}

func (ns *nodeService) GetNodes(ctx context.Context) ([]node.Node, error) {
	nodes, err := ns.nodeRepository.GetAllNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %w", err)
	}
	return nodes, nil
}
