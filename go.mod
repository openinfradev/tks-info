module github.com/sktelecom/tks-info

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/sktelecom/tks-contract v0.1.0
	github.com/sktelecom/tks-proto v0.0.4-0.20210419050352-2299e8d5d653
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/sktelecom/tks-info => ./
