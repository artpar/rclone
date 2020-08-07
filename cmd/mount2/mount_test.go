// +build linux darwin,amd64

package mount2

import (
	"testing"

	"github.com/artpar/rclone/cmd/mountlib/mounttest"
)

func TestMount(t *testing.T) {
	mounttest.RunTests(t, mount)
}
