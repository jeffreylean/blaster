FROM golang:1.18 AS builder
LABEL maintainer="weifonglean@airasia.com"

# Copy the directory into the container outside of the gopath
ENV GOPRIVATE gitlab.airasiatech.com
ENV BASEPATH /go/src/gitlab.airasiatech.com/data/platform/blaster

RUN mkdir -p $BASEPATH
WORKDIR $BASEPATH
ADD . $BASEPATH
# netrc file is present in gitlab runner
ADD netrc /root/.netrc

# Download and install any required third party dependencies into the container.
RUN go build -o /go/bin/blast cmd/main.go

# Base image for runtime
FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates procps

WORKDIR /root/
COPY --from=builder /go/bin/blast .
RUN chmod +x ./blast

EXPOSE 8080
CMD ["./blast"]
