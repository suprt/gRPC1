package main

import (
	"context"
	"log"
	"time"

	"github.com/suprt/gRPC1/app/api/chat"
	pb "github.com/suprt/gRPC1/app/api/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func createUser(client pb.UserAPIClient, request *pb.CreateRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Create(ctx, request)
	if err != nil {
		log.Printf("client.CreateUser failed: %v", err)
		return
	}
	log.Printf("User created with id: %v", response)
}

func getUser(client pb.UserAPIClient, request *pb.GetRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Get(ctx, request)
	if err != nil {
		log.Printf("client.GetUser failed: %v", err)
		return
	}
	log.Printf("User retrieved by id:%d %v", request.GetId(), response)
}

func deleteUser(client pb.UserAPIClient, request *pb.DeleteRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Delete(ctx, request)
	if err != nil {
		log.Printf("client.DeleteUser failed: %v", err)
		return
	}
	log.Printf("User deleted by id:%d %v", request.GetId(), response)
}

func updateUser(client pb.UserAPIClient, request *pb.UpdateRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Update(ctx, request)
	if err != nil {
		log.Printf("client.UpdateUser failed: %v", err)
		return
	}
	log.Printf("User updated by id:%d %v", request.GetId(), response)
}
func longRequest(client pb.UserAPIClient, request *pb.UpdateRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	time.Sleep(2 * time.Second)
	response, err := client.Update(ctx, request)

	if err != nil {
		log.Printf("client.Update failed: %v", err)
		return
	}
	log.Printf("User updated by id:%d %v", request.GetId(), response)
}

func createChat(client chat.ChatAPIClient, request *chat.CreateRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Create(ctx, request)
	if err != nil {
		log.Printf("client.CreateChat failed: %v", err)
		return
	}
	log.Printf("Chat created with id: %v", response)
}

func deleteChat(client chat.ChatAPIClient, request *chat.DeleteRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.Delete(ctx, request)
	if err != nil {
		log.Printf("client.DeleteChat failed: %v", err)
		return
	}
	log.Printf("Chat deleted by id:%d %v", request.GetId(), response)
}

func sendMessage(client chat.ChatAPIClient, request *chat.SendMessageRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := client.SendMessage(ctx, request)
	if err != nil {
		log.Printf("client.SendMessage failed: %v", err)
		return
	}
	log.Printf("Message sent by id:%d %v", request.GetId(), response)
}
func main() {
	serverAddr := "localhost:50050"
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.NewClient failed: %v", err)
	}
	defer conn.Close()
	clientUser := pb.NewUserAPIClient(conn)
	clientChat := chat.NewChatAPIClient(conn)
	createUser(clientUser, &pb.CreateRequest{
		Name:            "user1",
		Email:           "user1@gmail.com",
		Password:        "qwerty",
		PasswordConfirm: "qwerty",
		Role:            pb.Role_USER,
	})
	createUser(clientUser, &pb.CreateRequest{
		Name:            "admin1",
		Email:           "admin1@gmail.com",
		Password:        "qwerty",
		PasswordConfirm: "qwerty",
		Role:            pb.Role_ADMIN,
	})
	getUser(clientUser, &pb.GetRequest{
		Id: 0,
	})
	getUser(clientUser, &pb.GetRequest{
		Id: 1,
	})
	getUser(clientUser, &pb.GetRequest{
		Id: 2,
	})
	deleteUser(clientUser, &pb.DeleteRequest{
		Id: 2,
	})
	time.Sleep(500 * time.Millisecond)
	updateUser(clientUser, &pb.UpdateRequest{
		Id:    1,
		Name:  wrapperspb.String("user 1 updated"),
		Email: nil,
	})
	updateUser(clientUser, &pb.UpdateRequest{
		Id:    0,
		Name:  wrapperspb.String("user 0 updated"),
		Email: nil,
	})
	getUser(clientUser, &pb.GetRequest{
		Id: 1,
	})

	createChat(clientChat, &chat.CreateRequest{
		Usernames: []string{"user1", "user2"},
	})
	createChat(clientChat, &chat.CreateRequest{
		Usernames: []string{"user1", "user2"},
	})
	deleteChat(clientChat, &chat.DeleteRequest{
		Id: 1,
	})
	deleteChat(clientChat, &chat.DeleteRequest{
		Id: 0,
	})
	sendMessage(clientChat, &chat.SendMessageRequest{
		From:      "user1",
		Text:      "Hello world!",
		Id:        2,
		Timestamp: timestamppb.Now(),
	})
	sendMessage(clientChat, &chat.SendMessageRequest{
		From:      "user1",
		Text:      "Hello world!",
		Timestamp: timestamppb.Now(),
	})

	longRequest(clientUser, &pb.UpdateRequest{})
}
