#!/bin/sh
FROM golang:alpine AS builder
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

# COPY --from=builder /go/src/go-social-network/main .

# FROM alpine:latest  
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /go/src/go-social-network/main .
# WORKDIR /go/src/go-social-network/


RUN ls
ENV PORT=8080
EXPOSE 8080
ENV WD=/go/src/go-social-network
# RUN chmod +x ./main
# CMD ["/go/src/go-social-network/main"]
ENTRYPOINT ["#!/bin/bash","/go/src/go-social-network"]