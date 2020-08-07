// Test Union filesystem interface
package union_test

import (
	"testing"

	_ "github.com/artpar/rclone/backend/local"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName:  "TestUnion:",
		NilObject:   nil,
		SkipFsMatch: true,
	})
}
