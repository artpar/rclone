// Test Yandex filesystem interface
package yandex_test

import (
	"testing"

	"github.com/artpar/rclone/backend/yandex"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestYandex:",
		NilObject:  (*yandex.Object)(nil),
	})
}
