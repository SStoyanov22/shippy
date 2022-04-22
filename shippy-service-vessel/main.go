package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	// Import the generated protobuf code
	pb "github.com/SStoyanov22/shippy/shippy-service-vessel/proto/vessel"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	mu      sync.RWMutex
	vessels []*pb.Vessel
}

type service struct {
	repo Repository
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}

	return nil, errors.New("No vessel found with that specs")
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification) (*pb.Response, error) {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Vessel: vessel}, nil
}

func main() {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels: vessels}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to liste %v", err)
	}

	s := grpc.NewServer()

	// Register reflection service on gRPC server.
	pb.RegisterVesselServiceServer(s, &service{repo})

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
