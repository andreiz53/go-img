package util

import (
	"io"
	"os"
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

func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}
