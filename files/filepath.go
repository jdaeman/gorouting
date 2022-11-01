package files

import (
	"path/filepath"
	"strings"
)

func ToDataPath(inputPath string, newExt string) string {
	dir, file := filepath.Split(inputPath)
	ext := filepath.Ext(file)

	idx := strings.Index(file, ext)

	newFile := file[:idx] + newExt

	return filepath.Join(dir, newFile)
}
