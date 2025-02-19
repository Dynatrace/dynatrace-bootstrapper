package container

type ImageInfo struct {
	Registry    string `json:"container_image.registry"`
	Repository  string `json:"container_image.repository"`
	Tag         string `json:"container_image.tags"`
	ImageDigest string `json:"container_image.digest"`
}

func (img ImageInfo) ToURI() string {
	var imageName string

	registry := img.Registry
	repository := img.Repository
	tag := img.Tag
	digest := img.ImageDigest

	if registry != "" {
		if repository != "" {
			imageName = registry + "/" + repository
		} else {
			imageName = repository
		}

		if digest != "" {
			imageName += "@" + digest
		} else if tag != "" {
			imageName += ":" + tag
		}
	}

	return imageName
}
