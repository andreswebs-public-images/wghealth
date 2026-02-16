# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /build
COPY src/go.mod .
# COPY src/go.sum .
RUN go mod download
COPY src/ ./

ENV OUT_DIR="/root-layer/etc/s6-overlay/s6-rc.d/init-mod-wireguard-wghealth-install/bin"
ENV APP_NAME="wghealth"

COPY root-layer/ /root-layer/
RUN mkdir --parents "${OUT_DIR}"

ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 go build -a -o "${OUT_DIR}/${APP_NAME}" .

FROM scratch
LABEL maintainer=andreswebs@pm.me
COPY --from=build /root-layer/ /
