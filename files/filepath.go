package files

import (
	"path/filepath"
	"strings"
)

const (
	GEONODE     = ".node"
	NBGEDGE     = ".edge"
	GEOMETRY    = ".geo"
	ANNOTATION  = ".anno"
	RESTRICTION = ".restriction"
)

func ToDataPath(inputPath string, newExt string) string {
	dir, file := filepath.Split(inputPath)
	ext := filepath.Ext(file)

	idx := strings.Index(file, ext)

	newFile := file[:idx] + newExt

	return filepath.Join(dir, newFile)
}
