# FROM golang:1.23.4-alpine3.21 AS builder
FROM registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.docker.io.library.golang.1.23.4-alpine3.21 AS builder

WORKDIR /workspace

ENTRYPOINT ["./scripts/build.sh"]