package grpc

import (
	"context"
	"github.com/Totarae/GophKeeper/internal/client/model"
	pb "github.com/Totarae/GophKeeper/internal/pkg/proto_gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn       *grpc.ClientConn
	AuthClient pb.AuthServiceClient
	DataClient pb.DataServiceClient
	authToken  string
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:       conn,
		AuthClient: pb.NewAuthServiceClient(conn),
		DataClient: pb.NewDataServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, login, password string, masterPassword []byte) (string, error) {
	resp, err := c.AuthClient.Login(ctx, &pb.LoginRequest{
		Login:          login,
		Password:       password,
		MasterPassword: string(masterPassword),
	})
	if err != nil {
		return "", err
	}
	c.authToken = resp.Token
	return resp.Token, nil
}

func (c *Client) Register(ctx context.Context, login, password string, masterPassword []byte) (string, error) {
	resp, err := c.AuthClient.Register(ctx, &pb.RegisterRequest{
		Login:          login,
		Password:       password,
		MasterPassword: string(masterPassword),
	})
	if err != nil {
		return "", err
	}
	c.authToken = resp.Token
	return resp.Token, nil
}

func (c *Client) Merge(ctx context.Context, data *model.UserData) (*pb.DataResponse, error) {
	ctx = c.withAuth(ctx)
	return c.DataClient.Merge(ctx, &pb.MergeRequest{
		DataKey:   data.DataKey,
		DataValue: data.DataValue,
		UpdatedAt: timestamppb.New(data.UpdatedAt),
		DeletedAt: timestamppb.New(data.DeletedAt),
	})
}

func (c *Client) GetUpdates(ctx context.Context, updatedAfter time.Time) (*pb.DataListResponse, error) {
	ctx = c.withAuth(ctx)
	return c.DataClient.GetUpdates(ctx, &pb.GetUpdatesRequest{
		UpdatedAfter: timestamppb.New(updatedAfter),
	})
}

func (c *Client) withAuth(ctx context.Context) context.Context {
	if c.authToken != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.authToken)
	}
	return ctx
}
