package generator

import (
	"fmt"
	"go-img/internal/img"
	"go-img/internal/util"
	"log"
	"os"
	"strconv"
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
			imageWidth, err := strconv.ParseInt(i.Width, 10, 64)
			util.Check(err)
			resizeWidth, err := strconv.ParseInt(width, 10, 64)
			util.Check(err)
			if resizeWidth < imageWidth {
				_, err := i.Resize(width)
				if err != nil {
					fmt.Println("failed to resize image", err)
					return err
				}
			}
		}
	}
	return nil
}

func (ig ImageGenerator) GenerateHTMLs() error {

	if _, err := os.Stat("tmp/go-images.md"); err == nil {
		os.Remove("tmp/go-images.md")
	}
	file, err := os.OpenFile("tmp/go-images.md", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0661)
	if err != nil {
		log.Fatal("could not generate file with html templates")

	}
	for _, i := range ig.Images {
		i.AlternateWidths = ig.Widths
		node := i.HTML(ig.Widths)

		file.WriteString(fmt.Sprintf("### %s\n", i.Path+"/"+i.Name+i.Extension))
		file.WriteString(node + "\n\n")
	}
	return nil
}

func (ig ImageGenerator) GenerateTempl() error {
	err := util.CopyFile("internal/util/image.templ", "tmp/image.templ")
	if err != nil {
		return err
	}
	return nil
}
