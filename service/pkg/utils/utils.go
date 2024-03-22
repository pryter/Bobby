package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	return as, nil
}

// GetProjectRoot gives the project root.
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)

	root := filepath.Join(filepath.Dir(b), "../../..")
	return root
}

// ZipDirectory recursively compress a given directory to file io.
// Snippet from https://stackoverflow.com/a/69225665 by Clark McCauley
func ZipDirectory(dir string) (io.Reader, error) {
	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)
	err := filepath.Walk(
		dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			// Strip the absolute path up to the current directory, then trim off a leading
			// path separator (for Windows) and replace all instances of Windows path separators
			// with forward slashes as required by the w.Create method.
			f, err := w.Create(
				strings.Replace(
					strings.TrimPrefix(
						strings.TrimPrefix(
							path, dir,
						), string(filepath.Separator),
					), "\\", "/", -1,
				),
			)
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
		},
	)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
