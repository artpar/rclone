package tree

import (
	"bytes"
	"testing"

	"github.com/a8m/tree"
	_ "github.com/artpar/rclone/backend/local"
	"github.com/artpar/rclone/fs"
	"github.com/artpar/rclone/fstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	fstest.Initialise()

	buf := new(bytes.Buffer)

	f, err := fs.NewFs("testfiles")
	require.NoError(t, err)
	err = Tree(f, buf, new(tree.Options))
	require.NoError(t, err)
	assert.Equal(t, `/
├── file1
├── file2
├── file3
└── subdir
    ├── file4
    └── file5

1 directories, 5 files
`, buf.String())
}
