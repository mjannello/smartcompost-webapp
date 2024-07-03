package test

import (
	"context"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
	"github.com/stretchr/testify/mock"
	"time"
)

type NodeServiceMock struct {
	mock.Mock
}

func (ns *NodeServiceMock) GetNodes(ctx context.Context) ([]node.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (ns *NodeServiceMock) GetNodeByID(_ context.Context, nodeID uint64) (node.Node, error) {
	args := ns.Called(nodeID)
	n, _ := args.Get(0).(node.Node)
	e, _ := args.Get(1).(error)
	return n, e
}

func (ns *NodeServiceMock) UpdateNode(ctx context.Context, node node.Node) (node.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (ns *NodeServiceMock) UpdateNodeLastUpdated(_ context.Context, nodeID uint64, lastUpdated time.Time) error {
	args := ns.Called(nodeID, lastUpdated)
	e, _ := args.Get(0).(error)
	return e
}

func (ns *NodeServiceMock) DeleteNode(ctx context.Context, nodeID uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
