package img

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

type Image struct {
	Name      string
	Path      string
	Extension string
}

const (
	ExtensionJPG = ".jpg"
	ExtensionPNG = ".png"
)

func New(path string) *Image {
	if path == "" {
		return nil
	}
	splitPath := strings.Split(path, "/")

	name := splitPath[len(splitPath)-1]
	extension := filepath.Ext(name)
	if extension == "" {
		return nil
	}

	name = name[:len(name)-len(extension)]
	newPath := strings.Join(splitPath[:len(splitPath)-1], "/")

	return &Image{
		Extension: extension,
		Name:      name,
		Path:      newPath,
	}
}

func NewFromSlice(paths []string) []*Image {
	if len(paths) == 0 {
		return nil
	}
	images := []*Image{}
	for _, path := range paths {
		image := New(path)
		if image == nil {
			continue
		}
		images = append(images, image)
	}
	return images
}

func (i Image) FullPath() string {
	return i.Path + "/" + i.Name + i.Extension
}

func Filter(images []*Image, suffixes []string) []*Image {
	filteredImages := []*Image{}
	for _, image := range images {
		if !image.hasSuffix(suffixes) {
			filteredImages = append(filteredImages, image)
		}
	}
	return filteredImages
}

func (i Image) hasSuffix(suffixes []string) bool {
	for _, suffix := range suffixes {
		imageSuffix := i.Name[len(i.Name)-len(suffix):]
		if suffix == imageSuffix {
			return true
		}
	}
	return false
}

func (i Image) Resize(width string) (Image, error) {
	newWidth, err := strconv.ParseInt(width, 10, 64)
	if err != nil {
		return Image{}, err
	}
	initialFile, err := os.ReadFile(i.FullPath())
	if err != nil {
		return Image{}, err
	}

	initialImage, _, err := image.Decode(bytes.NewReader(initialFile))
	if err != nil {
		return Image{}, err
	}

	newImage := resize.Resize(uint(newWidth), 0, initialImage, resize.Lanczos3)
	newName := i.Name + "_" + width
	newPath := i.Path + "/" + newName + i.Extension
	newFile, err := os.Create(newPath)

	if err != nil {
		return Image{}, err
	}
	defer newFile.Close()

	switch i.Extension {
	case ExtensionJPG:
		jpeg.Encode(newFile, newImage, nil)
	case ExtensionPNG:
		png.Encode(newFile, newImage)
	default:
		return Image{}, errors.New("unsupported file extension")
	}
	return Image{
		Name:      newName,
		Path:      newPath,
		Extension: i.Extension,
	}, nil
}