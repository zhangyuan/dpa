package cmd

import (
	"dp/pkg/compile"
	"dp/pkg/errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// compileCmd represents the build command
var compileCmd = &cobra.Command{
	Use: "compile",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getWorkingDirectory(args)
		if err != nil {
			log.Fatalln("invalid directory", err)
		}

		if err := compile.Compile(path); err != nil {
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
	rootCmd.AddCommand(compileCmd)
}
