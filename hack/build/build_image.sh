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
arch=${4:-"amd64,arm64"}

IFS=',' read -ra supported_architectures <<< "$arch"
echo "Creating image-index manifest for ${supported_architectures[*]}"

commit=$(git rev-parse HEAD)
go_linker_args=$(hack/build/create_go_linker_args.sh "${tag}" "${commit}" "${debug}")
out_image="${image}:${tag}"

if ! command -v docker 2>/dev/null; then
  CONTAINER_CMD=podman
else
  CONTAINER_CMD=docker
fi

for architecture in "${supported_architectures[@]}"; do
  ${CONTAINER_CMD} build "--platform=linux/${architecture}" . -f ./Dockerfile -t "${out_image}-${architecture}" \
    --build-arg "GO_LINKER_ARGS=${go_linker_args}" \
    --label "quay.expires-after=14d"
done
