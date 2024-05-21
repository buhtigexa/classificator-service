package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	bayesService "github.com/buhtigexa/classificator-service/protos"
	"github.com/buhtigexa/naive-bayes/algorithms/bayes"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

type Client struct {
	bayesService.BayesServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		bayesService.NewBayesServiceClient(conn),
	}
}

func (c *Client) trainModel(trainPath, savePath string) error {
	f, err := os.Open(trainPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var corpus []bayes.Document
	if err := dec.Decode(&corpus); err != nil {
		log.Fatal(err)
	}
	stream, err := c.Train(context.Background())
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

func (c *Client) predict(predictPath, dataPath string) {

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

	stream, err := c.Predict(context.Background())
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
	trainPath := flag.String("train", "", "path to the training dataset")
	savePath := flag.String("save", "", "path to save the trained model")
	predictPath := flag.String("predict", "", "path to data to be for prediction")

	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := NewClient(conn)
	if *trainPath != "" && *savePath != "" {
		client.trainModel(*trainPath, *savePath)
	} else if *predictPath != "" && *savePath != "" {
		client.predict(*predictPath, *savePath)
	} else {
		log.Fatal("You must provide either training options (-train and -save) or prediction options (-predict and -data).")
	}

}
