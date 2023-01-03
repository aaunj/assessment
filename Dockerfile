FROM golang:1.19-alpine

# Set working directory
WORKDIR /go/src/target

# Run tests
CMD CGO_ENABLED=0 go test -v --tags=integration ./... -cover