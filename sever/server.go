package main

import (
	"context"
	"fmt"

	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	hello "grpcServer/proto"
	"net"

	"google.golang.org/grpc"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "pass"
	dbname   = "postgres"
)

type Server struct {
	hello.UnimplementedGreeterServer
	hello.UnimplementedUserServiceServer
}

func (s *Server) AddUser(ctx context.Context, req *hello.UserRequest) (*hello.UserResponse, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	var id int32
	err = db.QueryRow("INSERT INTO users(name, age) VALUES($1, $2) RETURNING id", req.Name, req.Age).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user into database: %v", err)
	}

	return &hello.UserResponse{Message: fmt.Sprintf("User added with ID: %d", id)}, nil
}

func (s *Server) GetUser(ctx context.Context, req *hello.UserID) (*hello.User, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	var name string
	var age int32
	err = db.QueryRow("SELECT name, age FROM users WHERE id=$1", req.Id).Scan(&name, &age)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from database: %v", err)
	}

	return &hello.User{Id: req.Id, Name: name, Age: age}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterUserServiceServer(s, &Server{})
	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
