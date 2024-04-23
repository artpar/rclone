//go:build darwin && !cmount

package nfsmount

import (
	"testing"

	"github.com/artpar/rclone/vfs/vfscommon"
	"github.com/artpar/rclone/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeMinimal, false, mount)
}
