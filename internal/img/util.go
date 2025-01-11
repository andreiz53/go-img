package img

import (
	"bytes"
	"image"
	"log"
	"os"
)

func getImageSizes(path string) (width int, height int) {

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(file))
	if err != nil {
		log.Fatal(err)
	}
	return cfg.Width, cfg.Height
}
