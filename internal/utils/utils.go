package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)

	root := filepath.Join(filepath.Dir(b), "../..")
	return root
}

func ZipDirectory(dir string) (io.Reader, error) {
	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// Strip the absolute path up to the current directory, then trim off a leading
		// path separator (for Windows) and replace all instances of Windows path separators
		// with forward slashes as required by the w.Create method.
		f, err := w.Create(strings.Replace(strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator)), "\\", "/", -1))
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
