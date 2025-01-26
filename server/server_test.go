package main

import (
	"context"
	"testing"

	"sray/ticket"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func createTestServer() *server {
	return &server{
		users:                    make([]map[*ticket.User]*ticket.TicketReceipt, 0),
		availableTicketsSectionA: 5,
		availableTicketsSectionB: 10,
	}
}

func TestServer(t *testing.T) {
	// Create the server instance
	s := createTestServer()

	// Test PurchaseTicket
	t.Run("PurchaseTicket", func(t *testing.T) {
		t.Run("Successful ticket purchase", func(t *testing.T) {
			req := &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "bob@example.com"},
				PricePaid: 1500,
				From:      "CityA",
				To:        "CityB",
			}
			resp, err := s.PurchaseTicket(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "CityA", resp.From)
			assert.Equal(t, "CityB", resp.To)
			assert.Equal(t, "bob@example.com", resp.User.Email)
			assert.Contains(t, resp.Seat, "A_")
		})

		t.Run("Duplicate user purchase attempt", func(t *testing.T) {
			req := &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "bob@example.com"},
				PricePaid: 1500,
				From:      "CityA",
				To:        "CityB",
			}
			resp, err := s.PurchaseTicket(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "already has a ticket")
		})

		t.Run("No seats available", func(t *testing.T) {
			for i := 0; i < 5; i++ {
				req := &ticket.PurchaseRequest{
					User:      &ticket.User{Email: "userA@example.com"},
					PricePaid: 1500,
					From:      "CityA",
					To:        "CityB",
				}
				req.User.Email = "userA" + string(rune(i+'0')) + "@example.com"
				_, _ = s.PurchaseTicket(context.Background(), req)
			}

			req := &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "overflow@example.com"},
				PricePaid: 1500,
				From:      "CityA",
				To:        "CityB",
			}
			resp, err := s.PurchaseTicket(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no seats available")
		})
	})

	// Test GetTicketReceipt
	t.Run("GetTicketReceipt", func(t *testing.T) {
		t.Run("Valid user email", func(t *testing.T) {
			req := &ticket.ReceiptRequest{Email: "bob@example.com"}
			resp, err := s.GetTicketReceipt(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "bob@example.com", resp.User.Email)
		})

		t.Run("Invalid user email", func(t *testing.T) {
			req := &ticket.ReceiptRequest{Email: "invalid@example.com"}
			resp, err := s.GetTicketReceipt(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ticket not found")
		})
	})

	// Test ViewUsersInSection
	t.Run("ViewUsersInSection", func(t *testing.T) {
		t.Run("Users in Section A", func(t *testing.T) {
			req := &ticket.SectionRequest{Section: "A"}
			resp, err := s.ViewUsersInSection(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.True(t, len(resp.Users) > 0)
		})

		t.Run("Users in Section B", func(t *testing.T) {
			req := &ticket.SectionRequest{Section: "B"}
			resp, err := s.ViewUsersInSection(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	})

	// Test RemoveUser
	t.Run("RemoveUser", func(t *testing.T) {
		t.Run("Remove existing user", func(t *testing.T) {
			req := &ticket.RemoveUserRequest{Email: "bob@example.com"}
			resp, err := s.RemoveUser(context.Background(), req)

			assert.NoError(t, err)
			assert.True(t, resp.Success)
		})

		t.Run("Remove non-existing user", func(t *testing.T) {
			req := &ticket.RemoveUserRequest{Email: "nonexistent@example.com"}
			resp, err := s.RemoveUser(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user with email")
		})
	})

	// Test ModifySeat
	t.Run("ModifySeat", func(t *testing.T) {
		t.Run("Modify seat for existing user", func(t *testing.T) {
			_, _ = s.PurchaseTicket(context.Background(), &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "alice@example.com"},
				PricePaid: 500,
				From:      "CityX",
				To:        "CityY",
			})

			req := &ticket.ModifySeatRequest{Email: "alice@example.com", NewSeat: "B_99"}
			resp, err := s.ModifySeat(context.Background(), req)

			assert.NoError(t, err)
			assert.True(t, resp.Success)
		})

		t.Run("Modify seat for non-existing user", func(t *testing.T) {
			req := &ticket.ModifySeatRequest{Email: "nonexistent@example.com", NewSeat: "B_99"}
			resp, err := s.ModifySeat(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user with email")
		})
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		t.Run("All users retrieved", func(t *testing.T) {
			_, _ = s.PurchaseTicket(context.Background(), &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "alice@example.com"},
				PricePaid: 1500,
				From:      "CityX",
				To:        "CityY",
			})
			_, _ = s.PurchaseTicket(context.Background(), &ticket.PurchaseRequest{
				User:      &ticket.User{Email: "charlie@example.com"},
				PricePaid: 500,
				From:      "CityX",
				To:        "CityY",
			})

			resp, err := s.GetAllUsers(context.Background(), &emptypb.Empty{})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, 6, len(resp.Users))
		})

		t.Run("No users available", func(t *testing.T) {
			emptyServer := createTestServer()

			resp, err := emptyServer.GetAllUsers(context.Background(), &emptypb.Empty{})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, 0, len(resp.Users))
		})
	})
}
