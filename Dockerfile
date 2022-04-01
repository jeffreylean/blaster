FROM golang:1.18 AS builder
LABEL maintainer="weifonglean@airsia.com"

# Copy the directory into the container outside of the gopath
RUN mkdir -p /go/src/gitlab.airasiatech.com/data-coe/data-engineering/blaster.git
WORKDIR /go/src/gitlab.airasiatech.com/data-coe/data-engineering/blaster/
ADD . /go/src/gitlab.airasiatech.com/data-coe/data-engineering/blaster/

# Download and install any required third party dependencies into the container.
RUN go build -o /go/bin/blaster ./cmd

# Base image for runtime
FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /root/
COPY --from=builder /go/bin/blaster .
RUN chmod +x ./blaster

EXPOSE 8081
CMD ["./blaster"]
