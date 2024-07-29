### Builder
FROM golang:1.22.5 AS builder

RUN mkdir -p /app /src/src
WORKDIR /gopath_dir/src/tgapi
ENV GOPATH /gopath_dir
ENV GOBIN $GOPATH/bin

COPY go.sum $GOPATH/src/tgapi
COPY go.mod $GOPATH/src/tgapi
RUN go mod download
COPY . $GOPATH/src/tgapi

### Go opts
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go get -u ./... &&\ 
    go build -o $BUILD_BIN_PATH ./cmd/app &&\
    rm -rf /src

CMD $BIN_PATH
