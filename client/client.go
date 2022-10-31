package main

import (
	"bufio"
	chittyChat "chittyChat/chatServer"
	"context"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var lamportTimeStamp = int64(0)

func main() {
	//Removes time from Log prints
	log.SetFlags(0)

	username := enterUsername()
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatal("Error at 27:", err)
	}

	client := chittyChat.NewChatServiceClient(conn)

	go joinChat(ctx, client, username)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text(), username)
	}
}

func enterUsername() string {
	log.Print("Please enter a username: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	return input.Text()
}

func joinChat(ctx context.Context, client chittyChat.ChatServiceClient, username string) {
	Chat := chittyChat.Chat{Name: "chittyChat", User: username}
	stream, err := client.JoinChat(ctx, &Chat)
	if err != nil {
		log.Fatal("Error:", err)
	} else {
		log.Println("Server Connected, chat started")
		log.Println("-- To leave chat write 'leave chat'")
		log.Println("------------------------------")
	}

	sendMessage(ctx, client, "joined Chitty-Chat", username)

	waitc := make(chan struct{})
	go func() {
		for {
			messageIncoming, err := stream.Recv()
			if err != nil {
				log.Fatal("Error:", err)
			}
			if username != messageIncoming.User {
				lamportTimeStamp = int64(math.Max(float64(messageIncoming.LamportTimeStamp), float64(lamportTimeStamp)))
				LPTS := strconv.FormatInt(lamportTimeStamp, 10)
				if messageIncoming.Message == "joined Chitty-Chat" {
					log.Println("- Participant", messageIncoming.User, messageIncoming.Message, "at Lamport time", LPTS, "-")
				} else if messageIncoming.Message == "left chat" {
					log.Println("-", messageIncoming.User, messageIncoming.Message, "at Lamport time", LPTS, "-")
				} else {
					log.Println("["+LPTS+"]["+messageIncoming.User+"]", messageIncoming.Message)
				}
				lamportTimeStamp++
			}
		}
	}()
	<-waitc
}

func leaveChat(ctx context.Context, client chittyChat.ChatServiceClient, username string) {
	Chat := chittyChat.Chat{Name: "chittyChat", User: username}
	client.LeaveChat(ctx, &Chat)
	time.Sleep(10 * time.Millisecond)
	os.Exit(0)
}

func sendMessage(ctx context.Context, client chittyChat.ChatServiceClient, message string, username string) {
	if 128 < len(message) {
		log.Println("Message is too long. Max 128 characters")
		return
	}

	lamportTimeStamp++
	if message == "leave chat" {
		message = "left chat"
	}
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Fatal("Error:", err)
	}
	chatMessage := chittyChat.Message{
		Chat: &chittyChat.Chat{
			Name: "chittyChat",
			User: username},
		Message:          message,
		User:             username,
		LamportTimeStamp: lamportTimeStamp,
	}
	stream.Send(&chatMessage)
	if message == "left chat" {
		time.Sleep(10 * time.Millisecond)
		leaveChat(ctx, client, username)
	}
}
