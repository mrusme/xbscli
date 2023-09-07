FROM golang:1.19-alpine

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /xbscli

ENTRYPOINT ["/xbscli", "-f", "-i", "-p", "-s"]
