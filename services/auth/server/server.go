package server

import (
	"github.com/bhrg3se/flahmingo-homework/services/auth/pb/proto"
	"github.com/bhrg3se/flahmingo-homework/services/auth/store"
)

func NewServer(store store.GenericStore) *Server {
	return &Server{store: store}
}

type Server struct {
	pb.UnimplementedAuthServiceServer
	store store.GenericStore
}
