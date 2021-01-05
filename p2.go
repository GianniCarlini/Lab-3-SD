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
	address = "localhost:50050" //Broker
	dns1 = "localhost:50051"
	dns2 = "localhost:50052"
	dns3 = "localhost:50053"
)

type server struct {
}

var ipcambio int64 = 9999999999
var ipconect string

//-----no imp--------------
func (s *server) Merge(ctx context.Context, in *pb.MergeRequest) (*pb.MergeReply, error) {
	r := []byte("Hola mundo!\n")
	return &pb.MergeReply{Logresp: r}, nil
}
func (s *server) PMerge(ctx context.Context, in *pb.PMergeRequest) (*pb.PMergeReply, error) {
	return &pb.PMergeReply{Mresp: "Gracias!"}, nil
}
func main() {
	fmt.Println("Bienvenido Administrador")
	for{
		fmt.Println("Ingrese el comando")
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		option := strings.Split(comando," ")[0]

		if (strings.ToLower(option) == ("create")) || (strings.ToLower(option) == ("update")) || (strings.ToLower(option) == ("delete")){
			fmt.Println("Conectando con el broker")
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())//broker
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewCrudClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.CreateB(ctx, &pb.CreateBRequest{Comandob: comando})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Ip: %s", r.GetIpb())
			fmt.Println("Connectando con el servidor DNS")
			if r.GetContador() == ipcambio {
				fmt.Println("conectando...")
			}else{
				ipconect = r.GetIpb()
				ipcambio = r.GetContador()
			}
			//--------------------------------------------
			conn2, err2 := grpc.Dial(dns2, grpc.WithInsecure(), grpc.WithBlock()) //dns
			if err2 != nil {
				log.Fatalf("did not connect: %v", err2)
			}
			defer conn2.Close()
			c2 := pb.NewCrudClient(conn2)
			ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r2, err3 := c2.CreateD(ctx2, &pb.CreateDRequest{Comandod: comando})
			if err3 != nil {
				log.Fatalf("could not greet: %v", err3)
			}
			fmt.Println("Reloj:",r2.GetReloj())
		}else{
			fmt.Println("Comando ingresado no valido")
		}
	}
}
