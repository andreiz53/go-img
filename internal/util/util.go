package util

import (
	"path/filepath"
	"slices"
	"strings"
)

func IsValidFileExtension(path string, acceptedExtensions []string) bool {
	fileExtension := filepath.Ext(path)
	return slices.Contains(acceptedExtensions, fileExtension)
}

func RenameImage(path string, suffix string) string {
	dirs := strings.Split(path, "/")
	fileName := dirs[len(dirs)-1]
	fileFragments := strings.Split(fileName, ".")
	fileName = fileFragments[0] + "_" + suffix + "." + fileFragments[1]

	dirs[len(dirs)-1] = fileName

	return strings.Join(dirs, "/")
}
