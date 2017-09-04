package genautocomplete

import (
	"log"
	"github.com/artpar/rclone/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Root.AddCommand(completionDefinition)
}

var completionDefinition = &cobra.Command{
	Use:   "genautocomplete [shell]",
	Short: `Output completion script for a given shell.`,
	Long: `
Generates a shell completion script for rclone.
Run with --help to list the supported shells.
`,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(0, 1, command, args)
		out := "/etc/bash_completion.d/rclone"
		if len(args) > 0 {
			out = args[0]
		}
		err := cmd.Root.GenBashCompletionFile(out)
		if err != nil {
			log.Print(err)
		}
	},
}
