package cmd

import (
	"dp/commands/initializer"
	"log"
	"os"

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
			log.Fatalln("fail to init", err)
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
