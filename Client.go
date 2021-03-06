package main

import (
	"context"
	"log"
	"time"
	"fmt"
	"strings"
	"os"
	"bufio"

	"google.golang.org/grpc"
	pb "github.com/GianniCarlini/Lab-3-SD/proto"
)

const (
	address = "10.10.28.67:50050" //Broker
	dns1 = "10.10.28.68:50051"
	dns2 = "10.10.28.69:50052"
	dns3 = "10.10.28.7:50053"
)
type server struct {
}

var ipcambio int64 = 0

//-------------no implementados----------------
func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	return &pb.CreateBReply{Ipb: "null"}, nil
}
func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	reloj := []int64{1,0,0}
	return &pb.CreateDReply{Reloj:reloj}, nil
}
func (s *server) Merge(ctx context.Context, in *pb.MergeRequest) (*pb.MergeReply, error) {
	r := []byte("Hola mundo!\n")
	return &pb.MergeReply{Logresp: r}, nil
}
func (s *server) PMerge(ctx context.Context, in *pb.PMergeRequest) (*pb.PMergeReply, error) {
	return &pb.PMergeReply{Mresp: "Gracias!"}, nil
}

func main() {
	fmt.Println("Bienvenido Cliente")
	for{
		fmt.Println("Ingrese el comando")
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		option := strings.Split(comando," ")[0]

		if (strings.ToLower(option) == ("get")){
			fmt.Println("Conectando con el broker")
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock()) //broker
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewCrudClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.ConnectC(ctx, &pb.ConnectCRequest{Comandoc: comando})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Ip: %s", r.GetIpc())
			fmt.Println("Reloj:",r.GetRelojc())
			//--------------------------------------------
		}else{
			fmt.Println("Comando ingresado no valido")
		}
	}
}
