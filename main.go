package main

import (
	_ "embed"
	"go-img/internal/constants"
	"go-img/internal/generator"
	"go-img/internal/img"
	"go-img/internal/search"
	"log"
)

func main() {
	s := search.NewFileSearcher(constants.AssetsDir, constants.IncludedExtensions)
	files, err := s.Search()
	if err != nil {
		log.Fatal("could not search through the given directory")
	}

	images := img.NewFromSlice(files)
	if images == nil {
		log.Fatal("could not find any images")
	}
	images = img.Filter(images, constants.WidthsToConvert)

	generator := generator.NewImageGenerator(images, constants.WidthsToConvert)
	err = generator.GenerateImages()
	if err != nil {
		log.Fatal("could not generate images")
	}
	err = generator.GenerateHTMLs()
	if err != nil {
		log.Fatal("could not generate HTML")
	}
	err = generator.GenerateTempl()
	if err != nil {
		log.Fatal("could not generate templ", err)
	}
}
