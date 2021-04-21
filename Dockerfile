FROM golang:1.16.3-stretch AS builder
LABEL AUTHOR Seungkyu Ahn (seungkyua@gmail.com)

RUN mkdir -p /build
WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/tksinfo ./cmd/server.go

RUN mkdir -p /dist
WORKDIR /dist
RUN cp /build/bin/tksinfo ./tksinfo



FROM golang:alpine3.13

RUN mkdir -p /app
WORKDIR /app

COPY --chown=0:0 --from=builder /dist /app/
EXPOSE 9111

ENTRYPOINT ["/app/tksinfo"]
CMD ["--port", "9111"]
