package main

import (
	"context"
	hello "grpcServer/proto"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//c := hello.NewGreeterClient(conn)
	cl := hello.NewUserServiceClient(conn)
	// resp, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "soham"})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp.Message)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add user
	age := 30 // Example age
	addUserResp, err := cl.AddUser(ctx, &hello.UserRequest{Name: name, Age: int32(age)})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("Add user response: %s", addUserResp.Message)

	// Get user
	id := 1 // Example user ID
	getUserResp, err := cl.GetUser(ctx, &hello.UserID{Id: int32(id)})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("Get user response: %s (Age: %d)", getUserResp.Name, getUserResp.Age)

}
