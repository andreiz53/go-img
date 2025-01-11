package img

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Image struct {
	Name            string
	Path            string
	Extension       string
	Width           string
	Height          string
	AlternateWidths []string
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

	// calculate width and height
	width, height := getImageSizes(path)

	return &Image{
		Extension: extension,
		Name:      name,
		Path:      newPath,
		Width:     fmt.Sprint(width),
		Height:    fmt.Sprint(height),
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
	return *New(newPath), nil
}

func (i Image) HTML(widths []string) string {
	attributes := []html.Attribute{
		{Key: "width", Val: i.Width},
		{Key: "height", Val: i.Height},
		{Key: "alt", Val: ""},
		{Key: "title", Val: i.Name},
		i.HTMLSizes(),
		i.HTMLSrcset(),
	}

	node := html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Img,
		Data:     "img",
		Attr:     attributes,
	}

	var b bytes.Buffer
	err := html.Render(&b, &node)
	if err != nil {
		log.Fatal("could not read image html")
	}
	return b.String()
}

func (i Image) HTMLSizes() html.Attribute {
	if len(i.AlternateWidths) == 0 {
		return html.Attribute{}
	}
	attr := html.Attribute{
		Key: "sizes",
	}
	attrVal := fmt.Sprintf("(min-width: 0px) and (max-width: %spx) %spx, ", i.AlternateWidths[0], i.AlternateWidths[0])
	if len(i.AlternateWidths) == 1 {
		attr.Val = attrVal + i.Width + "px"
		return attr
	}
	for idx, width := range i.AlternateWidths {
		if idx > 0 {
			prevWidth := i.AlternateWidths[idx-1]
			attrVal += fmt.Sprintf("(min-width: %spx) and (max-width: %spx) %spx, ", prevWidth, width, width)
		}
	}
	attr.Val = attrVal + i.Width + "px"
	return attr
}

func (i Image) HTMLSrcset() html.Attribute {
	if len(i.AlternateWidths) == 0 {
		return html.Attribute{}
	}
	attr := html.Attribute{
		Key: "srcset",
	}
	var attrVal string
	for idx, width := range i.AlternateWidths {
		altImgSrc := i.Path + "/" + i.Name + "_" + width + i.Extension
		attrVal += altImgSrc + " " + width + "w"
		if idx < len(i.AlternateWidths)-1 {
			attrVal += ", "
		}
	}
	attr.Val = attrVal
	return attr
}
