package bucket

import (
	"errors"
	"os"
)

type Bucket struct {
	RootPath string
}

var (
	ErrNoLatestFileFound = errors.New("no latest file found in the given query")
)

func (b Bucket) ReadFile(resolved ResolvedQuery) ([]byte, error) {
	content, err := os.ReadFile(resolved.Path)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
