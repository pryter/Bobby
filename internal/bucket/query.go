package bucket

import (
	"errors"
	"os"
	"path"
	"sort"
	"strings"
)

type DataContainerType string

var (
	ArtifactDataContainer   DataContainerType = "artifacts"
	RepositoryDataContainer DataContainerType = "repo"
)
var (
	ErrInvalidContainerType = errors.New("invalid container type")
	ErrInvalidQueryString   = errors.New("invalid query string \n ex. bucket/{repoID}/{type}/{query}")
)

type Query struct {
	ID            string
	ContainerType DataContainerType
	QueryStr      string
}

type ResolvedQuery struct {
	Path     string
	Filename string
}

func (q Query) getLatest(bucket Bucket) (string, error) {
	dirPath := path.Join(bucket.RootPath, q.ID, string(q.ContainerType))
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", ErrInvalidQueryString
	}
	sort.Slice(
		files, func(i, j int) bool {
			return files[i].Name() > files[j].Name()
		},
	)

	if len(files) <= 0 {
		return "", ErrNoLatestFileFound
	}

	latestFile := path.Join(dirPath, files[0].Name())
	return latestFile, nil
}

func (q Query) Resolve(bucket Bucket) (ResolvedQuery, error) {

	var resolved ResolvedQuery
	var err error

	switch q.QueryStr {
	case "latest":
		var rPath string
		rPath, err = q.getLatest(bucket)
		resolved = ResolvedQuery{Path: rPath, Filename: path.Base(rPath)}
		break
	default:
		resolved = ResolvedQuery{
			Path: path.Join(
				bucket.RootPath, q.ID, string(q.ContainerType), q.QueryStr,
			),
			Filename: q.QueryStr,
		}
	}

	return resolved, err
}

func NewQuery(queryStr string) (*Query, error) {
	splitFn := func(c rune) bool {
		return c == '/'
	}

	data := strings.FieldsFunc(queryStr, splitFn)

	if len(data) < 4 {
		return &Query{}, ErrInvalidQueryString
	}

	var containerType DataContainerType

	switch data[2] {
	case "artifacts":
		containerType = ArtifactDataContainer
		break
	case "repo":
		containerType = RepositoryDataContainer
		break
	default:
		return &Query{}, ErrInvalidContainerType
	}

	return &Query{ID: data[1], ContainerType: containerType, QueryStr: data[3]}, nil
}
