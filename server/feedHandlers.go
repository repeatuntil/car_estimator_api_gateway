package server

import (
	"log/slog"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

type FeedHandler struct {
	r *mux.Router
	logger *slog.Logger
}


func (h *FeedHandler) setupgRPC(conn *grpc.ClientConn) {

}

func (h *FeedHandler) setupRoutes() {

}