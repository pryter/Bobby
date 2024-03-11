package resources

import (
	"os"
	"path"
)

// Folder serves as a Factory for File
type Folder struct {
	Name         string
	ResourcePath string
}

func (f Folder) MapFile(name string) File {
	return File{Filename: name, AbsolutePath: path.Join(f.ResourcePath, f.Name)}
}

type File struct {
	Filename     string
	AbsolutePath string
}

func (r File) Open() ([]byte, error) {
	fullPath := path.Join(r.AbsolutePath, r.Filename)
	return os.ReadFile(fullPath)
}
