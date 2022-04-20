package cmd

import (
	"fmt"
	"log"
	"os"

	"dp/commands/initializer"
	"dp/pkg"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := findDirectory(args)
		if err != nil {
			log.Fatalln("invalid directory", err)
		}

		if err := initializer.Initialize(path); err != nil {
			if err, ok := err.(pkg.StackTracer); ok {
				for _, f := range err.StackTrace() {
					fmt.Printf("%+s:%d\n", f, f)
				}
			}

			log.Fatalln("fatal error: ", err.Error())
		}
	},
}

func findDirectory(args []string) (string, error) {
	if len(args) == 0 {
		path, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "fail to get working directory")
		}

		return path, nil
	}

	return args[0], nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
