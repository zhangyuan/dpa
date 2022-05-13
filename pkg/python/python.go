package python

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type PythonModule struct {
	Name    string
	Version string
}

type PythonModules []PythonModule

func GetPythonRequirements(projectDirectory string) (PythonModules, error) {
	if projectDirectory == "" {
		projectDirectory = "."
	}

	modules := []PythonModule{}

	requirementsPath := path.Join(projectDirectory, "requirements.txt")
	if FileExists(requirementsPath) {
		bytes, err := ioutil.ReadFile(requirementsPath)
		if err != nil {
			return nil, errors.Wrap(err, "fail to read requirements.txt")
		}
		content := string(bytes)

		for _, line := range strings.Split(content, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			nameVersionPair := strings.Split(line, "==")

			var pythonModule PythonModule
			if len(nameVersionPair) == 2 {
				pythonModule = PythonModule{Name: nameVersionPair[0], Version: nameVersionPair[1]}
			} else {
				pythonModule = PythonModule{Name: nameVersionPair[0]}
			}

			modules = append(modules, pythonModule)
		}
	}

	return modules, nil
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (pythonModule *PythonModule) ToString() string {
	if pythonModule.Version == "" {
		return pythonModule.Name
	}

	return fmt.Sprintf("%s==%s", pythonModule.Name, pythonModule.Version)
}

func (pythonModules *PythonModules) ToString() string {
	return strings.Join(
		lo.Map(*pythonModules, func(t PythonModule, i int) string {
			return t.ToString()
		}), ",")
}
