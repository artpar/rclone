// Test AzureBlob filesystem interface

// +build !freebsd,!netbsd,!openbsd,!plan9,!solaris,go1.8

package azureblob_test

import (
	"testing"

	"github.com/artpar/rclone/backend/azureblob"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestAzureBlob:",
		NilObject:  (*azureblob.Object)(nil),
	})
}
