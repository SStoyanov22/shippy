package main

import (
	"context"

	pb "github.com/SStoyanov22/shippy/shippy-service-vessel/proto/vessel"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error)
	Create(ctx context.Context, vessel *Vessel) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		MaxWeight: spec.MaxWeight,
		Capacity:  spec.Capacity,
	}
}

func UnmarshalSpecification(spec *Specification) *pb.Specification {
	return &pb.Specification{
		MaxWeight: spec.MaxWeight,
		Capacity:  spec.Capacity,
	}
}

type Vessel struct {
	ID        string
	Capacity  int32
	Name      string
	Available bool
	OwnerID   string
	MaxWeight int32
}

type Specification struct {
	Capacity  int32
	MaxWeight int32
}
