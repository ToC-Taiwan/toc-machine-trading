package log

import (
	"os"
	"path/filepath"
)

func getExecPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Clean(filepath.Dir(ex))
}
