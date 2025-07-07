package container

type ImageInfo struct {
	Registry    string `json:"container_image.registry,omitempty"`
	Repository  string `json:"container_image.repository,omitempty"`
	Tag         string `json:"container_image.tags,omitempty"`
	ImageDigest string `json:"container_image.digest,omitempty"`
}

func (img ImageInfo) ToURI() string {
	var imageName string

	registry := img.Registry
	repository := img.Repository
	tag := img.Tag
	digest := img.ImageDigest

	if registry != "" {
		imageName = registry + "/" + repository
	} else {
		imageName = repository
	}

	if tag != "" {
		imageName += ":" + tag
	}

	if digest != "" {
		imageName += "@" + digest
	}

	return imageName
}
