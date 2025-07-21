package grpc

import (
	__ "GophKeeper/internal/pkg/proto_gen"
	"GophKeeper/internal/server/manager"
)

type Server struct {
	__.UnimplementedAuthServiceServer
	__.UnimplementedDataServiceServer

	userManager UserManagerInterface
	dataManager UserDataManagerInterface
}

func NewServer(
	userManager *manager.UserManager,
	dataManager *manager.UserDataManager,
) *Server {
	return &Server{
		userManager: userManager,
		dataManager: dataManager,
	}
}
