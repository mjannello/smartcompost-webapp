package test

import (
	"github.com/google/uuid"
	"github.com/mjannello/smartcompost-webapp/backend/internal/node"
	"time"
)

func MakeNode(id uint64) node.Node {
	serial, _ := uuid.NewUUID()
	dateCreated := time.Now()
	return node.Node{
		ID:           id,
		SerialNumber: serial.String(),
		Model:        "Model",
		Description:  "Description",
		DateCreated:  dateCreated,
		LastUpdated:  dateCreated,
	}
}
