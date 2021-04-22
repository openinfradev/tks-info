package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-info/pkg/cert"
	pb "github.com/sktelecom/tks-proto/pbgo"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

	doAddApp(cc)
	doDeleteApp(cc)
	doGetAppIDs(cc)
	doGetAllAppsByClusterID(cc)
	doGetAppsByName(cc)
	doGetAppsByType(cc)
	doGetApp(cc)
	doUpdateApp(cc)
	doUpdateAppStatus(cc)
	doUpdateEndpoints(cc)
}

func doAddApp(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.AddAppRequest{
		ClusterId: "ccc",
		ServiceApp: &pb.ServiceApp{
			AppName: "my_service_mesh",
			Type:    pb.AppType_SERVICE_MESH,
			Owner:   "ccc",
			Status:  pb.AppStatus_APP_RUNNING,
			Endpoints: []*pb.Endpoint{
				{
					Type: pb.EpType_KIALI,
					Url:  "kiali.istio-system.svc.cluster.k1",
				},
				{
					Type: pb.EpType_JAEGER,
					Url:  "jaeger.istio-system.svc.cluster.k1",
				},
			},
			ExternalLabel: "service_mesh",
			CreatedTs:     timestamppb.New(time.Now()),
		},
	}
	res, err := c.AddApp(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling AddApp RPC", err)
	}
	log.Info("Response from AddApp: ", res.GetId())
}

func doDeleteApp(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.DeleteAppRequest{
		ClusterId: "ccc",
		AppId:     "111",
	}
	res, err := c.DeleteApp(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling DeleteApp RPC", err)
	}
	log.Info("Response from DeleteApp: ", res.GetCode())
}

func doGetAppIDs(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.IDRequest{
		Id: "ccc",
	}
	res, err := c.GetAppIDs(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppId RPC", err)
	}
	log.Info("Response from GetAppId: ", res.GetIds())
}

func doGetAllAppsByClusterID(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.IDRequest{
		Id: "ccc",
	}
	res, err := c.GetAllAppsByClusterID(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAllAppsByClusterID RPC", err)
	}
	log.Info("Response from GetAllAppsByClusterID: ", res)
}

func doGetAppsByName(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.GetAppsRequest{
		ClusterId: "ccc",
		AppName:   "my",
	}
	res, err := c.GetAppsByName(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppsByName RPC", err)
	}
	log.Info("Response from GetAppsByName: ", res)
}

func doGetAppsByType(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.GetAppsRequest{
		ClusterId: "ccc",
		Type:      pb.AppType_SERVICE_MESH,
	}
	res, err := c.GetAppsByType(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppsByType RPC", err)
	}
	log.Info("Response from GetAppsByType: ", res)
}

func doGetApp(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.GetAppRequest{
		ClusterId: "ccc",
		AppId:     "111",
	}
	res, err := c.GetApp(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetApp RPC", err)
	}
	log.Info("Response from GetApp: ", res)
}

func doUpdateApp(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.UpdateAppRequest{
		ClusterId: "ccc",
		ServiceApp: &pb.ServiceApp{
			AppId:   "111",
			AppName: "my_service_mesh",
			Type:    pb.AppType_SERVICE_MESH,
			Owner:   "ccc",
			Status:  pb.AppStatus_APP_RUNNING,
			Endpoints: []*pb.Endpoint{
				{
					Type: pb.EpType_KIALI,
					Url:  "kiali.istio-system.svc.cluster.k1",
				},
				{
					Type: pb.EpType_JAEGER,
					Url:  "jaeger.istio-system.svc.cluster.k1",
				},
			},
			ExternalLabel: "service_mesh",
			LastUpdatedTs: timestamppb.New(time.Now()),
		},
	}
	res, err := c.UpdateApp(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling UpdateApp RPC", err)
	}
	log.Info("Response from UpdateApp: ", res.GetCode())
}

func doUpdateAppStatus(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.UpdateAppStatusRequest{
		ClusterId: "ccc",
		AppId:     "111",
		Status:    pb.AppStatus_APP_RUNNING,
	}
	res, err := c.UpdateAppStatus(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling UpdateAppStatus RPC", err)
	}
	log.Info("Response from UpdateAppStatus: ", res.GetCode())
}

func doUpdateEndpoints(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.UpdateEndpointsRequest{
		ClusterId: "ccc",
		AppId:     "111",
		Endpoints: []*pb.Endpoint{
			{
				Type: pb.EpType_KIALI,
				Url:  "kiali.istio-system.svc.cluster.k1",
			},
			{
				Type: pb.EpType_JAEGER,
				Url:  "jaeger.istio-system.svc.cluster.k1",
			},
		},
	}
	res, err := c.UpdateEndpoints(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling UpdateEndpoints RPC", err)
	}
	log.Info("Response from UpdateEndpoints: ", res.GetCode())
}
