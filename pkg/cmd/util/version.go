package util

import (
	"fmt"
	"strings"
)

func TrimVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}

func GetReleaseBranchName(version string) string {
	return fmt.Sprintf("release-%s", version)
}

func GetReleaseTagName(version string) string {
	return fmt.Sprintf("v%s", version)
}
