// Test Linkbox filesystem interface
package linkbox_test

import (
	"testing"

	"github.com/artpar/rclone/backend/linkbox"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestLinkbox:",
		NilObject:  (*linkbox.Object)(nil),
	})
}
