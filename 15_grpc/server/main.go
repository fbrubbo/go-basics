package main

import (
	"context"
	"net"

	pb "github.com/fbrubbo/go-basics/15_grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a, b := req.GetA(), req.GetB()
	result := a + b
	return &pb.Response{Result: result}, nil
}

func (s *server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a, b := req.GetA(), req.GetB()
	result := a * b
	return &pb.Response{Result: result}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	pb.RegisterMathServiceServer(srv, &server{})
	reflection.Register(srv)
	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}
