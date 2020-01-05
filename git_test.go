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

func TestCommit(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	err := c.Init()
	assert.NoError(t, err)
	assert.NoError(t, c.ConfigureUser("barry", "barry@starlabs.org"))
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey"), 0744))
	assert.NoError(t, c.IndexAll())
	assert.NoError(t, c.Commit("first"))
	s, err := c.Status()
	assert.NoError(t, err)
	assert.Empty(t, s)
}

func TestLog(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	err := c.Init()
	assert.NoError(t, err)
	assert.NoError(t, c.ConfigureUser("barry", "barry@starlabs.org"))
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey"), 0744))
	assert.NoError(t, c.IndexAll())
	assert.NoError(t, c.Commit("first"))

	l, err := c.Log()
	assert.NoError(t, err)
	assert.Len(t, l, 1)

	// Create the second commit
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey hey"), 0744))
	assert.NoError(t, c.IndexAll())
	assert.NoError(t, c.Commit("second"))
	l, err = c.Log()
	assert.NoError(t, err)
	assert.Len(t, l, 2)
}

func TestRevParse(t *testing.T) {
	setup()
	c := NewCLI(DataPath)
	assert.NoError(t, c.Init())
	assert.NoError(t, c.ConfigureUser("barry", "barry@starlabs.org"))
	assert.NoError(t, ioutil.WriteFile(path.Join(DataPath, "readme.md"), []byte("#hey"), 0744))
	assert.NoError(t, c.IndexAll())
	assert.NoError(t, c.Commit("first"))
	sha, err := c.RevParse("head")
	assert.NoError(t, err)
	assert.NotEmpty(t, sha)
}
