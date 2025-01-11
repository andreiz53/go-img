package main

import (
	_ "embed"
	"go-img/internal/generator"
	img "go-img/internal/image"
	"go-img/internal/search"
	"log"
)

var (
	AssetsDir          = "assets/images"
	FSDir              = "assets/images"
	TmpDir             = "tmp"
	IncludedExtensions = []string{".png", ".jpg", ".jpeg"}
	WidthsToConvert    = []string{"200", "600"}
)

func main() {
	s := search.NewFileSearcher(AssetsDir, IncludedExtensions)
	files, err := s.Search()
	if err != nil {
		log.Fatal("could not search through the given directory")
	}

	images := img.NewFromSlice(files)
	if images == nil {
		log.Fatal("could not find any images")
	}
	images = img.Filter(images, WidthsToConvert)

	generator := generator.NewImageGenerator(images, WidthsToConvert)
	err = generator.GenerateImages()
	if err != nil {
		log.Fatal("could not generate images")
	}
	err = generator.GenerateHTMLs()
	if err != nil {
		log.Fatal("could not generate HTML")
	}
}
