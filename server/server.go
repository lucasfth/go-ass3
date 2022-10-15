package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"
	proto "twitter/proto"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedTimeAskServer
	name string
	port int
}

var port = flag.Int("port", 8080, "Server port number")

func main() {
	flag.Parse()

	server := &Server{
		name: "TwitterServer",
		port: *port,
	}


	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen on port 8080: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 8080: %v", err)
	}
	
	go startServer(server)

	for {}
}

func startServer (server *Server) {
	grpcserver := grpc.NewServer()

	listener, err := net.Listen("tcp", ":" + strconv.Itoa(server.port))
	
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", server.port, err)
	}
	log.Printf("Started server at port: %d", server.port)

	proto.RegisterTimeAskServer(grpcserver, server)
	serveError := grpcserver.Serve(listener)
	if serveError != nil {
		log.Fatalf("Failed to serve gRPC server over port %d: %v", server.port, serveError)
	}
}

func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.ReplyTimeMessage, error) {
	log.Printf("Received request from %d", in.ClientId)
	return &proto.ReplyTimeMessage{
		Time: time.Now().String(),
		ServerName: c.name,
		}, nil
}