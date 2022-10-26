package main

import (
	"log"
	"net"
	"os"
	"chatster/chatserver"
	"google.golang.org/grpc"
)

func main() {

	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080" //default Port set to 8080 if PORT is not set in env
	}

	//init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcserver, &cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}

}
