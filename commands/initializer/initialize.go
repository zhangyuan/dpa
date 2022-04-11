package initializer

import (
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

func Initialize(directoryPath string) error {
	if _, err := os.Stat(directoryPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(directoryPath, 0755); err != nil {
				return errors.Wrap(err, "fail to mkdir all")
			}
		} else {
			return errors.Wrap(err, "fail to stat path")
		}
	}

	_ = pkger.Include("templates/")

	var makeFullPath = func(s string) string {
		return strings.TrimSuffix(directoryPath, "/") + "/" + strings.TrimPrefix(s, "/")
	}
	err := pkger.Walk("/templates", func(fullPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		templatePath := fullPath[len("dp:/templates"):]

		targetPath := makeFullPath(templatePath)

		if info.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return errors.Wrapf(err, "fail to make dir %s", targetPath)
			}
		} else {
			templateFile, err := pkger.Open(fullPath)
			if err != nil {
				return errors.Wrap(err, "fail to open template")
			}

			defer templateFile.Close()

			buf := make([]byte, 10)

			targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				return errors.Wrapf(err, "fail to open file %s", targetPath)
			}

			defer targetFile.Close()

			for {
				n, err := templateFile.Read(buf)

				if err != nil && err != io.EOF {
					return errors.Wrapf(err, "fail to read from %s", templateFile.Path())
				}
				if n == 0 {
					break
				}

				if _, err := targetFile.Write(buf[:n]); err != nil {
					return errors.Wrapf(err, "fail to write to %s", targetFile.Name())
				}
			}
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "fail to walk")
	}

	return nil
}
