module github.com/openinfradev/tks-info

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/jackc/pgx/v4 v4.15.0 // indirect
	github.com/openinfradev/tks-common v0.0.0-20221122025625-be9f8957ec3c
	github.com/openinfradev/tks-proto v0.0.6-0.20230209014521-c44086e732d8
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.17.0 // indirect
	google.golang.org/genproto v0.0.0-20220211171837-173942840c17 // indirect
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.28.1
	gorm.io/datatypes v1.0.5
	gorm.io/driver/mysql v1.2.3 // indirect
	gorm.io/driver/postgres v1.2.3
	gorm.io/driver/sqlite v1.1.4 // indirect
	gorm.io/driver/sqlserver v1.0.9 // indirect
	gorm.io/gorm v1.22.5
)

replace github.com/openinfradev/tks-info => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
//replace github.com/openinfradev/tks-common => ../tks-common
//replace github.com/openinfradev/tks-contract => ../tks-contract
