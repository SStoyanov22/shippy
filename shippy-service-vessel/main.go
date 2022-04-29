package main

import (
	"context"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/SStoyanov22/shippy/shippy-service-vessel/proto/vessel"
	"go-micro.dev/v4"
)

const (
	port = ":50052"
)

func main() {
	service := micro.NewService(
		micro.Name("shippy-service-vessel"),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	vesselCollection := client.Database("shippy").Collection("vessels")
	repository := &MongoRepository{vesselCollection}

	h := &handler{repository}

	// Register our implementation with
	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
