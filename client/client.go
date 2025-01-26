package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sray/ticket"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcAddress = "localhost:50051"

func PurchaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	var purchaseReq ticket.PurchaseRequest

	if err := json.NewDecoder(r.Body).Decode(&purchaseReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)
	resp, err := client.PurchaseTicket(context.Background(), &purchaseReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to purchase ticket: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetTicketReceiptHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)
	resp, err := client.GetTicketReceipt(context.Background(), &ticket.ReceiptRequest{Email: email})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get ticket receipt: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ViewUsersInSectionHandler(w http.ResponseWriter, r *http.Request) {
	section := r.URL.Query().Get("section")

	if section != "A" && section != "B" {
		http.Error(w, "Section must be either A or B", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)
	resp, err := client.ViewUsersInSection(context.Background(), &ticket.SectionRequest{Section: section})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to view users in section: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	var removeReq ticket.RemoveUserRequest

	if err := json.NewDecoder(r.Body).Decode(&removeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)
	resp, err := client.RemoveUser(context.Background(), &removeReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ModifySeatHandler(w http.ResponseWriter, r *http.Request) {
	var modifyReq ticket.ModifySeatRequest

	if err := json.NewDecoder(r.Body).Decode(&modifyReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to gRPC server: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := ticket.NewTicketServiceClient(conn)
	resp, err := client.ModifySeat(context.Background(), &modifyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to modify seat: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/purchase-ticket", PurchaseTicketHandler)
	http.HandleFunc("/get-ticket-receipt", GetTicketReceiptHandler)
	http.HandleFunc("/view-users-in-section", ViewUsersInSectionHandler)
	http.HandleFunc("/remove-user", RemoveUserHandler)
	http.HandleFunc("/modify-seat", ModifySeatHandler)

	fmt.Println("HTTP server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
