// Test AmazonCloudDrive filesystem interface

//go:build acd
// +build acd

package amazonclouddrive_test

import (
	"testing"

	"github.com/artpar/rclone/backend/amazonclouddrive"
	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.NilObject = fs.Object((*amazonclouddrive.Object)(nil))
	fstests.RemoteName = "TestAmazonCloudDrive:"
	fstests.Run(t)
}
