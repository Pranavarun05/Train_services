// client.go (gRPC client)
package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"grpc/grpc_train/train"  // Use a relative import path
)

// Rest of the code remains the same...


func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	// Example Usage
	user := &train.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

	// Purchase Ticket
	ticketRequest := &train.Ticket{From: "London", To: "France", User: user, PricePaid: 20.0}
	purchasedTicket, err := client.PurchaseTicket(context.Background(), ticketRequest)
	if err != nil {
		log.Fatalf("Failed to purchase ticket: %v", err)
	}
	fmt.Println("Purchased Ticket:", purchasedTicket)

	// Get Receipt
	receipt, err := client.GetReceipt(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to get receipt: %v", err)
	}
	fmt.Println("Receipt:", receipt)

	// Get Users by Section
	sectionRequest := &train.SectionRequest{Section: "A"}
	usersInSection, err := client.GetUsersBySection(context.Background(), sectionRequest)
	if err != nil {
		log.Fatalf("Failed to get users by section: %v", err)
	}
	for {
		userTicket, err := usersInSection.Recv()
		if err != nil {
			break
		}
		fmt.Println("User in Section A:", userTicket)
	}

	// Remove User
	removedTicket, err := client.RemoveUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to remove user: %v", err)
	}
	fmt.Println("Removed User Ticket:", removedTicket)

	// Modify User Seat
	modifiedTicket, err := client.ModifyUserSeat(context.Background(), purchasedTicket)
	if err != nil {
		log.Fatalf("Failed to modify user seat: %v", err)
	}
	fmt.Println("Modified User Seat:", modifiedTicket)
}
