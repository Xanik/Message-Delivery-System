package test

import (
	"context"
	"errors"
	"io"
	"log"
	"message/proto/pb"
)

type mockMessage struct {
	Message string
}

func (s *TestServer) GetIdentity(ctx context.Context, m *pb.GetIdentityRequest) (*pb.GetIdentityResponse, error) {
	return &pb.GetIdentityResponse{
		UserID: "userID",
	}, nil
}

func (s *TestServer) ListConnectedUsers(ctx context.Context, m *pb.ListConnectedUsersRequest) (*pb.ListConnectedUsersResponse, error) {
	return &pb.ListConnectedUsersResponse{
		UserID: []string{
			"userID",
			"userID2",
		},
	}, nil
}

func (s *TestServer) RelayMessage(m pb.MessageService_RelayMessageServer) error {
	ctx := m.Context()
	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := m.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		// continue if number reveived from stream
		// less than max
		if len(req.UserID) == 0 {
			return errors.New("No user specified")
		}
		// send it to stream
		resp := pb.RelayMessageResponse{Message: req.Message}
		if err := m.Send(&resp); err != nil {
			return err
		}
		// return nil to only run once
		return nil
	}
}
