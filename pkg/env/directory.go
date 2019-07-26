package env

import (
	"os"
	"path/filepath"
)

func GetExecutableDir() (string, error)  {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(executablePath), nil
}
