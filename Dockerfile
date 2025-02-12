# setup build image
FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:c25964d301e6c50174d29deadbbaa5ea6443e94b61087b6d89e8f41ef4ebca35 AS build

WORKDIR /app

COPY main.go go.mod go.sum ./
RUN go mod download -x

ARG GO_LINKER_ARGS
ARG TARGETARCH
ARG TARGETOS

COPY pkg ./pkg
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags -trimpath -ldflags="${GO_LINKER_ARGS}" \
    -o ./build/_output/bin/dynatrace-bootstrapper

FROM --platform=$TARGETPLATFORM public.ecr.aws/dynatrace/dynatrace-codemodules:1.299.73.20250109-153420 AS codemodules

# copy bootstrapper binary
COPY --from=build /app/build/_output/bin /opt/dynatrace/oneagent/agent/lib64/

LABEL name="Dynatrace Bootstrapper" \
      vendor="Dynatrace LLC" \
      maintainer="Dynatrace LLC"

ENV USER_UID=1001 \
    USER_NAME=dynatrace-bootstrapper

USER ${USER_UID}:${USER_UID}

ENTRYPOINT ["/opt/dynatrace/oneagent/agent/lib64/dynatrace-bootstrapper"]