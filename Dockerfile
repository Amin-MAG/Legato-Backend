FROM  golang:1.17.7

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
