# FROM golang:1.23.4-alpine3.21 AS builder
FROM registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.docker.io.library.golang.1.23.4-alpine3.21 AS builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /workspace/brotato-unlocker

# FROM gcr.io/distroless/static-debian12
FROM registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.gcr.io.distroless.static-debian12.latest

WORKDIR /workspace
COPY --from=builder /workspace/brotato-unlocker /brotato-unlocker
ENTRYPOINT ["/brotato-unlocker"]