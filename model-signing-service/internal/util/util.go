package util

import (
	"archive/zip"
	"io"
	"strings"
	"os"
	"path/filepath"
)

func FilterModeArgs(args []string) []string {
	var filtered []string
	skip := false
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		if arg == "--mode" {
			skip = true
			continue
		}
		if strings.HasPrefix(arg, "--mode=") {
			continue
		}
		filtered = append(filtered, arg)
		_ = i
	}
	return filtered
}

func UnzipTemp(src string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "bundle_extract")
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(tmpDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return "", err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(outFile, rc); err != nil {
			return "", err
		}
		outFile.Close()
		rc.Close()
	}
	return tmpDir, nil
}
