package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	pb "github.com/GianniCarlini/Lab-3-SD/proto"
)

const (
	port = ":50050" //Broker
	dns1 = "localhost:50051"
	dns2 = "localhost:50052"
	dns3 = "localhost:50053"
)

type server struct {
}


func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	fmt.Println("Broker iniciado")
	log.Printf("Recibido: %v", in.GetComandob())
	//------ip aleatoria asignada---------
	rand.Seed(time.Now().Unix())
	ips := []string{dns1, dns2, dns3}
	n := rand.Int() % len(ips)
	pick := ips[n]
	return &pb.CreateBReply{Ipb: pick}, nil
}
//-------------no implementados----------------
func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	return &pb.CreateDReply{Reloj: "null"}, nil
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCrudServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
