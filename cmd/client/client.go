package main

import (
	"context"
	"fmt"
	bayesService "github.com/buhtigexa/classificator-service/protos"
	"github.com/buhtigexa/naive-bayes/algorithms/bayes"
	"google.golang.org/grpc"
	"io"
	"log"
)

func createCorpus() []bayes.Document {
	doc1 := bayes.NewDocument([]string{"dear", "friend", "launch", "money"}, "normal")
	doc2 := bayes.NewDocument([]string{"dear", "friend", "launch"}, "normal")
	doc3 := bayes.NewDocument([]string{"dear", "friend", "launch"}, "normal")
	doc4 := bayes.NewDocument([]string{"dear", "friend"}, "normal")
	doc5 := bayes.NewDocument([]string{"dear", "friend"}, "normal")
	doc6 := bayes.NewDocument([]string{"dear"}, "normal")
	doc7 := bayes.NewDocument([]string{"dear"}, "normal")
	doc8 := bayes.NewDocument([]string{"dear"}, "normal")

	doc9 := bayes.NewDocument([]string{"dear", "dear", "friend", "money"}, "spam")
	doc10 := bayes.NewDocument([]string{"money"}, "spam")
	doc11 := bayes.NewDocument([]string{"money"}, "spam")
	doc12 := bayes.NewDocument([]string{"money"}, "spam")

	corpus := []bayes.Document{doc1, doc2, doc3, doc4, doc5, doc6, doc7, doc8, doc9, doc10, doc11, doc12}
	return corpus
}

func train(client bayesService.BayesServiceClient) error {
	corpus := createCorpus()
	stream, err := client.Train(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < len(corpus); i++ {
		doc := &bayesService.Document{
			Term:  corpus[i].Terms,
			Class: corpus[i].Class,
		}
		if err := stream.Send(doc); err != nil {
			return err
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", response)
	return nil
}

func predict(client bayesService.BayesServiceClient) {

	waitc := make(chan struct{})
	doc1 := &bayesService.Document{
		Term:  []string{"launch", "money", "money", "money"},
		Class: "",
	}
	doc2 := &bayesService.Document{
		Term:  []string{"dear", "friend"},
		Class: "",
	}
	docs := []*bayesService.Document{doc1, doc2}
	stream, err := client.Predict(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer close(waitc)
		for {
			prediction, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Failed to receive a prediction: %v", err)
				break
			}
			fmt.Printf("Received prediction: %v\n", prediction)
		}
	}()

	for _, doc := range docs {
		if err := stream.Send(doc); err != nil {
			log.Printf("Failed to send a document: %v", err)
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	<-waitc

}
func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	client := bayesService.NewBayesServiceClient(conn)
	err = train(client)
	if err != nil {
		log.Fatal(err)
	}

}
