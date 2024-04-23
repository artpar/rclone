// Package logger implements testing for the sync (and bisync) logger
package logger

import (
	_ "github.com/artpar/rclone/backend/all" // import all backends
	"github.com/artpar/rclone/cmd"
	_ "github.com/artpar/rclone/cmd/all"    // import all commands
	_ "github.com/artpar/rclone/lib/plugin" // import plugins
)

// Main enables the testscript package. See:
// https://bitfieldconsulting.com/golang/cli-testing
// https://pkg.go.dev/github.com/rogpeppe/go-internal@v1.11.0/testscript
func Main() int {
	cmd.Main()
	return 0
}
