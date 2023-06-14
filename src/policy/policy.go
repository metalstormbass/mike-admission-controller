package policy

import (
	"strings"
)

func ValidateContainerTag(image string) (valid bool) {

	imageTagMap := strings.SplitN(image, ":", 2)
	imageTag := string(imageTagMap[1])
	if imageTag != "latest" {
		return true
	} else {
		return false
	}

}
