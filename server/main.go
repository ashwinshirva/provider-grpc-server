package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/ashwinshirva/provider-grpc-server/proto"
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
	Items       []int32
}

type ListServiceServer struct {
	pb.UnimplementedListServiceServer
}

func (s *ListServiceServer) CreateList(ctx context.Context, in *pb.CreateListReq) (*pb.CreateListResp, error) {
	if Lists[in.Name] != nil {
		//list := &List{}
		return &pb.CreateListResp{Status: "FAILED"}, errors.New(in.GetName() + " name not available")
	}

	if len(Lists) == 0 {
		Lists = make(map[string]*List)
	}

	fmt.Println("CreateList:: Create request received for list: ", in.GetName())
	Lists[in.Name] = &List{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Items:       make([]int32, 0, 1),
	}

	fmt.Println("CreateList:: List: ", Lists[in.Name])
	return &pb.CreateListResp{Status: "CREATED"}, nil
}

func (s *ListServiceServer) UpdateListItems(ctx context.Context, in *pb.UpdateListItemsReq) (*pb.UpdateListItemsResp, error) {
	fmt.Println("UpdateListItems:: UpdateListItems called...: ", in.GetName())
	printLists()
	if Lists[in.Name] == nil {
		return &pb.UpdateListItemsResp{Status: "FAILED"}, errors.New(in.GetName() + " list does not exist")
	}

	//Lists[in.Name].Items = append(Lists[in.Name].Items, in.GetNewItems()...)
	Lists[in.Name].Items = in.GetNewItems()

	fmt.Printf("UpdateListItems:: List: %v\n", Lists[in.Name])

	return &pb.UpdateListItemsResp{Status: "UPDATED"}, nil
}

func (s *ListServiceServer) DeleteList(ctx context.Context, in *pb.DeleteListReq) (*pb.DeleteListResp, error) {
	delete(Lists, in.GetName())
	fmt.Println("Delete:: Lists: ", Lists)
	return &pb.DeleteListResp{Status: "DELETED"}, nil
}

func (s *ListServiceServer) GetList(ctx context.Context, in *pb.GetListReq) (*pb.GetListResp, error) {
	if Lists[in.Name] == nil {
		return &pb.GetListResp{Status: "FAILED"}, errors.New(in.GetName() + " list does not exist")
	}

	printLists()
	return &pb.GetListResp{Status: "SUCCESS", Items: Lists[in.Name].Items}, nil
}

func printLists() {
	fmt.Printf("Lists: ")
	for _, list := range Lists {
		fmt.Printf("%v ", list)
	}
	fmt.Println()
}
