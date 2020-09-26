// Test 1Fichier filesystem interface
package fichier

import (
	"testing"

	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fs.Config.LogLevel = fs.LogLevelDebug
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestFichier:",
	})
}
