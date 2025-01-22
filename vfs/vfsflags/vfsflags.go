// Package vfsflags implements command line flags to set up a vfs
package vfsflags

import (
	"github.com/artpar/rclone/fs/config/flags"
	"github.com/artpar/rclone/vfs/vfscommon"
	"github.com/spf13/pflag"
)

// AddFlags adds the non filing system specific flags to the command
func AddFlags(flagSet *pflag.FlagSet) {
	flags.AddFlagsFromOptions(flagSet, "", vfscommon.OptionsInfo)
}
