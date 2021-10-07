//go:build linux || (darwin && amd64)
// +build linux darwin,amd64

package mount2

import (
	"testing"

	"github.com/artpar/rclone/fstest/testy"
	"github.com/artpar/rclone/vfs/vfstest"
)

func TestMount(t *testing.T) {
	testy.SkipUnreliable(t)
	vfstest.RunTests(t, false, mount)
}
