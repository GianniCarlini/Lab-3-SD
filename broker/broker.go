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
func (s *server) ConnectC(ctx context.Context, in *pb.ConnectCRequest) (*pb.ConnectCReply, error) {
	fmt.Println("Conectando con el broker metodo get")
	conn, err := grpc.Dial(dns1, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCrudClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Get(ctx, &pb.GetRequest{Comandoget: in.GetComandoc()})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return &pb.ConnectCReply{Ipc: r.GetIpget(), Relojc: r.GetRelojget()}, nil
}

//-------------no implementados----------------
func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	reloj := []int64{1,0,0}
	return &pb.CreateDReply{Reloj: reloj}, nil
}
//-------------no implementados----------------
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	reloj := []int64{1,0,0}
	return &pb.GetReply{Ipget: "ip", Relojget: reloj}, nil
}
func (s *server) Merge(ctx context.Context, in *pb.MergeRequest) (*pb.MergeReply, error) {
	r := []byte("Hola mundo!\n")
	return &pb.MergeReply{Logresp: r}, nil
}

func (s *server) PMerge(ctx context.Context, in *pb.PMergeRequest) (*pb.PMergeReply, error) {
	return &pb.PMergeReply{Mresp: "Gracias!"}, nil
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
