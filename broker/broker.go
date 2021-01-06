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
	address = "10.10.28.67:50050" //Broker
	dns1 = "10.10.28.68:50051"
	dns2 = "10.10.28.69:50052"
	dns3 = "10.10.28.7:50053"
)

type server struct {
}

var ipcambio int64 = 0

func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	fmt.Println("Broker iniciado")
	fmt.Println(ipcambio)
	fmt.Println("Broker inisaddsadasciado")
	log.Printf("Recibido: %v", in.GetComandob())
	//------ip aleatoria asignada---------
	rand.Seed(time.Now().Unix())
	ips := []string{dns1, dns2, dns3}
	n := rand.Int() % len(ips)
	pick := ips[n]
	return &pb.CreateBReply{Ipb: pick, Contador: ipcambio}, nil
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

func (s *server) IpCambio(ctx context.Context, in *pb.IpCambioRequest) (*pb.IpCambioReply, error) {
	fmt.Println(in.GetCambio())
	ipcambio++
	fmt.Println("Hicieron un merge")
	fmt.Println(ipcambio)
	return &pb.IpCambioReply{Recibido: "Gracias! me llego el merge"}, nil
}
//-------------no implementados----------------
func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	reloj := []int64{1,0,0}
	return &pb.CreateDReply{Reloj: reloj}, nil
}
func (s *server) RelojCambio(stream pb.Crud_RelojCambioServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		asd := in.Domain
		fmt.Println(asd)
		resp := pb.RelojCambioReply{Aviso: "Hice un cambio en el reloj"}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
	}
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
