# Ticket Booking System

## Overview
This repository implements a ticket booking system using gRPC and REST APIs. Below are the commands to generate protocol buffer files and examples of cURL commands for testing the REST endpoints.

---

## Generate Protocol Buffers Code
Run the following command to generate the necessary Go files for gRPC and protocol buffer support:
```bash
protoc --go_out=. --go-grpc_out=. ticket.proto
```

---

## REST API Endpoints

### 1. Purchase Ticket
**Endpoint:** `POST /purchase-ticket`  
**Description:** Allows a user to purchase a ticket.  

**Sample Request:**
```bash
curl -X POST "http://localhost:8080/purchase-ticket" \
-H "Content-Type: application/json" \
-d '{
  "from": "Mumbai",
  "to": "Delhi",
  "user": {
    "first_name": "Ravi",
    "last_name": "Kumar",
    "email": "ravi.kumar@example.com"
  },
  "price_paid": 500
}'
```

---

### 2. Get Ticket Receipt
**Endpoint:** `GET /get-ticket-receipt`  
**Description:** Fetch the ticket receipt for a user.  

**Sample Request:**
```bash
curl -X GET "http://localhost:8080/get-ticket-receipt?email=ravi.kumar@example.com"
```

---

### 3. View Users in a Section
**Endpoint:** `GET /view-users-in-section`  
**Description:** View all users in a specific section (A or B).  

**Sample Requests:**
- **Section A:**
  ```bash
  curl -X GET "http://localhost:8080/view-users-in-section?section=A"
  ```
- **Section B:**
  ```bash
  curl -X GET "http://localhost:8080/view-users-in-section?section=B"
  ```

---

### 4. Remove User
**Endpoint:** `POST /remove-user`  
**Description:** Remove a user from the booking system.  

**Sample Request:**
```bash
curl -X POST http://localhost:8080/remove-user \
-H "Content-Type: application/json" \
-d '{
  "email": "ravi.kumar@example.com"
}'
```

---

### 5. Modify Seat
**Endpoint:** `POST /modify-seat`  
**Description:** Modify the seat assigned to a user.  

**Sample Request:**
```bash
curl -X POST http://localhost:8080/modify-seat \
-H "Content-Type: application/json" \
-d '{
  "email": "ravi.kumar@example.com",
  "new_seat": "A_1"
}'
```

---

## Notes
- Ensure the server is running on `localhost:8080` before making requests.
- Replace email and other parameters in the requests as necessary.
