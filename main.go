package main

import (
	_ "embed"
	"go-img/internal/constants"
	"go-img/internal/generator"
	"go-img/internal/img"
	"go-img/internal/search"
	"go-img/internal/util"
	"log"
)

func main() {
	s := search.NewFileSearcher(constants.AssetsDir, constants.IncludedExtensions)
	files, err := s.Search()
	util.Check(err)

	images := img.NewFromSlice(files)
	if images == nil {
		log.Fatal("could not find any images")
	}
	images = img.Filter(images, constants.WidthsToConvert)

	generator := generator.NewImageGenerator(images, constants.WidthsToConvert)
	err = generator.GenerateImages()
	util.Check(err)

	err = generator.GenerateHTMLs()
	util.Check(err)

	err = generator.GenerateTempl()
	util.Check(err)
}
