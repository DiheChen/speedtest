FROM golang:alpine as builder

WORKDIR /app
COPY main.go /app/

RUN apk add --no-cache xz curl && \
    go mod init speedtest && \
    go mod tidy && \
    wget $(curl -s https://api.github.com/repos/upx/upx/releases/latest | grep 'tag_name' | cut -d\" -f4 | xargs -I {} sh -c 'echo "https://github.com/upx/upx/releases/download/{}/upx-${1#v}-amd64_linux.tar.xz"' _ {}) -O upx.tar.xz && \
    tar -xvf upx.tar.xz && \
    mv $(find . -type f -name 'upx' | head -n 1) ./upx && \
    chmod +x ./upx && \
    go build -ldflags "-s -w" && \
    ./upx speedtest

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/speedtest /app/

ENTRYPOINT ["/app/speedtest"]
