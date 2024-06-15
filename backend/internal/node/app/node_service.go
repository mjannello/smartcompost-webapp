package app

import (
	"context"
	"fmt"
	"log"

	nodemodel "github.com/mjannello/smartcompost-webapp/backend/internal/node"
)

type NodeService interface {
	GetNodes(ctx context.Context) ([]nodemodel.Node, error)
	GetNodeByID(ctx context.Context, nodeID uint64) (nodemodel.Node, error)
	UpdateNode(ctx context.Context, node nodemodel.Node) (nodemodel.Node, error)
	DeleteNode(ctx context.Context, nodeID uint64) (uint64, error)
}

type nodeService struct {
	nodeRepository nodemodel.Repository
}

func NewNodeService(repository nodemodel.Repository) NodeService {
	return &nodeService{nodeRepository: repository}
}

func (ns *nodeService) GetNodes(ctx context.Context) ([]nodemodel.Node, error) {
	nodes, err := ns.nodeRepository.GetAllNodes(ctx)
	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		return nil, fmt.Errorf("error getting nodes: %w", err)
	}
	log.Printf("Fetched %d nodes", len(nodes))
	return nodes, nil
}

func (ns *nodeService) GetNodeByID(ctx context.Context, nodeID uint64) (nodemodel.Node, error) {
	node, err := ns.nodeRepository.GetNodeByID(ctx, nodeID)
	if err != nil {
		log.Printf("Error fetching node by ID %d: %v", nodeID, err)
		return nodemodel.Node{}, fmt.Errorf("error getting node by ID: %w", err)
	}
	log.Printf("Fetched node: %+v", node)
	return node, nil
}

func (ns *nodeService) UpdateNode(ctx context.Context, node nodemodel.Node) (nodemodel.Node, error) {
	updatedNode, err := ns.nodeRepository.UpdateNode(ctx, node)
	if err != nil {
		log.Printf("Error updating node with ID %d: %v", node.ID, err)
		return nodemodel.Node{}, fmt.Errorf("error updating node: %w", err)
	}
	log.Printf("Updated node: %+v", updatedNode)
	return updatedNode, nil
}

func (ns *nodeService) DeleteNode(ctx context.Context, nodeID uint64) (uint64, error) {
	deletedID, err := ns.nodeRepository.DeleteNode(ctx, nodeID)
	if err != nil {
		log.Printf("Error deleting node with ID %d: %v", nodeID, err)
		return 0, fmt.Errorf("error deleting node: %w", err)
	}
	log.Printf("Deleted node with ID %d", deletedID)
	return deletedID, nil
}
