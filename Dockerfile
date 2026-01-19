# check=skip=RedundantTargetPlatform
# setup build image
FROM --platform=$BUILDPLATFORM golang:1.25.6@sha256:ce63a16e0f7063787ebb4eb28e72d477b00b4726f79874b3205a965ffd797ab2 AS build

WORKDIR /app

COPY main.go go.mod go.sum ./
RUN go mod download -x

ARG GO_LINKER_ARGS
ARG TARGETARCH
ARG TARGETOS

COPY pkg ./pkg
COPY cmd ./cmd
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags -trimpath -ldflags="${GO_LINKER_ARGS}" \
    -o ./build/_output/bin/dynatrace-bootstrapper

# platform is required, otherwise the copy command will copy the wrong architecture files, don't trust GitHub Actions linting warnings
FROM --platform=$TARGETPLATFORM public.ecr.aws/dynatrace/dynatrace-codemodules:1.311.70.20250416-094918 AS codemodules

# copy bootstrapper binary
COPY --from=build /app/build/_output/bin /opt/dynatrace/oneagent/agent/lib64/

LABEL name="Dynatrace Bootstrapper" \
      vendor="Dynatrace LLC" \
      maintainer="Dynatrace LLC"

ENV USER_UID=1001 \
    USER_NAME=dynatrace-bootstrapper

USER ${USER_UID}:${USER_UID}

ENTRYPOINT ["/opt/dynatrace/oneagent/agent/lib64/dynatrace-bootstrapper"]
