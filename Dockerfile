FROM golang:1.19.4-alpine3.17 AS GO_BUILD
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
COPY ./expense ./expense
RUN go build -o /app/kbtg .

FROM alpine:3.17
WORKDIR /app
COPY --from=GO_BUILD /app/kbtg .

CMD sleep 10; /app/kbtg