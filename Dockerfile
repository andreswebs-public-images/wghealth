# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"
ARG GOPROXY="https://proxy.golang.org"

ENV GOARCH="${TARGETARCH}"
ENV GOOS="${TARGETOS}"
ENV GOPROXY="${GOPROXY}"
ENV OUT_DIR="/root-layer/etc/s6-overlay/s6-rc.d/init-mod-wireguard-wghealth-install/bin"
ENV APP_NAME="wghealth"

WORKDIR /src
COPY go.mod .
COPY main.go .
COPY root-layer/ /root-layer/
RUN mkdir --parents "${OUT_DIR}"
RUN CGO_ENABLED=0 go build -a -o "${OUT_DIR}/${APP_NAME}" .

FROM scratch
LABEL maintainer=andreswebs@pm.me
COPY --from=build /root-layer/ /
