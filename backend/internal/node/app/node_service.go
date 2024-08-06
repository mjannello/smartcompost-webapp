package app

import (
	"context"
	"fmt"
	serialnumbergeneratorapp "github.com/mjannello/smartcompost-webapp/backend/internal/serial_number_generator/app"
	"github.com/mjannello/smartcompost-webapp/backend/pkg/clock"
	"log"
	"time"

	nodemodel "github.com/mjannello/smartcompost-webapp/backend/internal/node"
)

type NodeService interface {
	GetNodes(ctx context.Context) ([]nodemodel.Node, error)
	GetNode(ctx context.Context, nodeID uint64) (nodemodel.Node, error)
	GetNodeIDBySerialNumber(ctx context.Context, serialNumber string) (uint64, error)
	CreateNode(ctx context.Context, description, model string) (nodemodel.Node, error)
	UpdateNode(ctx context.Context, node nodemodel.Node) (nodemodel.Node, error)
	UpdateNodeLastUpdated(ctx context.Context, nodeID uint64, lastUpdated time.Time) error
	DeleteNode(ctx context.Context, nodeID uint64) (uint64, error)
}

type nodeService struct {
	nodeRepository        nodemodel.Repository
	serialNumberGenerator serialnumbergeneratorapp.SerialNumberGeneratorService
	realClock             clock.Clock
}

func NewNodeService(repository nodemodel.Repository, sngs serialnumbergeneratorapp.SerialNumberGeneratorService, realClock clock.Clock) NodeService {
	return &nodeService{nodeRepository: repository,
		serialNumberGenerator: sngs,
		realClock:             realClock,
	}
}

func (ns *nodeService) GetNodes(ctx context.Context) ([]nodemodel.Node, error) {
	nodes, err := ns.nodeRepository.GetAllNodes(ctx)
	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		return nil, fmt.Errorf("error fetching nodes: %w", err)
	}
	log.Printf("Fetched %d nodes", len(nodes))
	return nodes, nil
}

func (ns *nodeService) GetNode(ctx context.Context, nodeID uint64) (nodemodel.Node, error) {
	node, err := ns.nodeRepository.GetNodeByID(ctx, nodeID)
	if err != nil {
		log.Printf("Error fetching node by ID %d: %v", nodeID, err)
		return nodemodel.Node{}, fmt.Errorf("error getting node by ID: %w", err)
	}
	log.Printf("Fetched node: %+v", node)
	return node, nil
}

func (ns *nodeService) GetNodeIDBySerialNumber(ctx context.Context, serialNumber string) (uint64, error) {
	nodeID, err := ns.nodeRepository.GetNodeIDBySerialNumber(ctx, serialNumber)
	if err != nil {
		log.Printf("Error fetching nodeID: %v", err)
		return 0, fmt.Errorf("error fetching nodeID: %w", err)
	}
	log.Printf("Fetched %d nodeID", nodeID)
	return nodeID, nil
}

func (ns *nodeService) CreateNode(ctx context.Context, description, model string) (nodemodel.Node, error) {
	dateCreated := ns.realClock.Time()
	serialNumber, err := ns.serialNumberGenerator.New()
	if err != nil {
		log.Printf("Error creating node's serial number: %v", err)
		return nodemodel.Node{}, fmt.Errorf("error creating node's serial number: %w", err)
	}
	createdNode, err := ns.nodeRepository.CreateNode(ctx, serialNumber, description, model, dateCreated)
	if err != nil {
		log.Printf("Error creating node: %v", err)
		return nodemodel.Node{}, fmt.Errorf("error creating node: %w", err)
	}
	log.Printf("Created node: %+v", createdNode)
	return createdNode, nil
}

func (ns *nodeService) UpdateNode(ctx context.Context, node nodemodel.Node) (nodemodel.Node, error) {
	lastUpdated := ns.realClock.Time()
	updatedNode, err := ns.nodeRepository.UpdateNode(ctx, node, lastUpdated)
	if err != nil {
		log.Printf("Error updating node with ID %d: %v", node.ID, err)
		return nodemodel.Node{}, fmt.Errorf("error updating node: %w", err)
	}
	log.Printf("Updated node: %+v", updatedNode)
	return updatedNode, nil
}

// UpdateNodeLastUpdated uses the lastUpdated timestamp coming from the AP. Not the real time on web server
func (ns *nodeService) UpdateNodeLastUpdated(ctx context.Context, nodeID uint64, lastUpdated time.Time) error {
	node, err := ns.nodeRepository.GetNodeByID(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	node.LastUpdated = lastUpdated

	_, err = ns.UpdateNode(ctx, node)
	if err != nil {
		log.Printf("Error updating node last_updated: %v", err)
		return fmt.Errorf("error updating node last_updated: %w", err)
	}
	log.Printf("Node last_updated successfully for nodeID %d to: %s", nodeID, lastUpdated.String())
	return nil
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
