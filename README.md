# tks-info

[![Go Report Card](https://goreportcard.com/badge/github.com/openinfradev/tks-info?style=flat-square)](https://goreportcard.com/report/github.com/openinfradev/tks-info)
[![Go Reference](https://pkg.go.dev/badge/github.com/openinfradev/tks-info.svg)](https://pkg.go.dev/github.com/openinfradev/tks-info)
[![Release](https://img.shields.io/github/release/sktelecom/tks-info.svg?style=flat-square)](https://github.com/openinfradev/tks-info/releases/latest)

TKS는 Taco Kubernetes Service의 약자로, SK Telecom이 만든 GitOps기반의 서비스 시스템을 의미합니다. 그 중 tks-info는 사용자 클러스터 및 서비스의 정보를 다루는 서비스이며, 다른 tks service들과 gRPC 기반으로 통신합니다. 

RPC 호출을 위한 proto 파일은 [tks-proto](https://github.com/openinfradev/tks-proto)에서 확인할 수 있습니다.


## Quick Start

### Prerequisite
* docker 20.x 설치
* tks-contract 설치
  * tks-contract: https://github.com/openinfradev/tks-contract
* postgresql을 설치하고 database를 초기화합니다.
  ```
  $ docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=password -d postgres
  $ docker cp scripts/script.sql postgres:/script.sql
  $ docker exec -ti postgres psql -U postgres -a -f script.sql
  ``` 

### 서비스 구동 
#### For go developers
```
$ go build -o bin/tks-info ./cmd/server/
$ bin/tks-info -port 9110
```

#### For docker users
```
$ docker pull sktcloud/tks-info:latest
$ docker run --name tks-info -p 9110:9110 -d \
   sktcloud/tks-info:latest -port 9110 
```

### gRPC API 호출 예제 (golang)

```go
import (
  "google.golang.org/grpc"
  pb "github.com/openinfradev/tks-proto/tks_pb"
  "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/types/known/timestamppb"
  ...
)
  func YOUR_FUNCTION(YOUR_PARAMS...) {
    var conn *grpc.ClientConn
    tksInfoUrl = viper.GetString("tksInfoUrl")
    if tksInfoUrl == "" {
      fmt.Println("You must specify tksInfoUrl at config file")
      os.Exit(1)
    }
    conn, err := grpc.Dial(tksInfoUrl, grpc.WithInsecure())
    if err != nil {
      log.Fatalf("Couldn't connect to tks-info: %s", err)
    }
    defer conn.Close()
    client := pb.NewClusterInfoServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()
    data := pb.GetClustersRequest{}
    data.ContractId = viper.GetString("contractId")
    m := protojson.MarshalOptions{
      Indent:        "  ",
      UseProtoNames: true,
    }
    jsonBytes, _ := m.Marshal(&data)
    r, err := client.GetClusters(ctx, &data)
    if err != nil {
      fmt.Println(err)
    } else {
      if len(r.Clusters) == 0 {
        fmt.Println("No cluster exists for specified contract!")
      } else {
        /* print cluster info from 'r' */
 
    }
  }
}

```
