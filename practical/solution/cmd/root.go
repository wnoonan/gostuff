package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wnoonan/gostuff/practical/solution/pkg/jobbr"
)

var rootCmd = &cobra.Command{
	Use:   "jobbr",
	Short: "Jobs and Job accessories",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := jobbr.ParseEnv()
		if err != nil {
			log.Fatal(err)
		}

		jobber := jobbr.NewJobbr(env)
		out := jobbr.NewOutput(jobber)

		return out.RootOpts()
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
