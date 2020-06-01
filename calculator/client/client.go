package main

import (
	"context"
	"fmt"
	"github/hoangduc02011998/golang-grpc/calculator/calculatorpb"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err : %v", err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	//callSum(client)
	//log.Printf("service client %s", client)
	// callPND(client)
	// callAverage(client)
	//callFindMax(client)
	//callSquare(client)
	callSumWithDeadline(client)
}

func callSum(client calculatorpb.CalculatorServiceClient) {
	resp, err := client.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 4,
		Num2: 10,
	})
	if err != nil {
		log.Fatalf("err : %v", err)
	}

	fmt.Printf("result : %v", resp.GetResult())
}

func callSumWithDeadline(client calculatorpb.CalculatorServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		Num1: 4,
		Num2: 10,
	})
	if err != nil {

		if errStatus, ok := status.FromError(err); ok {

			log.Fatalf("err %v", errStatus.Message())
		}

		log.Fatalf("err : %v", err)
	}

	fmt.Printf("result : %v", resp.GetResult())
}

func callPND(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.PriceNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 1231,
	})

	if err == io.EOF {
		log.Println("server finish streaming")
	}

	if err != nil {
		log.Fatalf("err : %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("Price number %v", resp.GetResult())
	}
}

func callAverage(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.Average(context.Background())

	if err != nil {

	}

	for i := 0; i < 10; i++ {
		err = stream.Send(&calculatorpb.AverageRequest{Number: float32(i)})
		if err != nil {
			log.Fatalf("err while send request %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Response err :  %v", err)
	}

	log.Printf("Result : %v", resp.GetResult())
}

func callFindMax(client calculatorpb.CalculatorServiceClient) {

	stream, err := client.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error while send request: %v", err)
	}

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	//waitc := make(chan struct{})
	list := []int32{1, 31, 32, 11, 22, 1, -1, 32, 11}
	go func() {

		for _, l := range list {
			err = stream.Send(&calculatorpb.FindMaxRequest{Number: l})
			if err != nil {
				log.Fatalf("Error while send request %v", err)
				break
			}
		}
		waitGroup.Done()
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("server finish")
				//close(waitc)
				break
			}

			if err != nil {
				log.Fatalf("Err while receive response %v", err)
				break
			}

			log.Printf("Max is : %v\n", resp.GetResult())
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
	//<-waitc
}

func callSquare(client calculatorpb.CalculatorServiceClient) {
	number := int32(-8)
	resp, err := client.Square(context.Background(), &calculatorpb.SquareRequest{
		Number: number,
	})

	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err message : %v \n", errStatus.Message())
			log.Printf("err code: %v \n", errStatus.Code())

			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument : %v", number)
				return

			}
		}

	}

	log.Printf("Data response : %v", resp.GetResult())
}
