# check=skip=RedundantTargetPlatform
# setup build image
FROM --platform=$BUILDPLATFORM golang:1.24.2@sha256:1ecc479bc712a6bdb56df3e346e33edcc141f469f82840bab9f4bc2bc41bf91d AS build

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
# FROM --platform=$TARGETPLATFORM public.ecr.aws/dynatrace/dynatrace-codemodules:1.307.57.20250217-152612 AS codemodules
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
