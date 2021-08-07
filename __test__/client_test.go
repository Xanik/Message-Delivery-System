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

// New Storage Client Connection Created
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

	mockData := mockMessage{
		Message: "Who Am I",
	}

	req := &pb.GetIdentityRequest{
		Message: mockData.Message,
	}

	res, err := client.GetIdentity(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if res.UserID != "userID" {
		t.Errorf("%v is not equal to %v", res, req)
	}
}

func TestListConnectedUsers(t *testing.T) {
	client := newClient()

	mockData := mockMessage{
		Message: "Who is Here",
	}

	req := &pb.ListConnectedUsersRequest{
		Message: mockData.Message,
	}

	res, err := client.ListConnectedUsers(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if len(res.UserID) != 2 {
		t.Errorf("%v is not equal to %v", res, req)
	}
}

func TestRelayMessage(t *testing.T) {
	client := newClient()

	// create stream
	stream, err := client.RelayMessage(context.Background())
	if err != nil {
		t.Error(err)
	}

	ctx := stream.Context()

	// Send Message to stream
	req := pb.RelayMessageRequest{UserID: []string{"1", "2"}, Message: "foobar"}
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
