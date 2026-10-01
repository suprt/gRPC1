package user_handler

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

type Handler struct {
	pb.UnimplementedUserAPIServer
	service UserService
}

type UserService interface {
	Create(ctx context.Context, name, email, password, passwordConfirm, role string) (int64, error)
	Get(ctx context.Context, id int64) (model.UserGetResponse, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

func New(service UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := h.service.Create(ctx, req.GetName(), req.GetEmail(), req.GetPassword(), req.GetPasswordConfirm(), req.GetRole().String())
	if err != nil {
		return nil, toStatus(err)
	}

	response := pb.CreateResponse{Id: id}
	return &response, nil
}

func (h *Handler) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	tempResponse, err := h.service.Get(ctx, req.GetId())
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

func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	var name, email *string
	if req.GetName() != nil {
		v := req.GetName().GetValue()
		name = &v
	}
	if req.GetEmail() != nil {
		v := req.GetEmail().GetValue()
		email = &v
	}
	if err := h.service.Update(ctx, req.GetId(), name, email); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.service.Delete(ctx, req.GetId()); err != nil {
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
