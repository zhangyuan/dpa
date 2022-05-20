package initializer

import (
	"dpa/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func createFiles(fromDirectory string, toDirectory string) error {
	if err := mkdirp(toDirectory); err != nil {
		return err
	}

	entries, e := template.Content.ReadDir(fromDirectory)

	if e != nil {
		return errors.Wrap(e, "fail to read directory")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			newFromDirectory := filepath.Join(fromDirectory, entry.Name())
			newToDirectory := filepath.Join(toDirectory, entry.Name())
			if err := createFiles(newFromDirectory, newToDirectory); err != nil {
				return err
			}
		} else {
			fromFilePath := filepath.Join(fromDirectory, entry.Name())
			input, err := template.Content.ReadFile(fromFilePath)

			if err != nil {
				return errors.Wrapf(err, "fail to read file %s", fromFilePath)
			}
			toFilePath := filepath.Join(toDirectory, entry.Name())
			err = ioutil.WriteFile(toFilePath, input, 0644)
			if err != nil {
				return errors.Wrapf(err, "fail to write file %s", toFilePath)
			}

		}
	}

	return nil
}

func mkdirp(projectDirectory string) error {
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

func Initialize(projectDirectory string) error {
	return createFiles("lakeformation", projectDirectory)
}
