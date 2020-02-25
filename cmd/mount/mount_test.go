// +build linux,go1.11 darwin,go1.11 freebsd,go1.11

package mount

import (
	"runtime"
	"testing"

	"github.com/artpar/rclone/cmd/mountlib/mounttest"
)

func TestMount(t *testing.T) {
	if runtime.NumCPU() <= 2 {
		t.Skip("FIXME skipping mount tests as they lock up on <= 2 CPUs - See: https://github.com/artpar/rclone/issues/3154")
	}
	mounttest.RunTests(t, mount)
}
