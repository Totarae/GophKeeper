package grpc

import (
	"GophKeeper/internal/common/logger"
	__ "GophKeeper/internal/pkg/proto_gen"
	"GophKeeper/internal/server/manager"
	"GophKeeper/internal/server/model"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Server) Register(ctx context.Context, req *__.RegisterRequest) (*__.TokenResponse, error) {
	token, err := s.userManager.Register(req.Login, req.Password, req.MasterPassword)
	if err != nil {
		return nil, convertError(err)
	}

	return &__.TokenResponse{Token: token}, nil
}

func (s *Server) Login(ctx context.Context, req *__.LoginRequest) (*__.TokenResponse, error) {
	token, err := s.userManager.Login(req.Login, req.Password, req.MasterPassword)
	if err != nil {
		return nil, convertError(err)
	}

	return &__.TokenResponse{Token: token}, nil
}

func (s *Server) Merge(ctx context.Context, req *__.MergeRequest) (*__.DataResponse, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	data := &model.UserData{
		UserID:    claims.SubjectID,
		DataKey:   req.DataKey,
		DataValue: req.DataValue,
		UpdatedAt: req.UpdatedAt.AsTime(),
		DeletedAt: req.DeletedAt.AsTime(),
	}
	err = s.dataManager.Merge(ctx, data)

	if err != nil {
		return nil, convertError(err)
	}

	return &__.DataResponse{
		DataKey:   data.DataKey,
		DataValue: data.DataValue,
		UpdatedAt: timestamppb.New(data.UpdatedAt),
		DeletedAt: timestamppb.New(data.DeletedAt),
	}, nil
}

func (s *Server) GetUpdates(ctx context.Context, req *__.GetUpdatesRequest) (*__.DataListResponse, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	updates, err := s.dataManager.GetUpdates(ctx, claims.SubjectID, req.UpdatedAfter.AsTime())
	if err != nil {
		return nil, convertError(err)
	}

	pbUpdates := make([]*__.DataResponse, 0, len(updates))
	for _, data := range updates {
		pbUpdates = append(pbUpdates, &__.DataResponse{
			DataKey:   data.DataKey,
			DataValue: data.DataValue,
			UpdatedAt: timestamppb.New(data.UpdatedAt),
			DeletedAt: timestamppb.New(data.DeletedAt),
		})
	}

	return &__.DataListResponse{Items: pbUpdates}, nil
}

func convertError(err error) error {
	switch {
	case errors.Is(err, manager.ErrUserExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, manager.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		logger.Logger.Error("error occurred", zap.Error(err))

		return status.Error(codes.Internal, "internal server error")
	}
}
