package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/suprt/gRPC1/app/api/chat"
	pb "github.com/suprt/gRPC1/app/api/user"
	"github.com/suprt/gRPC1/internal/repository/chat_repository"
	"github.com/suprt/gRPC1/internal/repository/user_repository"
	"github.com/suprt/gRPC1/internal/service/chat_service"
	"github.com/suprt/gRPC1/internal/service/user_service"
	"github.com/suprt/gRPC1/internal/transport/chat_handler"
	"github.com/suprt/gRPC1/internal/transport/user_handler"
	"google.golang.org/grpc"
)

func main() {
	port := 50050
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	log.Println("Server starting...")
	if err != nil {
		log.Fatal(err)
	}
	userRepo := user_repository.New()
	userService := user_service.New(userRepo)
	userHandler := user_handler.New(userService)

	chatRepo := chat_repository.New()
	chatService := chat_service.New(chatRepo)
	chatHandler := chat_handler.New(chatService)

	grpcServer := grpc.NewServer()
	pb.RegisterUserAPIServer(grpcServer, userHandler)
	chat.RegisterChatAPIServer(grpcServer, chatHandler)
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}

	}()
	log.Println("Server started")
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
	log.Println("Shutting down...")
	grpcServer.GracefulStop()

}
