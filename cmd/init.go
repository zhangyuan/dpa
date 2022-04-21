package cmd

import (
	"fmt"
	"log"

	"dp/pkg/errors"
	"dp/pkg/initializer"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getWorkingDirectory(args)
		if err != nil {
			log.Fatalln("invalid directory", err)
		}

		if err := initializer.Initialize(path); err != nil {
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
	rootCmd.AddCommand(initCmd)
}
