package main

import (
	chittyChat "chittyChat/chatServer"
	"io"
	"log"
	"net"
	"google.golang.org/grpc"
)

type server struct {
	chittyChat.UnimplementedChatServiceServer
	Chats []chan *chittyChat.Message
}

func main() {
	log.SetFlags(0)
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("Error:", err)
	}else{
		log.Println("Server set up succesfully at port: 8080")
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chittyChat.RegisterChatServiceServer(grpcServer, &server{
		//Server keeps track of channels to communicate to via this array
		Chats: make([]chan *chittyChat.Message, 10),
	})

	nErr := grpcServer.Serve(lis);
	if nErr != nil {
		log.Fatal("Error:", err)
	}
}

func (s *server) SendMessage(msgStream chittyChat.ChatService_SendMessageServer) error {
	msg, err := msgStream.Recv()
	if err == io.EOF {
		return nil
	}else if err != nil {
		return err
	}

	if msg.Message != "was added to the chat" {
		log.Println("Message:", "\""+msg.Message+"\"", "received from", msg.User)
	}else{
		log.Println(msg.User, "joined", msg.Chat.Name)
	}

	go func() {
		streams := s.Chats
		for i := 0; i < len(streams); i++ {
			select {
			case streams[i] <- msg:
			default:
				streams[i] = streams[len(streams)-1]
				streams[len(streams)-1] = nil
				streams = streams[:len(streams)-1]
			}
		}
	}()

	return nil
}

func (s *server) JoinChat(Chat *chittyChat.Chat, msgStream chittyChat.ChatService_JoinChatServer) error {
	msgChannel := make(chan *chittyChat.Message)
	s.Chats = append(s.Chats, msgChannel)
	for {
		select {
			case msg := <-msgChannel:
			msgStream.Send(msg)
		}
	}
}

func (s *server) LeaveChat(Chat *chittyChat.Chat, msgStream chittyChat.ChatService_LeaveChatServer) error {
	log.Println(Chat.User, "left", Chat.Name)
	return nil
}
