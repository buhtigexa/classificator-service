package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	bayesService "quadtree/protos"
	"strconv"
)

type Server struct {
	bayesService.UnimplementedBayesServiceServer
}

func (s *Server) TrainModel(stream bayesService.BayesService_TrainModelServer) error {
	class := bayesService.Class{
		Label:           "label",
		PriorLikelihood: 12.0,
		Terms:           1321,
	}
	classes := []*bayesService.Class{&class}
	for {
		data, err := stream.Recv()
		fmt.Printf("%v", data)
		if err == io.EOF {
			return stream.SendAndClose(&bayesService.SummaryResponse{Classes: classes})
		}
	}

	return nil
}

func (s *Server) Classify(stream bayesService.BayesService_ClassifyServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Printf("[SERVER] Received from client  %s\n", data.Valor)

		for i := 0; i < 3; i++ {
			fmt.Printf("[SERVER] enviando al client \n")
			stream.Send(&bayesService.DocumentClassification{Resultado: strconv.Itoa(i) + ") Resultado de clasificacion "})
		}
	}
	return nil
}

func (s *Server) SendTerm(ctx context.Context, request *bayesService.TermRequest) (*bayesService.TermResponse, error) {
	data := fmt.Sprintf("Hola %s !! ", request.Value)
	return &bayesService.TermResponse{Data: data}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bayesService.RegisterBayesServiceServer(s, &Server{})
	log.Println("Server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}
}
