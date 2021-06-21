FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client
RUN go mod download -x
RUN go build -o httpServer ./cmd/httpserver/main.go

EXPOSE 8080 8080

CMD ["./httpServer"]