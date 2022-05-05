package main

import (
	"context"
	"fmt"
	"log"
	"os"

	proto "github.com/SStoyanov22/shippy/shippy-service-user/proto/user"
	micro "go-micro.dev/v4"
)

func createUser(ctx context.Context, service micro.Service, user *proto.User) error {
	client := proto.NewUserService("shippy.service.user", service.Client())
	rsp, err := client.Create(ctx, user)
	if err != nil {
		return err
	}

	// print the response
	fmt.Println("Response: ", rsp.User)

	return nil
}

func main() {
	service := micro.NewService(micro.Name("shippy-cli-user"))
	service.Init()

	client := proto.NewUserService("shippy-service-user", service.Client())
	name := "Stoyan Stoyanov"
	password := "test1"
	company := "LTD"
	email := "ss@tt.com"
	user := &proto.User{
		Name:     name,
		Password: password,
		Company:  company,
		Email:    email,
	}
	r, err := client.Create(context.Background(), user)

	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}

	log.Printf("Created: %s", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &proto.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResponse, err := client.Auth(context.TODO(), &proto.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
	}

	log.Printf("Your access token is: %s \n", authResponse.Token)

	// let's just exit because
	os.Exit(0)
}
