package lsd

import (
	"os"

	"github.com/artpar/rclone/cmd"
	"github.com/artpar/rclone/fs"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Root.AddCommand(commandDefintion)
}

var commandDefintion = &cobra.Command{
	Use:   "lsd remote:path",
	Short: `List all directories/containers/buckets in the path.`,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(1, 1, command, args)
		fsrc := cmd.NewFsSrc(args)
		cmd.Run(false, false, command, func() error {
			return fs.ListDir(fsrc, os.Stdout)
		})
	},
}
