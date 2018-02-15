package all

import (
	// Active file systems
	_ "github.com/artpar/rclone/backend/alias"
	_ "github.com/artpar/rclone/backend/amazonclouddrive"
	_ "github.com/artpar/rclone/backend/azureblob"
	_ "github.com/artpar/rclone/backend/b2"
	_ "github.com/artpar/rclone/backend/box"
	_ "github.com/artpar/rclone/backend/cache"
	_ "github.com/artpar/rclone/backend/crypt"
	_ "github.com/artpar/rclone/backend/drive"
	_ "github.com/artpar/rclone/backend/dropbox"
	_ "github.com/artpar/rclone/backend/ftp"
	_ "github.com/artpar/rclone/backend/googlecloudstorage"
	_ "github.com/artpar/rclone/backend/http"
	_ "github.com/artpar/rclone/backend/hubic"
	_ "github.com/artpar/rclone/backend/local"
	_ "github.com/artpar/rclone/backend/onedrive"
	_ "github.com/artpar/rclone/backend/pcloud"
	_ "github.com/artpar/rclone/backend/qingstor"
	_ "github.com/artpar/rclone/backend/s3"
	_ "github.com/artpar/rclone/backend/sftp"
	_ "github.com/artpar/rclone/backend/swift"
	_ "github.com/artpar/rclone/backend/webdav"
	_ "github.com/artpar/rclone/backend/yandex"
)
