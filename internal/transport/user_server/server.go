package user_server

import (
	"context"
	"errors"

	pb "github.com/suprt/gRPC1/app/api/user"
	"github.com/suprt/gRPC1/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedUserAPIServer
	service UserService
}

type UserService interface {
	Create(ctx context.Context, name, email, password, passwordConfirm, role string) (int64, error)
	Get(ctx context.Context, id int64) (model.UserGetResponse, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

func New(service UserService) *Server {
	return &Server{service: service}
}

func (s *Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := s.service.Create(ctx, req.GetName(), req.GetEmail(), req.GetPassword(), req.GetPasswordConfirm(), req.GetRole().String())
	if err != nil {
		return nil, toStatus(err)
	}

	response := pb.CreateResponse{Id: id}
	return &response, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	tempResponse, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, toStatus(err)
	}
	response := pb.GetResponse{
		Id:        tempResponse.ID,
		Name:      tempResponse.Name,
		Email:     tempResponse.Email,
		Role:      pb.Role(pb.Role_value[tempResponse.Role]),
		CreatedAt: ts.New(tempResponse.Created),
		UpdatedAt: ts.New(tempResponse.Updated),
	}
	return &response, nil
}

func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	var name, email *string
	if req.GetName() != nil {
		v := req.GetName().GetValue()
		name = &v
	}
	if req.GetEmail() != nil {
		v := req.GetEmail().GetValue()
		email = &v
	}
	if err := s.service.Update(ctx, req.GetId(), name, email); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.service.Delete(ctx, req.GetId()); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, model.ErrEmptyName),
		errors.Is(err, model.ErrEmptyEmail),
		errors.Is(err, model.ErrPasswordMatch),
		errors.Is(err, model.ErrPasswordSpace),
		errors.Is(err, model.ErrUserNotFound):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
