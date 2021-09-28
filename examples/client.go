package main

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/google/uuid"
	"github.com/openinfradev/tks-contract/pkg/log"
	"github.com/openinfradev/tks-info/pkg/cert"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	port       = flag.Int("port", 9111, "The gRPC server port")
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile     = flag.String("ca_file", "", "The TLS ca file")
	contractID = uuid.New().String()
	cspID      = uuid.New().String()
	clusterID  string
	appGroupID string
	clusterName       = "testCluster"
)

func main() {
	log.Info("Hello I'm a application client")
	log.Info("new Cluster ID: ", clusterID)

	flag.Parse()

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

	doCreateCluster(cc)
	doCreateAppGroup(cc)
	doGetAppGroupsByClusterID(cc)
	doGetAppGroups(cc)
	doGetAppGroup(cc)
	doUpdateAppGroupStatus(cc)
	doUpdateApp(cc)
	doGetAppsByAppGroupID(cc)
	doDeleteAppGroup(cc)
}

func doCreateCluster(cc *grpc.ClientConn) {
	c := pb.NewClusterInfoServiceClient(cc)

	dummyConf := pb.ClusterConf{
		MasterFlavor: "tiny",
		MasterReplicas: 3,
		MasterRootSize: 50,
		WorkerFlavor: "medium",
		WorkerReplicas: 5,
		WorkerRootSize: 50,
		K8SVersion: "1.18.8",
	}

	req := &pb.AddClusterInfoRequest{
		ContractId: contractID,
		CspId: cspID,
		Name: clusterName,
		Conf: &dummyConf,
	}

	res, err := c.AddClusterInfo(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling AddClusterInfo RPC", err)
	}

	clusterID = res.GetId()
	log.Info("Response from AddClusterInfo: ", clusterID)
}

func doCreateAppGroup(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.CreateAppGroupRequest{
		ClusterId: clusterID,
		AppGroup: &pb.AppGroup{
			AppGroupName:  "my_service_mesh",
			Type:          pb.AppGroupType_SERVICE_MESH,
			Status:        pb.AppGroupStatus_APP_GROUP_INSTALLING,
			ExternalLabel: "service_mesh",
		},
	}
	res, err := c.CreateAppGroup(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling CreateAppGroup RPC", err)
	}
	appGroupID = res.GetId()
	log.Info("Response from CreateAppGroup: ", res.GetId())
}

func doGetAppGroupsByClusterID(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.IDRequest{
		Id: clusterID,
	}
	res, err := c.GetAppGroupsByClusterID(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppGroupsByClusterID RPC", err)
	}
	log.Info("Response from GetAppGroupsByClusterID: ", res)
}

func doGetAppGroups(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.GetAppGroupsRequest{
		AppGroupName: "my_service_mesh",
		Type:         pb.AppGroupType_SERVICE_MESH,
	}
	res, err := c.GetAppGroups(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppGroups RPC", err)
	}
	log.Info("Response from GetAppGroups: ", res)
}

func doGetAppGroup(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.GetAppGroupRequest{
		AppGroupId: appGroupID,
	}
	res, err := c.GetAppGroup(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppGroup RPC", err)
	}
	log.Info("Response from GetAppGroup: ", res)
}

func doUpdateAppGroupStatus(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.UpdateAppGroupStatusRequest{
		AppGroupId: appGroupID,
		Status:     pb.AppGroupStatus_APP_GROUP_ERROR,
	}
	res, err := c.UpdateAppGroupStatus(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling UpdateAppGroupStatus RPC", err)
	}
	log.Info("Response from UpdateAppGroupStatus: ", res)
}

func doUpdateApp(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.UpdateAppRequest{
		AppGroupId: appGroupID,
		AppType:    pb.AppType_KIALI,
		Endpoint:   "https://localhost:20001",
		Metadata:   "",
	}
	res, err := c.UpdateApp(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling UpdateApp RPC", err)
	}
	log.Info("Response from UpdateApp: ", res.GetCode())
}

func doGetAppsByAppGroupID(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.IDRequest{
		Id: appGroupID,
	}
	res, err := c.GetAppsByAppGroupID(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling GetAppsByAppGroupID RPC", err)
	}
	log.Info("Response from GetAppsByAppGroupID: ")
	for _, app := range res.Apps {
		log.Info(fmt.Sprintf("id %s, type %d, endpoint %s", app.AppId, app.Type, app.Endpoint))
	}
}

func doDeleteAppGroup(cc *grpc.ClientConn) {
	c := pb.NewAppInfoServiceClient(cc)

	req := &pb.DeleteAppGroupRequest{
		AppGroupId: appGroupID,
	}
	res, err := c.DeleteAppGroup(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling DeleteApp RPC", err)
	}
	log.Info("Response from DeleteApp: ", res.GetCode())
}
