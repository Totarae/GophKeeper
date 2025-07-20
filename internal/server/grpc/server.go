package grpc

import "GophKeeper/internal/server/manager"

type Server struct {
	proto.UnimplementedAuthServiceServer
	proto.UnimplementedDataServiceServer

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
