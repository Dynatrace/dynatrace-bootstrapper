IMAGE ?= quay.io/dynatrace/dynatrace-bootstrapper
DEBUG ?= false
BOOTSTRAPPER_BUILD_ARCHS ?= amd64,arm64

#Needed for the e2e pipeline to work
BRANCH ?= $(shell git branch --show-current)
SNAPSHOT_SUFFIX ?= $(shell echo "${BRANCH}" | sed "s/[^a-zA-Z0-9_-]/-/g")
ifneq ($(BRANCH), main)
	TAG ?= snapshot-${SNAPSHOT_SUFFIX}
else
	TAG ?= snapshot
endif

#use the digest if digest is set
ifeq ($(DIGEST),)
	IMAGE_URI ?= "$(IMAGE):$(TAG)"
else
	IMAGE_URI ?= "$(IMAGE):$(TAG)@$(DIGEST)"
endif



ensure-tag-not-snapshot:
ifeq ($(TAG), snapshot)
	$(error "Image tag is snapshot, please set TAG to a valid tag")
endif

## Builds an Bootstrapper image with a given IMAGE and TAG
images/build: ensure-tag-not-snapshot
	./hack/build/build_image.sh "${IMAGE}" "${TAG}" "${DEBUG}" "${BOOTSTRAPPER_BUILD_ARCHS}"

## Pushes an ALREADY BUILT Bootstrapper image with a given IMAGE and TAG
images/push: ensure-tag-not-snapshot
	./hack/build/push_image.sh "${IMAGE}" "${TAG}" "${BOOTSTRAPPER_BUILD_ARCHS}"

## Builds an Bootstrapper image and pushes it
images/build/push: images/build images/push

