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
	"time"

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
func merge(){
    for{
        time.Sleep(time.Duration(60)*time.Second)
		fmt.Println("Iniciando merge")
		//----------log 2-----------------------
		conn, err := grpc.Dial(dns2, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewCrudClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Merge(ctx, &pb.MergeRequest{Peticionlog: "gimme your log"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		bit := r.GetLogresp()
		err2 := ioutil.WriteFile("log2.txt", bit, 0644)
		if err2 != nil {
			log.Fatal(err2)
		}
		
		//----------log 3-----------------------
		conn2, err3 := grpc.Dial(dns3, grpc.WithInsecure(), grpc.WithBlock())
		if err3 != nil {
			log.Fatalf("did not connect: %v", err3)
		}
		defer conn2.Close()
		c2 := pb.NewCrudClient(conn2)
		ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r2, err4 := c2.Merge(ctx2, &pb.MergeRequest{Peticionlog: "gimme your log"})
		if err4 != nil {
			log.Fatalf("could not greet: %v", err4)
		}
		bit2 := r2.GetLogresp()
		err5 := ioutil.WriteFile("log3.txt", bit2, 0644)
		if err5 != nil {
			log.Fatal(err5)
		}
		//-------magia-----------------------
		var create []string
		var update []string
		var delete []string

		var create1 []string
		var update1 []string
		var delete1 []string

		var mergecito []string
		//----------------------------------------------
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
				create = append(create, scanner.Text())
				create1 = append(create1, scanner.Text())
			}else if strings.ToLower(option) == "update"{
				update = append(update, scanner.Text())
				update1 = append(update1, scanner.Text())
			}else if strings.ToLower(option) == "delete"{
				delete = append(delete, scanner.Text())
				delete1 = append(delete1, scanner.Text())
			}
		}
		//----------------------------------------------
		file2, err8 := os.Open("log2.txt")
		if err8 != nil {
			log.Fatal(err8)
		}
		defer file2.Close()
		scanner2 := bufio.NewScanner(file2)
		for scanner2.Scan() {             // internally, it advances token based on sperator
			fmt.Println(scanner2.Text()) 
			option := strings.Split(scanner2.Text()," ")[0]
			if strings.ToLower(option) == "create"{
				create = append(create, scanner2.Text())
			}else if strings.ToLower(option) == "update"{
				update = append(update, scanner2.Text())
			}else if strings.ToLower(option) == "delete"{
				delete = append(delete, scanner2.Text())
			}
			nd := strings.Split(scanner2.Text()," ")[1]
			nd = strings.TrimSuffix(nd, "\n")
			domain := strings.Split(nd,".")[1]
			existe,indice := existeEnArreglo(dom,domain)
			if existe {
				watch[indice][1] += 1 //por dominio
			}else{
				watch = append(watch,[]int64{0,1,0})
				dom = append(dom,domain)
			}
		}
		//----------------------------------------------
		file3, err9 := os.Open("log3.txt")
		if err9 != nil {
			log.Fatal(err9)
		}
		defer file3.Close()
 
		scanner3 := bufio.NewScanner(file3)
		for scanner3.Scan() {             // internally, it advances token based on sperator
			fmt.Println(scanner3.Text()) 
			option := strings.Split(scanner3.Text()," ")[0]
			if strings.ToLower(option) == "create"{
				create = append(create, scanner3.Text())
			}else if strings.ToLower(option) == "update"{
				update = append(update, scanner3.Text())
			}else if strings.ToLower(option) == "delete"{
				delete = append(delete, scanner3.Text())
			}
			nd := strings.Split(scanner3.Text()," ")[1]
			nd = strings.TrimSuffix(nd, "\n")
			domain := strings.Split(nd,".")[1]
			existe,indice := existeEnArreglo(dom,domain)
			if existe {
				watch[indice][2] += 1 //por dominio
			}else{
				watch = append(watch,[]int64{0,0,1})
				dom = append(dom,domain)
			}
		}
		for _,x := range create{
			mergecito = append(mergecito, x)
		}
		for _,x := range update{
			mergecito = append(mergecito, x)
		}
		for _,x := range delete{
			mergecito = append(mergecito, x)
		}
		fmt.Println(mergecito)
		//--------envio log a 2---------
		conn5, err11 := grpc.Dial(dns2, grpc.WithInsecure(), grpc.WithBlock())
		if err11 != nil {
			log.Fatalf("did not connect: %v", err11)
		}
		defer conn5.Close()
		c5 := pb.NewCrudClient(conn5)
		ctx5, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r5, err22 := c5.PMerge(ctx5, &pb.PMergeRequest{Mergecito: mergecito})
		if err22 != nil {
			log.Fatalf("could not greet: %v", err22)
		}
		fmt.Println("Respuesta merge2:",r5.GetMresp())
		//--------envio log a 3---------
		conn6, err111 := grpc.Dial(dns3, grpc.WithInsecure(), grpc.WithBlock())
		if err111 != nil {
			log.Fatalf("did not connect: %v", err111)
		}
		defer conn6.Close()
		c55 := pb.NewCrudClient(conn6)
		ctx6, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r6, err222 := c55.PMerge(ctx6, &pb.PMergeRequest{Mergecito: mergecito})
		if err222 != nil {
			log.Fatalf("could not greet: %v", err222)
		}
		fmt.Println("Respuesta merge3:",r6.GetMresp())
		//--------creo mis archivos---------
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
		//----------reloj 2-----------------------
		conn133, err133 := grpc.Dial(dns2, grpc.WithInsecure())
		if err133 != nil {
			log.Fatalf("failed to connect: %s", err133)
		}
		defer conn133.Close()

		client := pb.NewCrudClient(conn133)
		stream, _ := client.RelojCambio(context.Background())
		for i,_ := range dom{
			msg := &pb.RelojCambioRequest{Relojito : watch[i], Domain: dom[i]}
			stream.Send(msg)
			resp, err1332 := stream.Recv()
			if err1332 != nil {
				log.Fatalf("can not receive %v", err1332)
			}
			fmt.Println(resp.Aviso)
		}
		stream.CloseSend()
		//----------reloj 3-----------------------
		conn1333, err1333 := grpc.Dial(dns3, grpc.WithInsecure())
		if err1333 != nil {
			log.Fatalf("failed to connect: %s", err1333)
		}
		defer conn1333.Close()

		client3 := pb.NewCrudClient(conn1333)
		stream3, _ := client3.RelojCambio(context.Background())
		for i,_ := range dom{
			msg3 := &pb.RelojCambioRequest{Relojito : watch[i], Domain: dom[i]}
			stream3.Send(msg3)
			resp3, err13323 := stream3.Recv()
			if err13323 != nil {
				log.Fatalf("can not receive %v", err13323)
			}
			fmt.Println(resp3.Aviso)
		}
		stream3.CloseSend()
		//-------borrado de logs-------------
		err6 := os.Remove("log2.txt")
		if err6 != nil {
		  fmt.Printf("Error eliminando archivo: %v\n", err6)
		} else {
		  fmt.Println("Eliminado correctamente")
		}
		
		err7 := os.Remove("log3.txt")
		if err7 != nil {
		  fmt.Printf("Error eliminando archivo: %v\n", err7)
		} else {
		  fmt.Println("Eliminado correctamente")
		}

		err777 := os.Remove("log.txt")
		if err777 != nil {
		  fmt.Printf("Error eliminando archivo: %v\n", err777)
		} else {
		  fmt.Println("Eliminado correctamente")
		}
		_, err12312 := os.Create("log.txt")
		if err12312 != nil {
			fmt.Printf("Error eliminando archivo: %v\n",err12312)
		}
		//-----aviso al broker-------------
		conn63, err1112 := grpc.Dial("localhost:50050", grpc.WithInsecure(), grpc.WithBlock())
		if err1112 != nil {
			log.Fatalf("did not connect: %v", err1112)
		}
		defer conn63.Close()
		c553 := pb.NewCrudClient(conn63)
		ctx63, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r63, err2221 := c553.IpCambio(ctx63, &pb.IpCambioRequest{Cambio: "Hice un merge"})
		if err2221 != nil {
			log.Fatalf("could not greet: %v", err2221)
		}
		fmt.Println("Respuesta broker:",r63.GetRecibido())
		//--------------asd-------------------
		fmt.Println("Merge watch")
		fmt.Println(watch)
		fmt.Println(dom)
    }
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
//-------------no implementados----------------
func (s *server) CreateB(ctx context.Context, in *pb.CreateBRequest) (*pb.CreateBReply, error) {
	return &pb.CreateBReply{Ipb: "null"}, nil
}
func (s *server) ConnectC(ctx context.Context, in *pb.ConnectCRequest) (*pb.ConnectCReply, error) {
	reloj := []int64{1,0,0}
	return &pb.ConnectCReply{Ipc: "ip", Relojc: reloj}, nil
}
func (s *server) Merge(ctx context.Context, in *pb.MergeRequest) (*pb.MergeReply, error) {
	r := []byte("Hola mundo!\n")
	return &pb.MergeReply{Logresp: r}, nil
}
func (s *server) PMerge(ctx context.Context, in *pb.PMergeRequest) (*pb.PMergeReply, error) {
	return &pb.PMergeReply{Mresp: "Gracias!"}, nil
}

func (s *server) IpCambio(ctx context.Context, in *pb.IpCambioRequest) (*pb.IpCambioReply, error) {
	return &pb.IpCambioReply{Recibido: "Gracias!"}, nil
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
	}
}
func main() {

	go merge()

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
