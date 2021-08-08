package hub_test

import (
	"io"
	"log"
	"math/rand"
	"message/models"
	"message/proto/pb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// New Hub Server Connection Created
func NewHubServer(server pb.MessageServiceServer) {
	port := ":3030"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to Listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterMessageServiceServer(s, server)

	log.Println("Starting Server On Port:" + port)

	e := s.Serve(lis)
	if e != nil {
		log.Fatalf("failed to Serve: %v", e)
	}
}

type hub struct {
	mutex *sync.RWMutex
	store map[pb.MessageService_SendMessageServer]uint64
}

//Initialized a constructor  Of hub struct
func NewHub() *hub {
	return &hub{mutex: &sync.RWMutex{}, store: make(map[pb.MessageService_SendMessageServer]uint64)}
}

//Generate Ascending UserIDs
func (s hub) generateID() uint64 {
	UserIDs := uint64(rand.Intn(1000))
	return UserIDs
}

// SendMessage Listens for sendmessage stream from the server
func (s hub) SendMessage(m pb.MessageService_SendMessageServer) error {
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
			// check stream and store stream if new
			s.mutex.RLock()
			_, ok := s.store[m]
			s.mutex.RUnlock()

			if !ok {
				key := s.generateID()
				s.mutex.Lock()
				s.store[m] = key
				s.mutex.Unlock()
			}
			// return if number reveived from stream
			// less than 0
			if len(req.UserIDs) == 0 {
				log.Printf(models.Message(models.UsersNotFound))
				continue
			}
			for k, v := range s.store {
				for _, UserIDs := range req.UserIDs {
					if v == UserIDs {
						// send it to stream
						resp := pb.MessageResponse{Message: req.Message}
						if err := k.Send(&resp); err != nil {
							log.Println(err.Error())
						}
					}
				}
			}
		case "identity":
			// check stream and store stream if new
			s.mutex.RLock()
			data, ok := s.store[m]
			s.mutex.RUnlock()

			if !ok {
				key := s.generateID()
				s.mutex.Lock()
				s.store[m] = key
				s.mutex.Unlock()
				data = key
			}
			resp := pb.MessageResponse{UserIDs: []uint64{data}}
			if err := m.Send(&resp); err != nil {
				log.Println(err.Error())
			}
		case "list":
			// check stream and store stream if new
			s.mutex.RLock()
			_, ok := s.store[m]
			s.mutex.RUnlock()

			if !ok {
				key := s.generateID()
				s.mutex.Lock()
				s.store[m] = key
				s.mutex.Unlock()
			}

			s.mutex.RLock()
			if len(s.store) == 0 {
				log.Println(models.Message(models.UsersNotFound))
				s.mutex.RUnlock()
			}
			users := []uint64{}

			for k, v := range s.store {
				if k != m {
					users = append(users, v)
				}
			}
			s.mutex.RUnlock()
			// send it to stream
			resp := pb.MessageResponse{UserIDs: users}
			if err := m.Send(&resp); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
