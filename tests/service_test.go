package tests

import (
	"encoding/json"
	"fmt"
	bayesService "github.com/buhtigexa/classificator-service/protos"
	"github.com/buhtigexa/naive-bayes/algorithms/bayes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func createCorpus() []bayes.Document {
	//doc1 := bayes.NewDocument([]string{"dear", "friend", "launch", "money"}, "normal")
	//doc2 := bayes.NewDocument([]string{"dear", "friend", "launch"}, "normal")
	//doc3 := bayes.NewDocument([]string{"dear", "friend", "launch"}, "normal")
	//doc4 := bayes.NewDocument([]string{"dear", "friend"}, "normal")
	//doc5 := bayes.NewDocument([]string{"dear", "friend"}, "normal")
	//doc6 := bayes.NewDocument([]string{"dear"}, "normal")
	//doc7 := bayes.NewDocument([]string{"dear"}, "normal")
	//doc8 := bayes.NewDocument([]string{"dear"}, "normal")
	//
	//doc9 := bayes.NewDocument([]string{"dear", "dear", "friend", "money"}, "spam")
	//doc10 := bayes.NewDocument([]string{"money"}, "spam")
	//doc11 := bayes.NewDocument([]string{"money"}, "spam")
	//doc12 := bayes.NewDocument([]string{"money"}, "spam")
	//
	//corpus := []bayes.Document{doc1, doc2, doc3, doc4, doc5, doc6, doc7, doc8, doc9, doc10, doc11, doc12}
	//
	f, err := os.Open("corpus.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var corpus []bayes.Document
	if err := dec.Decode(&corpus); err != nil {
		log.Fatal(err)
	}

	return corpus
}

func setUp() (func(), bayesService.BayesServiceClient) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	closeFn := func() {
		if err := conn.Close(); err != nil {
			log.Printf("%v\n", err)
		}
	}
	client := bayesService.NewBayesServiceClient(conn)
	return closeFn, client
}

func train(client bayesService.BayesServiceClient) {
	corpus := createCorpus()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()
	stream, err := client.Train(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < len(corpus); i++ {
		doc := &bayesService.Document{
			Term:  corpus[i].Terms,
			Class: corpus[i].Class,
		}
		if err := stream.Send(doc); err != nil {
			log.Println(err)
			return
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%v\n", response)
}

func TestTrain(t *testing.T) {
	var tearDown, client = setUp()
	defer tearDown()
	train(client)

}

func TestPredict(t *testing.T) {
	tearDown, client := setUp()
	train(client)
	defer tearDown()
	waitc := make(chan struct{})

	docs := []*bayesService.Document{{
		Term:  []string{"launch", "money", "money", "money"},
		Class: "",
	},
		{
			Term:  []string{"dear", "friend"},
			Class: "",
		},
	}

	stream, err := client.Predict(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer func() {
			close(waitc)
			fmt.Printf(" Closing channel \n")
		}()
		for {
			prediction, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("%s\n", err)
			}
			fmt.Printf("%v\n", prediction)
		}
	}()

	for _, doc := range docs {
		stream.Send(doc)
	}

	assert.Nil(t, stream.CloseSend())
	<-waitc
	fmt.Printf(" Stream finished ")

}
