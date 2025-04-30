# Print help text
help:
    @just --list --unsorted

# Building dynamically linked binary
build:
    go build -C ./src -o ../todayiwork

# Building statically linked binary
build-static os="linux" arch="amd64":
    CGO_ENABLED=0 GOOS={{ os }} GOARCH={{ arch }} go build -C ./src  -o ../todayiwork
