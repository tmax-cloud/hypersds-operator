# Build the manager binary
FROM golang:1.13 as builder

ARG CEPH_VERSION=15.2.8
RUN apt-get update && apt-get install -y \
    wget software-properties-common lsb-release \
    sshpass 
RUN wget https://download.ceph.com/keys/release.asc &&\
apt-key add release.asc &&\
rm release.asc 
RUN add-apt-repository \
"deb [arch=amd64] https://download.ceph.com/debian-${CEPH_VERSION}/ \
$(lsb_release -cs) main" &&\
apt-get update &&\
apt-get install -y ceph-common
RUN apt-get purge -y wget software-properties-common lsb-release &&\
apt-get autoremove -y &&\
rm -rf /var/lib/apt/lists/*

RUN mkdir -p /working/config/

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
