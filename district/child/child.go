package child

import (
	"context"
	"log"
	"time"

	pb "github.com/leeif/mercury/district/index/index"
	"google.golang.org/grpc"
)

func RegisterMember() {
	address := ":8080"
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIndexNodeClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RegisterMember(ctx, &pb.RegisterMemberRequest{Mid: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Register Member: %s", r.Message)
}
