// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	_ "github.com/artpar/rclone/backend/all" // import all backends
	"github.com/artpar/rclone/cmd"
	_ "github.com/artpar/rclone/cmd/all"    // import all commands
	_ "github.com/artpar/rclone/lib/plugin" // import plugins
)

func main() {
	cmd.Main()
}
