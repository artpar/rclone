//go:build !plan9
// +build !plan9

// Test Tardigrade filesystem interface
package tardigrade_test

import (
	"testing"

	"github.com/artpar/rclone/backend/tardigrade"
	"github.com/artpar/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestTardigrade:",
		NilObject:  (*tardigrade.Object)(nil),
	})
}
