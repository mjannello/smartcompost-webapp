package port

import "net/http"

type Handler interface {
	GetNode(w http.ResponseWriter, r *http.Request) error
}

type handler struct {
}

func NewRouterHandler() Handler {
	return &handler{}
}

func (h *handler) GetNode(w http.ResponseWriter, r *http.Request) error {
	//TODO implement me
	panic("implement me")
}
