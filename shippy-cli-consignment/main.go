// shippy/shippy-cli-consignment/main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/SStoyanov22/shippy/shippy-service-consignment/proto/consignment"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	service := micro.NewService(micro.Name("shippy-cli-consignment"))
	service.Init()

	client := pb.NewShippingService("shippy-service-consignment", service.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	var token string
	log.Println(os.Args)
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	if len(os.Args) < 3 {
		log.Fatal(errors.New("Not enough arguments, expecing file and token."))
	}
	file = os.Args[1]
	token = os.Args[2]

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	// Create a new context which contains our given token.
	// This same context will be passed into both the calls we make
	// to our consignment-service.
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})
	// First call using our tokenised context
	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
