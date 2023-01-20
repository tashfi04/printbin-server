package utils

import (
	"github.com/tashfi04/printbin-server/config"
	"os"
	"path"
)

func BuildPath(keywords ...string) string {

	dest := append([]string{config.App().StoragePath}, keywords...)

	return path.Join(dest...)
}

func CreateDirectory(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
