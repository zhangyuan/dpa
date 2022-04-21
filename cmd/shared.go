package cmd

import (
	"os"

	"github.com/pkg/errors"
)

func getWorkingDirectory(args []string) (string, error) {
	if len(args) == 0 {
		path, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "fail to get working directory")
		}

		return path, nil
	}

	return args[0], nil
}
