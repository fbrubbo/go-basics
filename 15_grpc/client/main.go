package main

import (
	"context"
	"log"
	"time"

	pb "github.com/fbrubbo/go-basics/15_grpc/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:4040"
)

func sum(client pb.MathServiceClient, a int64, b int64) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Request{A: a, B: b}
	resp, err := client.Add(ctx, req)
	if err != nil {
		panic(err)
	}
	return resp.GetResult()
}

func multiply(client pb.MathServiceClient, a int64, b int64) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Request{A: a, B: b}
	resp, err := client.Multiply(ctx, req)
	if err != nil {
		panic(err)
	}
	return resp.GetResult()
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewMathServiceClient(conn)

	log.Print("SUM IS: ", sum(client, 2, 5))
	log.Print("MULTIPLY IS: ", multiply(client, 2, 5))
}
