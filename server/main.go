package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/provider-grpc-server/proto"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", Port)

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterListServiceServer(s, &ListServiceServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

const (
	Port = ":50050"
)

var Lists map[string]*List

type List struct {
	Name        string
	Description string
	Items       []string
}

type ListServiceServer struct {
	pb.UnimplementedListServiceServer
}

func (s *ListServiceServer) Create(ctx context.Context, in *pb.CreateReq) (*pb.CreateResp, error) {
	if Lists[in.Name] != nil {
		return &pb.CreateResp{Status: "FAILED"}, errors.New(in.GetName() + " name not available")
	}

	Lists[in.Name] = &List{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Items:       make([]string, 0, 1),
	}

	fmt.Println("Create:: List: ", Lists[in.Name])
	return &pb.CreateResp{Status: "CREATED"}, nil
}

func (s *ListServiceServer) AddItems(ctx context.Context, in *pb.AddItemsReq) (*pb.AddItemsResp, error) {
	if Lists[in.Name] != nil {
		return &pb.AddItemsResp{Status: "FAILED"}, errors.New(in.GetName() + " list does not exist")
	}

	Lists[in.Name].Items = append(Lists[in.Name].Items, in.GetNewItems()...)

	fmt.Println("AddItems:: List: ", Lists[in.Name])
	return &pb.AddItemsResp{Status: "UPDATED"}, nil
}

func (s *ListServiceServer) Delete(ctx context.Context, in *pb.DeleteReq) (*pb.DeleteResp, error) {
	delete(Lists, in.GetName())
	fmt.Println("Delete:: Lists: ", Lists)
	return &pb.DeleteResp{Status: "DELETED"}, nil
}
