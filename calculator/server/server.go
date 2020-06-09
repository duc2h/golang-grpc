package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/hoangduc02011998/golang-grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func (*Server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*Server) PriceNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PriceNumberDecompositionServer) error {
	log.Println("PriceNumberDecomposition is running")
	N := req.GetNumber()
	k := int32(2)
	for N > 1 {
		if N%k == 0 {
			N = N / k
			// send to client
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
		} else {
			k += 1
		}
	}
	defer log.Println("PriceNumberDecomposition is stopped")
	return nil
}

func (*Server) Average(stream calculatorpb.CalculatorService_AverageServer) error {

	log.Println("Average is called.")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(&resp)
		}

		if err != nil {
			log.Fatalf("err while Recv %v", err)
		}

		total += req.GetNumber()
		count++
	}
}

func (*Server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {

	log.Println("FindMax is called.")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while Recv : %v", err)
			return err
		}

		if max < req.GetNumber() {

			max = req.GetNumber()
			err = stream.Send(&calculatorpb.FindMaxResponse{
				Result: max,
			})

			if err != nil {
				log.Fatalf("Send FindMax err : %v", err)
				return err
			}
		}
	}
}

func (*Server) Square(ctx context.Context, req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Println("Square is called...")

	num := req.GetNumber()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "number can not >0 %v", num)
	}

	return &calculatorpb.SquareResponse{Result: math.Sqrt(float64(num))}, nil
}

func (*Server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	log.Println("SumWithDeadline is called ...")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("Context cancel")
			return nil, status.Error(codes.Canceled, "client cancel request")
		}
		time.Sleep(time.Second * 1)
	}

	resp := calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return &resp, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("err while create server : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})

	fmt.Println("Running")

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while serve : %v", err)
	}
}
