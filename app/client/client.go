package main

import (
	"context"
	"io"
	"log"
	"message/proto/pb"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	address  = "localhost:3030"
	deadline = 20
)

func main() {
	message := os.Args[1]

	messageType := os.Args[2]

	var user int
	var err error

	if messageType == "relay" {
		userID := os.Args[3]

		user, err = strconv.Atoi(userID)
		if err != nil {
			log.Fatalf("failed to convert userid: %v", err)
		}
	}

	log.Printf("Arg %v, %v", message, messageType)

	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(deadline)*time.Second))

	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}

	// defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	// create stream
	stream, err := client.SendMessage(context.Background())
	if err != nil {
		log.Println(err)
	}
	// Send Message to stream
	req := pb.MessageRequest{Message: message, Type: messageType, UserIDs: []uint64{uint64(user)}}
	if err := stream.Send(&req); err != nil {
		log.Println(err)
	}

	// second goroutine receives data from stream
	// and prints result
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Fatalf("can not receive %v", err)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("new message %v received", resp)
		}
	}()

	time.Sleep(5 * time.Minute)
}
