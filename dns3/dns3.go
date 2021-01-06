package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"bufio"

	"google.golang.org/grpc"
	pb "github.com/GianniCarlini/Lab-3-SD/proto"
)

const (
	port = ":50053"
	broker = "10.10.28.67:50050" //Broker
	dns1 = "10.10.28.68:50051"
	dns2 = "10.10.28.69:50052"
	dns3 = "10.10.28.7:50053"
)

type server struct {
}

var dom [] string
var watch [][] int64

func existeEnArreglo(arreglo []string, busqueda string) (bool,int) { //https://parzibyte.me/blog/2019/08/07/go-elemento-existe-en-arreglo/
	var ind int
	for i,numero := range arreglo {
		if numero == busqueda {
			return true,i
		}else{
			ind = i+1
		}
	}
	return false,ind
}
func existeEnArreglo2(arreglo []string, busqueda string) bool{ //https://parzibyte.me/blog/2019/08/07/go-elemento-existe-en-arreglo/
	for _,numero := range arreglo {
		if numero == busqueda {
			return true
		}
	}
	return false
}
func (s *server) CreateD(ctx context.Context, in *pb.CreateDRequest) (*pb.CreateDReply, error) {
	fmt.Println("DNS 3 iniciado")
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
		watch[indice][2] += 1 //por dominio
	}else{
		watch = append(watch,[]int64{0,0,1})
		dom = append(dom,domain)
	}
	fmt.Println(watch)
	fmt.Println(dom)
	//----------------------------------------------------------------------

	if strings.ToLower(option) == "create"{
		linea := nd+" IN A "+dns3+"\n"
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
                        lines[i] = cambio+" IN A "+dns3
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
	return &pb.CreateDReply{Reloj: watch[indice]}, nil
}
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	comando := in.GetComandoget()
	nd := strings.Split(comando," ")[1]
	nd = strings.TrimSuffix(nd, "\n")
	dominio := strings.Split(nd,".")[1]
	fmt.Println(dominio)
	//-----abro--------------
	file, err := os.Open(dominio+".txt")
       if err != nil {
           log.Fatal(err)
       }
       defer file.Close()

       scanner := bufio.NewScanner(file)

	   var ip string
	   for scanner.Scan() {            // internally, it advances token based on sperator
			linea := scanner.Text()
			res := strings.Split(linea, " ") 
			if nd == res[0]{
				fmt.Println("encontre la linea")
				ip = res[3]
			}
           fmt.Println(scanner.Text())  // token in unicode-char
	   }
	   existe,indice := existeEnArreglo(dom,dominio)
	   if existe{
		   fmt.Println("nombre.dominio encontrado")
	   }
	   reloj := watch[indice]
	return &pb.GetReply{Ipget: ip, Relojget: reloj}, nil
}
func (s *server) Merge(ctx context.Context, in *pb.MergeRequest) (*pb.MergeReply, error) {
	b, err := ioutil.ReadFile("log.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return &pb.MergeReply{Logresp: b}, nil
}
func (s *server) PMerge(ctx context.Context, in *pb.PMergeRequest) (*pb.PMergeReply, error) {
	var create1 []string
	var update1 []string
	var delete1 []string
	mergecito := in.GetMergecito()

	file, err0 := os.Open("log.txt")
	if err0 != nil {
		log.Fatal(err0)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {             // internally, it advances token based on sperator
		fmt.Println(scanner.Text()) 
		option := strings.Split(scanner.Text()," ")[0]
		if strings.ToLower(option) == "create"{
			create1 = append(create1, scanner.Text())
		}else if strings.ToLower(option) == "update"{
			update1 = append(update1, scanner.Text())
		}else if strings.ToLower(option) == "delete"{
			delete1 = append(delete1, scanner.Text())
		}
	}

	for _,c := range mergecito{
		if (existeEnArreglo2(create1, c) == true) || (existeEnArreglo2(update1, c)== true) || (existeEnArreglo2(delete1, c)== true){
			fmt.Println("ya me ejecutaron")
		}else{
			option := strings.Split(c," ")[0]
			nd := strings.Split(c," ")[1]
			nd = strings.TrimSuffix(nd, "\n")
			domain := strings.Split(nd,".")[1]
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
				cambio := strings.Split(c," ")[2]
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
		}
	}
	err777 := os.Remove("log.txt")
	if err777 != nil {
	  fmt.Printf("Error eliminando archivo: %v\n", err777)
	} else {
	  fmt.Println("Eliminado correctamente")
	}
	_, err12312 := os.Create("log.txt")
	if err12312 != nil {
		fmt.Printf("Error eliminando archivo: %v\n", err12312)
	}
	watch = nil
	dom = nil
	return &pb.PMergeReply{Mresp: "Gracias!"}, nil
}
func (s *server) RelojCambio(stream pb.Crud_RelojCambioServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		watch = append(watch,in.Relojito)
		dom = append(dom,in.Domain)

		resp := pb.RelojCambioReply{Aviso: "Hice un cambio en el reloj"}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
		fmt.Println("Merge watch")
		fmt.Println(watch)
		fmt.Println(dom)
	}
}
//-------------no implementados----------------
func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	return &pb.CreateBReply{Ipb: "null"}, nil
}
func (s *server) ConnectC(ctx context.Context, in *pb.ConnectCRequest) (*pb.ConnectCReply, error) {
	reloj := []int64{1,0,0}
	return &pb.ConnectCReply{Ipc: "ip", Relojc: reloj}, nil
}

func (s *server) IpCambio(ctx context.Context, in *pb.IpCambioRequest) (*pb.IpCambioReply, error) {
	return &pb.IpCambioReply{Recibido: "Gracias!"}, nil
}
func main() {
	_, err12312 := os.Create("log.txt")
		if err12312 != nil {
			fmt.Printf("Error eliminando archivo: %v\n",err12312)
		}
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
