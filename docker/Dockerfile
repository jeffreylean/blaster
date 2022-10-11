FROM golang:1.18 AS builder
LABEL maintainer="leanwf1117@gmail.com"

# Copy the directory into the container outside of the gopath
ENV BASEPATH /go/src/blaster

RUN mkdir -p $BASEPATH
WORKDIR $BASEPATH
ADD . $BASEPATH

RUN go build -o /go/bin/blast main.go

# Base image for runtime
FROM debian:latest
# Download and install any required third party dependencies into the container.
RUN apt-get update && apt-get install -y ca-certificates procps

WORKDIR /root/
COPY --from=builder /go/bin/blast .
RUN chmod +x ./blast

EXPOSE 1323
CMD [ "./blast","server","1323"]
