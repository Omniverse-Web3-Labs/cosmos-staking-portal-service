package utils

import (
	"os"
	"path/filepath"
)

var FileUtil = newFileUtil()

type fileUtil struct {
}

func newFileUtil() *fileUtil {
	return &fileUtil{}
}

func (*fileUtil) ExecuteDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
