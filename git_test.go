package git

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DataPath string = "./.data"

func setup() {
	os.RemoveAll(DataPath)
	os.Mkdir(DataPath, 0744)
}

func TestVersion(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	v, err := c.Version()
	assert.NoError(t, err)
	assert.Contains(t, v, "git version")
}

func TestInit(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	err := c.Init()
	assert.NoError(t, err)
	s, err := c.Status()
	assert.Empty(t, s)
}

func TestStatus(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	err := c.Init()
	assert.NoError(t, err)
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey"), 0744))
	s, err := c.Status()
	assert.NoError(t, err)
	assert.Equal(t, "readme.md", s[0].Path)
	assert.Equal(t, StatusUntracked, s[0].Status)
}

func TestIndexAll(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	err := c.Init()
	assert.NoError(t, err)
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey"), 0744))
	assert.NoError(t, c.IndexAll())
	s, err := c.Status()
	assert.NoError(t, err)
	assert.Equal(t, "readme.md", s[0].Path)
	assert.Equal(t, StatusAdded, s[0].Status)
}
