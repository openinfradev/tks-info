package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-contract/pkg/log"
	"github.com/openinfradev/tks-info/pkg/cert"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	port       = flag.Int("port", 9111, "The gRPC server port")
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	dbhost     = flag.String("dbhost", "localhost", "host of postgreSQL")
	dbport     = flag.String("dbport", "5432", "port of postgreSQL")
	dbuser     = flag.String("dbuser", "postgres", "postgreSQL user")
	dbpassword = flag.String("dbpassword", "password", "password for postgreSQL user")
)

func main() {
	log.Info("tksinfo server is starting...")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		// log.Fatalln("Failed to listen:", err)
		log.Fatal("failed to listen:", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = cert.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = cert.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatal("Failed to generate credentials", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s := grpc.NewServer(opts...)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=tks port=%s sslmode=disable TimeZone=Asia/Seoul",
		*dbhost, *dbuser, *dbpassword, *dbport)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	InitAppInfoHandler(db)
	InitClusterInfoHandler(db)
	InitCspInfoHandler(db)
	InitKeycloakInfoHandler(db)
	pb.RegisterAppInfoServiceServer(s, &AppInfoServer{})
	pb.RegisterClusterInfoServiceServer(s, &ClusterInfoServer{})
	pb.RegisterCspInfoServiceServer(s, &CspInfoServer{})
	pb.RegisterKeycloakInfoServiceServer(s, &KeycloakInfoServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
