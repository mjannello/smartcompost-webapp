package app_test

import (
	"context"
	"fmt"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"github.com/mjannello/smartcompost-webapp/backend/pkg/clock"
	"github.com/mjannello/smartcompost-webapp/backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNodeService_CreateNode(t *testing.T) {

	type depFields struct {
		nodeRepository               *test.NodeRepositoryMock
		serialNumberGeneratorService *test.SerialNumberGeneratorServiceMock
		realClock                    *clock.ClockMock
	}
	type input struct {
		description, model string
	}
	type output struct {
		node node.Node
		err  error
	}
	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)
	expectedNode := node.Node{ID: uint64(1)}
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error creating node's serial number",
			in:   input{description: "Some description", model: "Some model"},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.serialNumberGeneratorService.On("New").Return("", fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error creating node's serial number: test")
				assert.Equal(t, node.Node{}, out.node)
			},
		},
		{
			name: "error creating node",
			in:   input{description: "Some description", model: "Some model"},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.serialNumberGeneratorService.On("New").Return("serial", nil)
				df.nodeRepository.On("CreateNode", "serial", "Some description", "Some model", timeNow).Return(node.Node{}, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error creating node: test")
				assert.Equal(t, node.Node{}, out.node)
			},
		},
		{
			name: "create node successfully",
			in:   input{description: "Some description", model: "Some model"},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.serialNumberGeneratorService.On("New").Return("serial", nil)
				df.nodeRepository.On("CreateNode", "serial", "Some description", "Some model", timeNow).Return(expectedNode, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedNode, out.node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			serialNumberGeneratorServiceMock := &test.SerialNumberGeneratorServiceMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewNodeService(nodeRepositoryMock, serialNumberGeneratorServiceMock, clockMock)

			df := &depFields{nodeRepository: nodeRepositoryMock, serialNumberGeneratorService: serialNumberGeneratorServiceMock, realClock: clockMock}
			tt.on(df)

			// When
			nodeCreated, err := s.CreateNode(context.Background(), tt.in.description, tt.in.model)

			// Then
			tt.assert(t, &output{nodeCreated, err})
			nodeRepositoryMock.AssertExpectations(t)
			serialNumberGeneratorServiceMock.AssertExpectations(t)
			clockMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_GetNodes(t *testing.T) {

	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
	}

	type output struct {
		nodes []node.Node
		err   error
	}
	expectedNodes := []node.Node{
		{ID: uint64(1)},
		{ID: uint64(2)},
	}
	tests := []struct {
		name   string
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error fetching nodes",
			on: func(df *depFields) {
				df.nodeRepository.On("GetAllNodes").Return(nil, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error fetching nodes: test")
				assert.Nil(t, out.nodes)
			},
		},
		{
			name: "fetch nodes successfully",
			on: func(df *depFields) {
				df.nodeRepository.On("GetAllNodes").Return(expectedNodes, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedNodes, out.nodes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, nil)

			df := &depFields{nodeRepository: nodeRepositoryMock}
			tt.on(df)

			// When
			nodesFetched, err := s.GetNodes(context.Background())

			// Then
			tt.assert(t, &output{nodesFetched, err})
			nodeRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_GetNode(t *testing.T) {

	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
	}

	type input struct {
		nodeID uint64
	}

	type output struct {
		node node.Node
		err  error
	}
	expectedNode := test.MakeNode(1)
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error fetching node",
			in:   input{uint64(1)},
			on: func(df *depFields) {
				df.nodeRepository.On("GetNodeByID", uint64(1)).Return(node.Node{}, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error getting node by ID: test")
				assert.Equal(t, node.Node{}, out.node)
			},
		},
		{
			name: "fetch node successfully",
			in:   input{uint64(1)},
			on: func(df *depFields) {
				df.nodeRepository.On("GetNodeByID", expectedNode.ID).Return(expectedNode, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, expectedNode, out.node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, nil)

			df := &depFields{nodeRepository: nodeRepositoryMock}
			tt.on(df)

			// When
			nodeFetched, err := s.GetNode(context.Background(), tt.in.nodeID)

			// Then
			tt.assert(t, &output{nodeFetched, err})
			nodeRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_GetNodeIDBySerialNumber(t *testing.T) {

	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
	}

	type input struct {
		serialNumber string
	}

	type output struct {
		nodeID uint64
		err    error
	}
	expectedNode := test.MakeNode(1)
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error fetching node ID",
			in:   input{expectedNode.SerialNumber},
			on: func(df *depFields) {
				df.nodeRepository.On("GetNodeIDBySerialNumber", expectedNode.SerialNumber).Return(0, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error fetching nodeID: test")
				assert.Equal(t, uint64(0), out.nodeID)
			},
		},
		//{
		//	name: "fetch node successfully",
		//	in:   input{uint64(1)},
		//	on: func(df *depFields) {
		//		df.nodeRepository.On("GetNodeByID", expectedNode.ID).Return(expectedNode, nil)
		//	},
		//	assert: func(t *testing.T, out *output) {
		//		assert.NoError(t, out.err)
		//		assert.Equal(t, expectedNode, out.node)
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, nil)

			df := &depFields{nodeRepository: nodeRepositoryMock}
			tt.on(df)

			// When
			nodeID, err := s.GetNodeIDBySerialNumber(context.Background(), tt.in.serialNumber)

			// Then
			tt.assert(t, &output{nodeID, err})
			nodeRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_UpdateNode(t *testing.T) {
	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
		realClock      *clock.ClockMock
	}
	type input struct {
		node node.Node
	}
	type output struct {
		node node.Node
		err  error
	}
	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)
	nodeToUpdate := test.MakeNode(1)
	updatedNode := nodeToUpdate
	updatedNode.LastUpdated = timeNow
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error updating node",
			in:   input{node: nodeToUpdate},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.nodeRepository.On("UpdateNode", nodeToUpdate, timeNow).Return(node.Node{}, fmt.Errorf("test"))

			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error updating node: test")
				assert.Equal(t, node.Node{}, out.node)
			},
		},
		{
			name: "update node successfully",
			in:   input{node: nodeToUpdate},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.nodeRepository.On("UpdateNode", nodeToUpdate, timeNow).Return(updatedNode, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, updatedNode, out.node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, clockMock)

			df := &depFields{nodeRepository: nodeRepositoryMock, realClock: clockMock}
			tt.on(df)

			// When
			nodeCreated, err := s.UpdateNode(context.Background(), tt.in.node)

			// Then
			tt.assert(t, &output{nodeCreated, err})
			nodeRepositoryMock.AssertExpectations(t)
			clockMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_UpdateNodeLastUpdated(t *testing.T) {
	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
		realClock      *clock.ClockMock
	}
	type input struct {
		nodeID      uint64
		lastUpdated time.Time
	}
	type output struct {
		err error
	}
	newLastUpdated := time.Date(2024, 06, 20, 22, 30, 30, 0, time.UTC)
	timeNow := time.Date(2024, 06, 20, 20, 15, 30, 0, time.UTC)

	nodeToUpdate := test.MakeNode(1)
	updatedNode := nodeToUpdate
	updatedNode.LastUpdated = newLastUpdated
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error getting node",
			in:   input{nodeID: uint64(1), lastUpdated: newLastUpdated},
			on: func(df *depFields) {
				df.nodeRepository.On("GetNodeByID", uint64(1)).Return(node.Node{}, fmt.Errorf("test"))

			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "node not found: test")
			},
		},
		{
			name: "update node lastUpdated successfully",
			in:   input{nodeID: uint64(1), lastUpdated: newLastUpdated},
			on: func(df *depFields) {
				df.realClock.On("Time").Return(timeNow)
				df.nodeRepository.On("GetNodeByID", uint64(1)).Return(nodeToUpdate, nil)
				df.nodeRepository.On("UpdateNode", updatedNode, timeNow).Return(updatedNode, nil)

			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			clockMock := &clock.ClockMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, clockMock)

			df := &depFields{nodeRepository: nodeRepositoryMock, realClock: clockMock}
			tt.on(df)

			// When
			err := s.UpdateNodeLastUpdated(context.Background(), tt.in.nodeID, tt.in.lastUpdated)

			// Then
			tt.assert(t, &output{err})
			nodeRepositoryMock.AssertExpectations(t)
			clockMock.AssertExpectations(t)
		})
	}
}

func TestNodeService_DeleteNode(t *testing.T) {

	type depFields struct {
		nodeRepository *test.NodeRepositoryMock
	}

	type input struct {
		nodeID uint64
	}

	type output struct {
		deletedNodeID uint64
		err           error
	}
	tests := []struct {
		name   string
		in     input
		on     func(*depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "error deleting node",
			in:   input{uint64(1)},
			on: func(df *depFields) {
				df.nodeRepository.On("DeleteNode", uint64(1)).Return(node.Node{}, fmt.Errorf("test"))
			},
			assert: func(t *testing.T, out *output) {
				assert.ErrorContains(t, out.err, "error deleting node: test")
				assert.Equal(t, uint64(0), out.deletedNodeID)
			},
		},
		{
			name: "deleted node successfully",
			in:   input{uint64(1)},
			on: func(df *depFields) {
				df.nodeRepository.On("DeleteNode", uint64(1)).Return(uint64(1), nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, uint64(1), out.deletedNodeID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			nodeRepositoryMock := &test.NodeRepositoryMock{}
			s := app.NewNodeService(nodeRepositoryMock, nil, nil)

			df := &depFields{nodeRepository: nodeRepositoryMock}
			tt.on(df)

			// When
			deletedNodeID, err := s.DeleteNode(context.Background(), tt.in.nodeID)

			// Then
			tt.assert(t, &output{deletedNodeID, err})
			nodeRepositoryMock.AssertExpectations(t)
		})
	}
}
