package main

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-info/pkg/cert"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
	port   = flag.Int("port", 9111, "The gRPC server port")
	tls    = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile = flag.String("ca_file", "", "The TLS ca file")
)

func main() {
	log.Info("Hello I'm a application client")

	opts := grpc.WithInsecure()
	if *tls {
		if *caFile == "" {
			*caFile = cert.Path("x509/ca.crt")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, "")
		if err != nil {
			log.Fatal("Error while loading CA trust certificate: ", err)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	addr := fmt.Sprintf(":%d", *port)
	cc, err := grpc.Dial(addr, opts)
	if err != nil {
		log.Fatal("could not connect: ", err)
	}
	defer cc.Close()

	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.IDRequest{
		Id: "uuid",
	}
	res, err := c.GetAppIDs(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppId RPC", err)
	}
	log.Info("Response from GetAppId: ", res.GetIds())
}
