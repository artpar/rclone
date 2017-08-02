package config

import (
	"github.com/artpar/rclone/cmd"
	"github.com/artpar/rclone/fs"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Root.AddCommand(commandDefintion)
}

var commandDefintion = &cobra.Command{
	Use:   "config",
	Short: `Enter an interactive configuration session.`,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(0, 0, command, args)
		fs.EditConfig()
	},
}
