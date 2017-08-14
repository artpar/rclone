// oauthutil parts pre go1.8+

//+build !go1.8

package oauthutil

import "github.com/artpar/rclone/fs"

func (s *authServer) Stop() {
	fs.Debugf(nil, "Closing auth server")
	if s.code != nil {
		close(s.code)
		s.code = nil
	}
	_ = s.listener.Close()
}
