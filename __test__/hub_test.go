package test

import (
	"errors"
	"io"
	"log"
	"message/proto/pb"
)

type mockMessage struct {
	Type    string
	Message string
}

func (s *TestServer) SendMessage(m pb.MessageService_SendMessageServer) error {
	mockUsers := []uint64{1, 2}
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

		switch req.Type {
		case "relay":
			// return if number reveived from stream
			// less than 0
			if len(req.UserIDs) == 0 {
				return errors.New("No user specified")
			}
			// send it to stream
			resp := pb.MessageResponse{Message: req.Message, UserIDs: mockUsers}
			if err := m.Send(&resp); err != nil {
				return err
			}
			// return nil to only run once
			return nil
		case "identity":
			// send it to stream
			resp := pb.MessageResponse{Message: "1", UserIDs: []uint64{1}}
			if err := m.Send(&resp); err != nil {
				return err
			}
			return nil
		case "list":
			// send it to stream
			resp := pb.MessageResponse{UserIDs: mockUsers}
			if err := m.Send(&resp); err != nil {
				return err
			}
		}
	}
}
