package main

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/grpc_client"
	"github.com/openinfradev/tks-common/pkg/grpc_server"
	"github.com/openinfradev/tks-common/pkg/log"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	port              int
	tlsEnabled        bool
	tlsClientCertPath string
	tlsCertPath       string
	tlsKeyPath        string

	contractAddress string
	contractPort    int
	dbhost          string
	dbport          string
	dbuser          string
	dbpassword      string
)

var (
	contractClient pb.ContractServiceClient
)

func init() {
	flag.IntVar(&port, "port", 9111, "service port")
	flag.BoolVar(&tlsEnabled, "tlsEnabled", false, "enabled tls")
	flag.StringVar(&tlsClientCertPath, "tls-client-cert-path", "../../cert/tks-ca.crt", "path of ca cert file for tls")
	flag.StringVar(&tlsCertPath, "tls-cert-path", "../../cert/tks-server.crt", "path of cert file for tls")
	flag.StringVar(&tlsKeyPath, "tls-key-path", "../../cert/tks-server.key", "path of key file for tls")
	flag.StringVar(&contractAddress, "contract-address", "localhost", "service address for tks-contract")
	flag.IntVar(&contractPort, "contract-port", 9110, "service port for tks-contract")
	flag.StringVar(&dbhost, "dbhost", "localhost", "host of postgreSQL")
	flag.StringVar(&dbport, "dbport", "5432", "port of postgreSQL")
	flag.StringVar(&dbuser, "dbuser", "postgres", "postgreSQL user")
	flag.StringVar(&dbpassword, "dbpassword", "password", "password for postgreSQL user")
}

func main() {
	flag.Parse()

	log.Info("*** Arguments *** ")
	log.Info("port : ", port)
	log.Info("tlsEnabled : ", tlsEnabled)
	log.Info("tlsClientCertPath : ", tlsClientCertPath)
	log.Info("tlsCertPath : ", tlsCertPath)
	log.Info("tlsKeyPath : ", tlsKeyPath)
	log.Info("contractAddress : ", contractAddress)
	log.Info("contractPort : ", contractPort)
	log.Info("dbhost : ", dbhost)
	log.Info("dbport : ", dbport)
	log.Info("dbuser : ", dbuser)
	log.Info("dbpassword : ", dbpassword)
	log.Info("****************** ")

	// initialize database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=tks port=%s sslmode=disable TimeZone=Asia/Seoul",
		dbhost, dbuser, dbpassword, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open database ", err)
	}

	// initialize handlers
	InitAppInfoHandler(db)
	InitClusterInfoHandler(db)
	InitCspInfoHandler(db)
	InitKeycloakInfoHandler(db)

	// initialize clients
	if _, contractClient, err = grpc_client.CreateContractClient(contractAddress, contractPort, tlsEnabled, tlsClientCertPath); err != nil {
		log.Fatal("failed to create contract client : ", err)
	}

	// start server
	s, conn, err := grpc_server.CreateServer(port, tlsEnabled, tlsCertPath, tlsKeyPath)
	if err != nil {
		log.Fatal("failed to crate grpc_server : ", err)
	}

	pb.RegisterAppInfoServiceServer(s, &AppInfoServer{})
	pb.RegisterClusterInfoServiceServer(s, &ClusterInfoServer{})
	pb.RegisterCspInfoServiceServer(s, &CspInfoServer{})
	pb.RegisterKeycloakInfoServiceServer(s, &KeycloakInfoServer{})
	if err := s.Serve(conn); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
