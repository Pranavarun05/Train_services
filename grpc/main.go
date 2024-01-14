// main.go (gRPC server)
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"google.golang.org/grpc"
	//"grpc/train"
)

// Rest of the code remains the same...


type trainServer struct {
	tickets []*train.Ticket
}

func (s *trainServer) PurchaseTicket(ctx context.Context, request *train.Ticket) (*train.Ticket, error) {
	if len(s.tickets)%2 == 0 {
		request.SeatSection = "A"
	} else {
		request.SeatSection = "B"
	}

	s.tickets = append(s.tickets, request)
	return request, nil
}

func (s *trainServer) GetReceipt(ctx context.Context, user *train.User) (*train.Ticket, error) {
	for _, ticket := range s.tickets {
		if strings.EqualFold(ticket.User.Email, user.Email) {
			return ticket, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *trainServer) GetUsersBySection(req *train.SectionRequest, stream train.TrainService_GetUsersBySectionServer) error {
	for _, ticket := range s.tickets {
		if strings.EqualFold(ticket.SeatSection, req.Section) {
			if err := stream.Send(ticket); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *trainServer) RemoveUser(ctx context.Context, user *train.User) (*train.Ticket, error) {
	for i, ticket := range s.tickets {
		if strings.EqualFold(ticket.User.Email, user.Email) {
			removedTicket := ticket
			s.tickets = append(s.tickets[:i], s.tickets[i+1:]...)
			return removedTicket, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *trainServer) ModifyUserSeat(ctx context.Context, request *train.Ticket) (*train.Ticket, error) {
	for _, ticket := range s.tickets {
		if strings.EqualFold(ticket.User.Email, request.User.Email) {
			ticket.SeatSection = request.SeatSection
			return ticket, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	train.RegisterTrainServiceServer(s, &trainServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
