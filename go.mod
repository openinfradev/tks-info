module github.com/openinfradev/tks-info

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/openinfradev/tks-contract v0.1.1-0.20210928021110-fe2b666327cc
	github.com/openinfradev/tks-proto v0.0.6-0.20210924020717-178698d59e9d
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20210927181540-4e4d966f7476 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	google.golang.org/genproto v0.0.0-20210927142257-433400c27d05 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gorm.io/datatypes v1.0.2
	gorm.io/driver/postgres v1.1.1
	gorm.io/gorm v1.21.15
)

replace github.com/openinfradev/tks-info => ./
replace github.com/openinfradev/tks-proto => ../tks-proto2
