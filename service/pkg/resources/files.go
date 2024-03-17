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

func (f Folder) GetAbsolutePath() string {
	return path.Join(f.ResourcePath, f.Name)
}

func (f Folder) MapFile(name string) File {
	return File{Filename: name, AbsolutePath: path.Join(f.GetAbsolutePath(), name)}
}

func (f Folder) CreateIfNotExist() {
	_ = os.Mkdir(f.GetAbsolutePath(), 0777)
}

type File struct {
	Filename     string
	AbsolutePath string
}

func (r File) Open() ([]byte, error) {
	fullPath := r.AbsolutePath
	return os.ReadFile(fullPath)
}
