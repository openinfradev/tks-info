package client

import (
	"context"
	"sync"
	"fmt"
	"time"

	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//"google.golang.org/grpc/credentials"

	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

var (
	once sync.Once
	clusterInfoClient pb.ClusterInfoServiceClient
	cspInfoClient pb.CspInfoServiceClient
)

// FOR_TEST
func RequestLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		end := time.Now()

		log.Info(fmt.Sprintf("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("time.RFC3339"), end.Format(time.RFC3339), err))
		
		return err
	}
}

func GetClusterInfoClient(address string, port int, caller string) pb.ClusterInfoServiceClient {
	host := fmt.Sprintf("%s:%d", address, port)
	once.Do(func() {
		conn, _ := grpc.Dial(
			host,
			grpc.WithInsecure(),
		)

		clusterInfoClient = pb.NewClusterInfoServiceClient(conn)
	})
	return clusterInfoClient
}

func GetCspInfoClient(address string, port int, caller string) pb.CspInfoServiceClient {
	host := fmt.Sprintf("%s:%d", address, port)
	once.Do(func() {
		conn, _ := grpc.Dial(
			host,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(
					RequestLogging(),
				),
			),
		)

		cspInfoClient = pb.NewCspInfoServiceClient(conn)
	})
	return cspInfoClient
}
