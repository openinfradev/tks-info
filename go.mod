module github.com/openinfradev/tks-info

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/openinfradev/tks-contract v0.1.1-0.20210915081037-2fef4d86b728
	github.com/openinfradev/tks-proto v0.0.6-0.20210901093202-5e0db3fa3d4f
	github.com/stretchr/testify v1.7.0
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/datatypes v1.0.2
	gorm.io/driver/postgres v1.1.1
	gorm.io/gorm v1.21.15
)

replace github.com/openinfradev/tks-info => ./

replace github.com/openinfradev/tks-proto => ../tks-proto

replace github.com/openinfradev/tks-contract => ../tks-contract
