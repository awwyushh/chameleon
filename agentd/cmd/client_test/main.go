package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/awwyushh/chameleon/agentd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to UDS
	conn, err := grpc.Dial("unix:///tmp/chameleon.sock", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChameleonClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Benign Request
	fmt.Println("--- Sending Benign Request ---")
	r1, err := c.Classify(ctx, &pb.ClassifyRequest{
		SrcIp: "127.0.0.1",
		Path:  "/home",
		Body:  "hello world",
	})
	if err != nil {
		log.Fatal(err)
	}
	printDec(r1)

	// 2. SQLi Request
	fmt.Println("\n--- Sending SQLi Request ---")
	r2, err := c.Classify(ctx, &pb.ClassifyRequest{
		SrcIp: "192.168.1.50",
		Path:  "/search",
		Body:  "SELECT * FROM users WHERE id=1 OR 1=1 --",
	})
	if err != nil {
		log.Fatal(err)
	}
	printDec(r2)
}

func printDec(r *pb.DecisionResponse) {
	fmt.Printf("Action: %v\n", r.Action)
	fmt.Printf("Label: %s (%.2f)\n", r.Label, r.Confidence)
	fmt.Printf("Delay: %d ms\n", r.DelayMs)
	fmt.Printf("Message: %s\n", r.Message)
	if r.HoneypotPort > 0 {
		fmt.Printf("Honeypot: %s:%d\n", r.HoneypotHost, r.HoneypotPort)
	}
}