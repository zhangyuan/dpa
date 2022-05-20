package cmd

import (
	"dpa/pkg/build"
	"dpa/pkg/errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getWorkingDirectory(args)
		if err != nil {
			log.Fatalln("invalid directory", err)
		}

		if err := build.Build(path); err != nil {
			if err, ok := err.(errors.StackTracer); ok {
				for _, f := range err.StackTrace() {
					fmt.Printf("%+s:%d\n", f, f)
				}
			}

			log.Fatalln("fatal error: ", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
