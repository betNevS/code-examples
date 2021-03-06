package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/betNevS/code-examples/protobuf/pcbook/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/betNevS/code-examples/protobuf/pcbook/pb"
	"google.golang.org/grpc"
)

func main() {
	serverAddr := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddr)

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	laptopClient := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	searchLaptop(laptopClient, filter)
}

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Println("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		laptop := res.GetLaptop()
		log.Println("- found: ", laptop.GetId())
		log.Println("  + brand: ", laptop.GetBrand())
		log.Println("  + name: ", laptop.GetName())
		log.Println("  + cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Println("  + cpu min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Println("  + ram: ", laptop.GetMemory().GetValue(), laptop.GetMemory().GetUnit())
		log.Println("  + price: ", laptop.GetPriceUsd(), "usd")
	}
}

func createLaptop(client pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	laptop.Id = ""
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
	}
	log.Printf("created laptop with id: %v", res.Id)
}
