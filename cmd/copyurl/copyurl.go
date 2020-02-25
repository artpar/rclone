package copyurl

import (
	"context"
	"errors"
	"os"

	"github.com/artpar/rclone/cmd"
	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fs/config/flags"
	"github.com/artpar/rclone/fs/operations"
	"github.com/spf13/cobra"
)

var (
	autoFilename = false
	stdout       = false
)

func init() {
	cmd.Root.AddCommand(commandDefinition)
	cmdFlags := commandDefinition.Flags()
	flags.BoolVarP(cmdFlags, &autoFilename, "auto-filename", "a", autoFilename, "Get the file name from the URL and use it for destination file path")
	flags.BoolVarP(cmdFlags, &stdout, "stdout", "", stdout, "Write the output to stdout rather than a file")
}

var commandDefinition = &cobra.Command{
	Use:   "copyurl https://example.com dest:path",
	Short: `Copy url content to dest.`,
	Long: `
Download a URL's content and copy it to the destination without saving
it in temporary storage.

Setting --auto-filename will cause the file name to be retreived from
the from URL (after any redirections) and used in the destination
path.

Setting --stdout or making the output file name "-" will cause the
output to be written to standard output.
`,
	RunE: func(command *cobra.Command, args []string) (err error) {
		cmd.CheckArgs(1, 2, command, args)

		var dstFileName string
		var fsdst fs.Fs
		if !stdout {
			if len(args) < 2 {
				return errors.New("need 2 arguments if not using --stdout")
			}
			if args[1] == "-" {
				stdout = true
			} else if autoFilename {
				fsdst = cmd.NewFsDir(args[1:])
			} else {
				fsdst, dstFileName = cmd.NewFsDstFile(args[1:])
			}
		}
		cmd.Run(true, true, command, func() error {
			if stdout {
				err = operations.CopyURLToWriter(context.Background(), args[0], os.Stdout)
			} else {
				_, err = operations.CopyURL(context.Background(), fsdst, dstFileName, args[0], autoFilename)
			}
			return err
		})
		return nil
	},
}
