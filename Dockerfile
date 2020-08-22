#!/bin/sh
FROM golang:alpine
# AS builder
ADD . /go/src/go-social-network
WORKDIR /go/src/go-social-network
# # Git is required for fetching the dependencies.
# RUN ls
RUN apk update && apk add --no-cache git
# RUN go get github.com/dsa0x/go-social-network
COPY . .
RUN go get
RUN go install

# # Build the binary.
# RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .



RUN ls
ENV PORT=8080
EXPOSE 8080
RUN chmod +x /go/src/go-social-network/main
CMD ["/go/src/go-social-network/main"]
ENTRYPOINT ["/go/src/go-social-network/main"]