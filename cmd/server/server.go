package main

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/grpc_server"
	"github.com/openinfradev/tks-common/pkg/log"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	port               int
	tls                bool
	tlsClientCertPath  string
	tlsCertPath        string
	tlsKeyPath         string

	dbhost             string
	dbport             string
	dbuser             string
	dbpassword         string
)

func init() {
	flag.IntVar(&port, "port", 9111, "service port")
	flag.BoolVar(&tls, "tls", false, "enabled tls")
	flag.StringVar(&tlsClientCertPath, "tls-client-cert-path", "../../cert/tks-ca.crt", "path of ca cert file for tls")
	flag.StringVar(&tlsCertPath, "tls-cert-path", "../../cert/tks-server.crt", "path of cert file for tls")
	flag.StringVar(&tlsKeyPath, "tls-key-path", "../../cert/tks-server.key", "path of key file for tls")
	flag.StringVar(&dbhost, "dbhost", "localhost", "host of postgreSQL")
	flag.StringVar(&dbport, "dbport", "5432", "port of postgreSQL")
	flag.StringVar(&dbuser, "dbuser", "postgres", "postgreSQL user")
	flag.StringVar(&dbpassword, "dbpassword", "password", "password for postgreSQL user")
}

func main() {
	flag.Parse()

	log.Info("*** Arguments *** ")
	log.Info("port : ", port)
	log.Info("tls : ", tls)
	log.Info("tlsClientCertPath : ", tlsClientCertPath)
	log.Info("tlsCertPath : ", tlsCertPath)
	log.Info("tlsKeyPath : ", tlsKeyPath)
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

	// start server
	s, conn, err := grpc_server.CreateServer(port, tls, tlsCertPath, tlsKeyPath)
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
