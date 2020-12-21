package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"strings"
	"io/ioutil"


	"google.golang.org/grpc"
	pb "github.com/GianniCarlini/Lab-3-SD/proto"
)

const (
	port = ":50051"
	broker = "localhost:50050"
	dns1 = "localhost:50051"
	dns2 = "localhost:50052"
	dns3 = "localhost:50053"
)

type server struct {
}

var dom [] string
var watch [][] int64

func existeEnArreglo(arreglo []string, busqueda string) (bool,int) { //https://parzibyte.me/blog/2019/08/07/go-elemento-existe-en-arreglo/
	for i, numero := range arreglo {
		if numero == busqueda {
			return true,i
		}
	}
	return false,0
}

func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	fmt.Println("DNS 1 iniciado")
	log.Printf("Recibido: %v", in.GetComandod())
	//------archivo---------
	comando := in.GetComandod()
	option := strings.Split(comando," ")[0]
	nd := strings.Split(comando," ")[1]
	nd = strings.TrimSuffix(nd, "\n")
	domain := strings.Split(nd,".")[1]

	//--------------------------------------------
	line := comando
	content, err := ioutil.ReadFile("log.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	
	content = append(content, []byte(line)...)

	err = ioutil.WriteFile(("log.txt"), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//----------------------------------------------------------
	existe,indice := existeEnArreglo(dom,domain)
	if existe {
		watch[indice][0] += 1 //por dominio
	}else{
		watch = append(watch,[]int64{1,0,0})
		dom = append(dom,domain)
	}
	fmt.Println(watch)
	fmt.Println(dom)
	//----------------------------------------------------------------------

	if strings.ToLower(option) == "create"{
		linea := nd+" IN A "+dns1+"\n"
		content, err := ioutil.ReadFile(domain+".txt") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		
		content = append(content, []byte(linea)...)

		err = ioutil.WriteFile((domain+".txt"), content, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}else if strings.ToLower(option) == "update"{
		cambio := strings.Split(comando," ")[2]
		cambio = strings.TrimSuffix(cambio, "\n")
		input, err := ioutil.ReadFile(domain+".txt")
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for i, line := range lines {
                if strings.Contains(line, nd) {
                        lines[i] = cambio+" IN A "+dns1
                }
        }
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(domain+".txt", []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
        }
	}else if strings.ToLower(option) == "delete"{
		input, err := ioutil.ReadFile(domain+".txt")
        if err != nil {
                log.Fatalln(err)
        }

        lines := strings.Split(string(input), "\n")

        for i, line := range lines {
                if strings.Contains(line, nd) {
                        lines[i] = " "
                }
        }
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(domain+".txt", []byte(output), 0644)
        if err != nil {
                log.Fatalln(err)
		}
	}
	fmt.Println("ANTES DE LA TRAGEDIA")
	fmt.Println(watch[indice])
	return &pb.CreateDReply{Reloj: watch[indice]}, nil
}
//-------------no implementados----------------
func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	return &pb.CreateBReply{Ipb: "null"}, nil
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
