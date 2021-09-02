module github.com/openinfradev/tks-info

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/lib/pq v1.10.2 // indirect
	github.com/openinfradev/tks-contract v0.1.1-0.20210902134454-132819708ac3
	github.com/openinfradev/tks-proto v0.0.6-0.20210901093202-5e0db3fa3d4f
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/openinfradev/tks-info => ./
