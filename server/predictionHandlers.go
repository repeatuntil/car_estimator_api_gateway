package server

import (
	"log/slog"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

type PredictionHandler struct {
	r *mux.Router
	logger *slog.Logger
}

func (h *PredictionHandler) setupgRPC(conn *grpc.ClientConn) {

}

func (h *PredictionHandler) setupRoutes() {

}