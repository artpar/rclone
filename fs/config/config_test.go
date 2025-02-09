// These are in an external package because we need to import configfile

package config_test

import (
	"testing"

	"github.com/artpar/rclone/fs/config"
	"github.com/artpar/rclone/fs/config/configfile"
	"github.com/stretchr/testify/assert"
)

func init() {
	configfile.Install()
}

func TestConfigLoad(t *testing.T) {
	oldConfigPath := config.GetConfigPath()
	assert.NoError(t, config.SetConfigPath("./testdata/plain.conf"))
	defer func() {
		assert.NoError(t, config.SetConfigPath(oldConfigPath))
	}()
	config.ClearConfigPassword()
	sections := config.Data().GetSectionList()
	var expect = []string{"RCLONE_ENCRYPT_V0", "nounc", "unc"}
	assert.Equal(t, expect, sections)

	keys := config.Data().GetKeyList("nounc")
	expect = []string{"type", "nounc"}
	assert.Equal(t, expect, keys)
}
