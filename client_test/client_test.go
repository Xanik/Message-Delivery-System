package clienttest

import (
	"context"
	"io"
	"log"
	"message/hub_test"
	"message/proto/pb"
	"testing"

	"google.golang.org/grpc"
)

// Initialize Server Service
func init() {
	go hub_test.NewHubServer(hub_test.NewHub())
}

// TestGetIdentity
func TestGetIdentity(t *testing.T) {
	port := ":3030"

	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}

	client := pb.NewMessageServiceClient(conn)
	// create stream
	stream, err := client.SendMessage(context.Background())
	if err != nil {
		t.Error(err)
	}
	ctx := stream.Context()
	// Send Message to stream
	req := pb.MessageRequest{Message: "Who Am I", Type: "identity"}
	if err := stream.Send(&req); err != nil {
		t.Error(err)
	}
	// check if message is received and close the stream
	resp, err := stream.Recv()
	if err == io.EOF {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	t.Logf("new message %v received", resp.UserIDs)
	ctx.Done()
}

// TestListConnectedUsers
func TestListConnectedUsersAndRelayMessage(t *testing.T) {
	port := ":3030"

	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}

	client := pb.NewMessageServiceClient(conn)

	// create stream
	stream, err := client.SendMessage(context.Background())
	if err != nil {
		t.Error(err)
	}

	ctx := stream.Context()

	// Send Message to stream
	req := pb.MessageRequest{Message: "Who Is Here?", Type: "list"}
	if err := stream.Send(&req); err != nil {
		t.Error(err)
	}

	// check if message is received and close the stream
	resp, err := stream.Recv()
	if err == io.EOF {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	t.Logf("new message %v received", resp.UserIDs)
	// Send Message to stream
	request := pb.MessageRequest{UserIDs: []uint64{1, 2}, Message: "foobar", Type: "relay"}
	if err := stream.Send(&request); err != nil {
		t.Error(err)
	}
	t.Log("relayed message to server")
	ctx.Done()
}
