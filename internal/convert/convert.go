package generator

import (
	"fmt"
	img "go-img/internal/image"
)

type ImageGenerator struct {
	Images []*img.Image
	Widths []string
}

func NewImageGenerator(images []*img.Image, widths []string) ImageGenerator {
	return ImageGenerator{
		Images: images,
		Widths: widths,
	}
}

func (ig ImageGenerator) GenerateImages() error {
	for _, i := range ig.Images {
		for _, width := range ig.Widths {
			_, err := i.Resize(width)
			if err != nil {
				fmt.Println("failed to resize image", err)
			}
		}
	}
	return nil
}

func (ig ImageGenerator) GenerateHTMLs() error {
	return nil
}
