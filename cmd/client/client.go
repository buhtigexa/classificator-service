package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	bayesService "quadtree/protos"
	"strconv"
)

func BidiStream(client bayesService.BayesServiceClient) {
	waitc := make(chan struct{})
	stream, err := client.Classify(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer close(waitc)
	labelLoop:
		for {
			data, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break labelLoop
				}
				log.Fatal(err)
			}
			fmt.Printf("[CLIENT] Received from server %s\n", data.Resultado)
		}
	}()
	for i := 0; i < 3; i++ {
		stream.Send(&bayesService.DocumentRequest{Valor: strconv.Itoa(i) + ") documento enviado"})
	}
	stream.CloseSend()
	<-waitc
	fmt.Println(" Listo , ")
}

func trainModel(client bayesService.BayesServiceClient) {
	stream, err := client.TrainModel(context.Background())
	if err != nil {
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&bayesService.TermRequest{Value: strconv.Itoa(i)}); err != nil {

		}
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(" Summary ", summary)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := bayesService.NewBayesServiceClient(conn)
	response, err := client.SendTerm(context.Background(), &bayesService.TermRequest{Value: "Marce"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	for i := 0; i < 5; i++ {
		trainModel(client)
	}
	for i := 0; i < 5; i++ {
		BidiStream(client)
	}

}
