//go:build noselfupdate

package selfupdate

import (
	"github.com/artpar/rclone/lib/buildinfo"
)

func init() {
	buildinfo.Tags = append(buildinfo.Tags, "noselfupdate")
}
