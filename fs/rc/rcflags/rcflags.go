// Package rcflags implements command line flags to set up the remote control
package rcflags

import (
	"github.com/artpar/rclone/cmd/serve/httplib/httpflags"
	"github.com/artpar/rclone/fs/config/flags"
	"github.com/artpar/rclone/fs/rc"
	"github.com/spf13/pflag"
)

// Options set by command line flags
var (
	Opt = rc.DefaultOpt
)

// AddFlags adds the remote control flags to the flagSet
func AddFlags(flagSet *pflag.FlagSet) {
	rc.AddOption("rc", &Opt)
	flags.BoolVarP(flagSet, &Opt.Enabled, "rc", "", false, "Enable the remote control server.")
	flags.StringVarP(flagSet, &Opt.Files, "rc-files", "", "", "Path to local files to serve on the HTTP server.")
	flags.BoolVarP(flagSet, &Opt.Serve, "rc-serve", "", false, "Enable the serving of remote objects.")
	flags.BoolVarP(flagSet, &Opt.NoAuth, "rc-no-auth", "", false, "Don't require auth for certain methods.")
	flags.BoolVarP(flagSet, &Opt.WebUI, "rc-web-gui", "", false, "Launch WebGUI on localhost")
	flags.BoolVarP(flagSet, &Opt.WebGUIUpdate, "rc-web-gui-update", "", false, "Update / Force update to latest version of web gui")
	flags.StringVarP(flagSet, &Opt.WebGUIFetchURL, "rc-web-fetch-url", "", "https://api.github.com/repos/rclone/rclone-webui-react/releases/latest", "URL to fetch the releases for webgui.")
	flags.StringVarP(flagSet, &Opt.AccessControlAllowOrigin, "rc-allow-origin", "", "", "Set the allowed origin for CORS.")
	flags.DurationVarP(flagSet, &Opt.JobExpireDuration, "rc-job-expire-duration", "", Opt.JobExpireDuration, "expire finished async jobs older than this value")
	flags.DurationVarP(flagSet, &Opt.JobExpireInterval, "rc-job-expire-interval", "", Opt.JobExpireInterval, "interval to check for expired async jobs")
	httpflags.AddFlagsPrefix(flagSet, "rc-", &Opt.HTTPOptions)
}
