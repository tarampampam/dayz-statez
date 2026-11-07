package mod

import (
	"io"
	"os"
	"path/filepath"
)

func ReadFile(dirPath string, fileName FileName) (io.ReadCloser, error) {
	return os.Open(filepath.Join(dirPath, string(fileName)))
}
