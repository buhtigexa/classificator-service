package main

import (
	bayesService "github.com/buhtigexa/classificator-service/protos"
	bayes "github.com/buhtigexa/naive-bayes/algorithms/bayes"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	bayesService.UnimplementedBayesServiceServer
	nb *bayes.NaiveBayes
}

func (s *Server) Predict(stream bayesService.BayesService_PredictServer) error {
	for {
		doc, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		predictions := s.nb.Predict(bayes.NewDocument(doc.Term, ""))
		for _, p := range predictions {
			resultPrediction := &bayesService.Prediction{
				Class: p.Class,
				Score: p.Prob,
			}
			if err := stream.Send(resultPrediction); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) Train(stream bayesService.BayesService_TrainServer) error {
	var corpus []bayes.Document
	for {
		doc, err := stream.Recv()
		if err == io.EOF {
			log.Printf(" Error %s\n", err)
			trainResult := s.nb.Train(corpus)
			response := to(trainResult)
			return stream.SendAndClose(response)
		}
		corpus = append(corpus, bayes.NewDocument(doc.Term, doc.Class))
	}
}

func to(trainResult *bayes.TrainResult) *bayesService.TrainResponse {
	classes := make(map[string]*bayesService.Class)
	for k, v := range trainResult.Classes {
		classes[k] = &bayesService.Class{
			Id:         v.Id,
			TotalWords: v.TotalWords,
			TotalDocs:  int32(v.TotalDocs),
			PriorProb:  float32(v.PriorProb),
		}
	}
	result := &bayesService.TrainResponse{
		Docs:    int32(trainResult.Docs),
		Classes: classes,
	}
	return result
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bayesService.RegisterBayesServiceServer(s, NewServer())
	log.Println("Server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}
}

func NewServer() *Server {
	return &Server{
		nb: bayes.NewNaiveBayes(),
	}
}
