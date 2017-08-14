package all

import (
	// Active file systems
	_ "github.com/artpar/rclone/amazonclouddrive"
	_ "github.com/artpar/rclone/b2"
	_ "github.com/artpar/rclone/box"
	_ "github.com/artpar/rclone/crypt"
	_ "github.com/artpar/rclone/drive"
	_ "github.com/artpar/rclone/dropbox"
	_ "github.com/artpar/rclone/ftp"
	_ "github.com/artpar/rclone/googlecloudstorage"
	_ "github.com/artpar/rclone/http"
	_ "github.com/artpar/rclone/hubic"
	_ "github.com/artpar/rclone/local"
	_ "github.com/artpar/rclone/onedrive"
	_ "github.com/artpar/rclone/qingstor"
	_ "github.com/artpar/rclone/s3"
	_ "github.com/artpar/rclone/sftp"
	_ "github.com/artpar/rclone/swift"
	_ "github.com/artpar/rclone/yandex"
)
