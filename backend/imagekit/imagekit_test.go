package imagekit

import (
	"testing"

	"github.com/artpar/rclone/fstest"
	"github.com/artpar/rclone/fstest/fstests"
)

func TestIntegration(t *testing.T) {
	debug := true
	fstest.Verbose = &debug
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestImageKit:",
		NilObject:       (*Object)(nil),
		SkipFsCheckWrap: true,
	})
}
