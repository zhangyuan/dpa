package helpers

import (
	"os"

	"github.com/pkg/errors"
)

func Mkdirp(projectDirectory string) error {
	if _, err := os.Stat(projectDirectory); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(projectDirectory, 0755); err != nil {
				return errors.Wrap(err, "fail to mkdir all")
			}
		} else {
			return errors.Wrap(err, "fail to stat path")
		}
	}
	return nil
}
