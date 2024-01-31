// Test smb filesystem interface
package smb_test

import (
	"testing"

	"github.com/artpar/rclone/backend/smb"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSMB:rclone",
		NilObject:  (*smb.Object)(nil),
	})
}
