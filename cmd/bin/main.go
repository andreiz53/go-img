package main

import (
	"bytes"
	_ "embed"
	"fmt"
	img "go-img/internal/image"
	"go-img/internal/util"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

type Config struct {
	AssetsDir          string
	IncludedExtensions []string
	Sizes              []string
}

func main() {

	config := Config{
		AssetsDir:          "./assets/images",
		IncludedExtensions: []string{".png", ".jpeg", ".jpg"},
		Sizes:              []string{"400", "200"},
	}

	var images []img.Image
	err := filepath.Walk(config.AssetsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && util.IsValidFileExtension(path, config.IncludedExtensions) {
			img := img.Image{
				Name: info.Name(),
				Path: path,
			}
			images = append(images, img)
		}
		return nil
	})
	if err != nil {
		log.Fatal("could not traverse the given directory")
	}
outer:
	for _, i := range images {
		for _, size := range config.Sizes {
			if strings.Contains(i.Name, size) {
				continue outer
			}
		}
		for _, size := range config.Sizes {
			imageWidth, err := strconv.ParseInt(size, 10, 64)
			if err != nil {
				log.Fatal("sizes in config are not integers")
			}

			file, err := os.ReadFile(i.Path)
			if err != nil {
				fmt.Println("err", err)
				break
			}

			img, _, err := image.Decode(bytes.NewReader(file))
			if err != nil {
				fmt.Println("here err", err)
				break
			}

			resizedImage := resize.Resize(uint(imageWidth), 0, img, resize.Lanczos3)
			newPath := util.RenameImage(i.Path, size)

			newFile, err := os.Create(newPath)
			if err != nil {
				log.Fatal("could not create file")
			}
			defer newFile.Close()

			imgExtension := filepath.Ext(i.Name)
			switch imgExtension {
			case ".jpg":
				jpeg.Encode(newFile, resizedImage, nil)
			case ".png":
				png.Encode(newFile, resizedImage)
			}
		}
	}

}
