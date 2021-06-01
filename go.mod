module github.com/sktelecom/tks-info

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/lib/pq v1.10.2 // indirect
	github.com/sktelecom/tks-contract v0.1.0
	github.com/sktelecom/tks-proto v0.0.5-0.20210601070539-28e30ba879cd
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/postgres v1.1.0 // indirect
	gorm.io/gorm v1.21.10 // indirect
)

replace github.com/sktelecom/tks-info => ./
