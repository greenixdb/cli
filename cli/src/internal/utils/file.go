package utils

import (
	"io"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func GetCurrentDir() (string, error) {
	return os.Getwd()
}

func GetExecutablePath() (string, error) {
	return os.Executable()
}

