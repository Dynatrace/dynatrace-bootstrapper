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
arch=${3:-"amd64,arm64"}

out_image="${image}:${tag}"


IFS=',' read -ra supported_architectures <<< "$arch"
images=()
echo "Creating image-index manifest for ${supported_architectures[*]}"

if ! command -v docker 2>/dev/null; then
  CONTAINER_CMD=podman
else
  CONTAINER_CMD=docker
fi

for architecture in "${supported_architectures[@]}"; do
    ${CONTAINER_CMD} push "${out_image}-${architecture}"
    images+=("${out_image}-${architecture}")
done

${CONTAINER_CMD} manifest rm "${out_image}" 2>/dev/null || true
${CONTAINER_CMD} manifest create "${out_image}" "${images[@]}"

sha256=$(${CONTAINER_CMD} manifest push "${out_image}")
echo "Image index created locally with digest ${sha256}"
