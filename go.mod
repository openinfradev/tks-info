module github.com/sktelecom/tks-info

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/jackc/pgx/v4 v4.11.0 // indirect
	github.com/lib/pq v1.10.2
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/sktelecom/tks-contract v0.1.1-0.20210604023929-73ffc015c1f1
	github.com/sktelecom/tks-proto v0.0.6-0.20210622012523-ded9f951101f // indirect
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/sktelecom/tks-info => ./
