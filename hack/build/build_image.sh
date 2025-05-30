#!/bin/bash

if [[ ! "${1}" ]]; then
  echo "first param is not set, should be the image without the tag"
  exit 1
fi
if [[ ! "${2}" ]]; then
  echo "second param is not set, should be the tag of the image"
  exit 1
fi

image=${1}
tag=${2}
debug=${3:-false}

commit=$(git rev-parse HEAD)
go_linker_args=$(hack/build/create_go_linker_args.sh "${tag}" "${commit}" "${debug}")
out_image="${image}:${tag}"

if ! command -v docker 2>/dev/null; then
  CONTAINER_CMD=podman
else
  CONTAINER_CMD=docker
fi

BOOTSTRAPPER_BUILD_PLATFORM="--platform=linux/amd64"
if [ -n "${BOOTSTRAPPER_DEV_BUILD_PLATFORM}" ]; then
  echo "overriding platform to ${BOOTSTRAPPER_DEV_BUILD_PLATFORM}"
  BOOTSTRAPPER_BUILD_PLATFORM="--platform=${BOOTSTRAPPER_DEV_BUILD_PLATFORM}"
fi

${CONTAINER_CMD} build "${BOOTSTRAPPER_BUILD_PLATFORM}" . -f ./Dockerfile -t "${out_image}" \
  --build-arg "GO_LINKER_ARGS=${go_linker_args}" \
  --label "quay.expires-after=14d"

