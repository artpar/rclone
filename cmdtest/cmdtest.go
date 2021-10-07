// Package cmdtest creates a testable interface to rclone main
//
// The interface is used to perform end-to-end test of
// commands, flags, environment variables etc.
//
package cmdtest

// The rest of this file is a 1:1 copy from rclone.go

import (
	_ "github.com/artpar/rclone/backend/all" // import all backends
	"github.com/artpar/rclone/cmd"
	_ "github.com/artpar/rclone/cmd/all"    // import all commands
	_ "github.com/artpar/rclone/lib/plugin" // import plugins
)

func main() {
	cmd.Main()
}
