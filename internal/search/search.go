package search

import (
	"go-img/internal/util"
	"io/fs"
	"path/filepath"
)

type Searcher interface {
	Search() ([]string, error)
}

type FileSearcher struct {
	AssetsDir          string
	IncludedExtensions []string
}

func NewFileSearcher(dir string, extensions []string) Searcher {
	return FileSearcher{
		AssetsDir:          dir,
		IncludedExtensions: extensions,
	}
}

func (sc FileSearcher) Search() ([]string, error) {
	var output []string
	err := filepath.Walk(sc.AssetsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && util.IsValidFileExtension(path, sc.IncludedExtensions) {
			output = append(output, path)
		}
		return nil
	})
	if err != nil {
		return output, err
	}
	return output, nil
}
