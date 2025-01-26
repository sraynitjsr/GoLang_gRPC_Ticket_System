package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sray/ticket"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	ticket.UnimplementedTicketServiceServer
	mu                       sync.Mutex
	users                    []map[*ticket.User]*ticket.TicketReceipt
	availableTicketsSectionA int
	availableTicketsSectionB int
}

func (s *server) PurchaseTicket(ctx context.Context, req *ticket.PurchaseRequest) (*ticket.TicketReceipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userTobeAdded := req.User

	for _, newUser := range s.users {
		for existingUser := range newUser {
			if existingUser.Email == userTobeAdded.Email {
				return nil, fmt.Errorf("user %s already has a ticket, let others also travel please", userTobeAdded)
			}
		}
	}

	var seat string
	if req.PricePaid > 1000 && s.availableTicketsSectionA > 0 {
		seat = fmt.Sprintf("A_%d", s.availableTicketsSectionA)
		s.availableTicketsSectionA--
	} else if req.PricePaid <= 1000 && req.PricePaid >= 100 && s.availableTicketsSectionB > 0 {
		seat = fmt.Sprintf("B_%d", s.availableTicketsSectionB)
		s.availableTicketsSectionB--
	} else {
		return nil, fmt.Errorf("no seats available for the given price")
	}

	receipt := &ticket.TicketReceipt{
		From:      req.From,
		To:        req.To,
		User:      req.User,
		PricePaid: req.PricePaid,
		Seat:      seat,
	}

	userTicketMap := map[*ticket.User]*ticket.TicketReceipt{
		req.User: receipt,
	}

	s.users = append(s.users, userTicketMap)

	return receipt, nil
}

func (s *server) GetTicketReceipt(ctx context.Context, req *ticket.ReceiptRequest) (*ticket.TicketReceipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, userTicketMap := range s.users {
		for user, receipt := range userTicketMap {
			if user.Email == req.Email {
				return receipt, nil
			}
		}
	}

	return nil, fmt.Errorf("ticket not found for user with email: %s", req.Email)
}

func (s *server) ViewUsersInSection(ctx context.Context, req *ticket.SectionRequest) (*ticket.UsersInSection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var usersInSection []*ticket.UserWithSeat
	for _, userTicketMap := range s.users {
		for user, receipt := range userTicketMap {
			if (req.Section == "A" && receipt.Seat[0] == 'A') || (req.Section == "B" && receipt.Seat[0] == 'B') {
				usersInSection = append(usersInSection, &ticket.UserWithSeat{
					User: user,
					Seat: receipt.Seat,
				})
			}
		}
	}

	return &ticket.UsersInSection{
		Users: usersInSection,
	}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *ticket.RemoveUserRequest) (*ticket.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, userTicketMap := range s.users {
		for user := range userTicketMap {
			if user.Email == req.Email {
				s.users = append(s.users[:i], s.users[i+1:]...)
				return &ticket.RemoveUserResponse{Success: true}, nil
			}
		}
	}

	return nil, fmt.Errorf("user with email %s not found", req.Email)
}

func (s *server) ModifySeat(ctx context.Context, req *ticket.ModifySeatRequest) (*ticket.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, userTicketMap := range s.users {
		for user, receipt := range userTicketMap {
			if user.Email == req.Email {
				receipt.Seat = req.NewSeat
				return &ticket.ModifySeatResponse{Success: true}, nil
			}
		}
	}

	return nil, fmt.Errorf("user with email %s not found", req.Email)
}

func (s *server) GetAllUsers(ctx context.Context, req *emptypb.Empty) (*ticket.GetAllUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []*ticket.UserWithSeat
	for _, userTicketMap := range s.users {
		for user, receipt := range userTicketMap {
			users = append(users, &ticket.UserWithSeat{
				User: user,
				Seat: receipt.Seat,
			})
		}
	}

	return &ticket.GetAllUsersResponse{
		Users: users,
	}, nil
}

func main() {
	availableSeatsSectionA := 5
	availableSeatsSectionB := 10

	s := &server{
		users:                    make([]map[*ticket.User]*ticket.TicketReceipt, 0),
		availableTicketsSectionA: availableSeatsSectionA,
		availableTicketsSectionB: availableSeatsSectionB,
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ticket.RegisterTicketServiceServer(grpcServer, s)

	fmt.Println("Server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
