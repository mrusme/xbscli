FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY lz77/*.go ./lz77/
COPY *.go ./

RUN go build -o /xbscli-docker

ENTRYPOINT ["/xbscli-docker", "-f", "-i", "-p", "-s"]