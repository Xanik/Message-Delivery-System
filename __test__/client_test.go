package test

import (
	"context"
	"io"
	"log"
	"message/proto/pb"
	"net"
	"testing"

	"google.golang.org/grpc"
)

// Initialize Server Service
func init() {
	go MockServer(InitializeMockServer())
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// New Hub Client Connection Created
func newClient() pb.MessageServiceClient {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
	client := pb.NewMessageServiceClient(conn)
	return client
}

func TestGetIdentity(t *testing.T) {
	client := newClient()

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
	log.Printf("new message %s received", resp.Message)
	ctx.Done()
}

func TestListConnectedUsers(t *testing.T) {
	client := newClient()

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
	log.Printf("new message %v received", resp.UserIDs)
	ctx.Done()
}

func TestRelayMessage(t *testing.T) {
	client := newClient()

	// create stream
	stream, err := client.SendMessage(context.Background())
	if err != nil {
		t.Error(err)
	}

	ctx := stream.Context()

	// Send Message to stream
	req := pb.MessageRequest{UserIDs: []uint64{1, 2}, Message: "foobar", Type: "relay"}
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
	log.Printf("new message %s received", resp.Message)
	ctx.Done()
}
