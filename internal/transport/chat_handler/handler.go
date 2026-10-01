package chat_handler

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/suprt/gRPC1/app/api/chat"
	"github.com/suprt/gRPC1/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedChatAPIServer
	service ChatService
}

type ChatService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, from, text string, id int64, timestamp time.Time) error
}

func New(service ChatService) *Handler {
	return &Handler{service: service}
}
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := h.service.Create(ctx, req.GetUsernames())
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: id}, nil
}

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (h *Handler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error) {
	err := h.service.SendMessage(ctx, req.GetFrom(), req.GetText(), req.GetId(), req.GetTimestamp().AsTime())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, model.ErrChatNotFound),
		errors.Is(err, model.ErrEmptyUsernames):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
