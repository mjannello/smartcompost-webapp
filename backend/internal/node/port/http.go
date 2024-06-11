package port

import (
	nodeapp "github.com/mjannello/smartcompost-webapp/backend/internal/node/app"
	"net/http"
)

type Handler interface {
	GetNode(w http.ResponseWriter, r *http.Request) error
}

type handler struct {
	nodeService nodeapp.NodeService
}

func NewRouterHandler(nodeService nodeapp.NodeService) Handler {
	return &handler{nodeService: nodeService}
}

func (h *handler) GetNode(w http.ResponseWriter, r *http.Request) error {
	//TODO implement me
	panic("implement me")
}
