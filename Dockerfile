FROM  golang:1.17.7

WORKDIR /usr/src/app

COPY . .
RUN go mod download
