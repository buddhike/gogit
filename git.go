package git

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

var statusStringRegexp *regexp.Regexp

func init() {
	statusStringRegexp = regexp.MustCompile(`^(\?\?|A|M|D|R)\s+(.*)$`)
}

// Status represents the kind of change made to tracked file
type Status string

const (
	// StatusUntracked file
	StatusUntracked Status = "U"
	// StatusAdded file
	StatusAdded Status = "A"
	// StatusModified file
	StatusModified Status = "M"
	// StatusRenamed file
	StatusRenamed Status = "R"
	// StatusDeleted file
	StatusDeleted Status = "D"
)

var statusTable map[string]Status = map[string]Status{
	"??": StatusUntracked,
	"A":  StatusAdded,
	"M":  StatusModified,
	"D":  StatusDeleted,
	"R":  StatusRenamed,
}

// CLI maintains the state of currently opened repository
type CLI struct {
	path string
}

// StatusEntry represents an entry reported by status command
type StatusEntry struct {
	Path   string
	Status Status
}

// NewCLI creates a new instance of CLI
func NewCLI(path string) CLI {
	return CLI{
		path: path,
	}
}

// Version returns the git cli version
func (c CLI) Version() (string, error) {
	r, err := c.runCommand("version")
	if err != nil {
		return "", err
	}

	return r[0], nil
}

// Init initializes a new repository
func (c CLI) Init() error {
	_, err := c.runCommand("init")
	return err
}

// Status runs git status command
func (c CLI) Status() ([]StatusEntry, error) {
	l, err := c.runCommand("status", "-s")
	if err != nil {
		return nil, err
	}

	r := make([]StatusEntry, len(l))
	for i, e := range l {
		matches := statusStringRegexp.FindStringSubmatch(e)
		if matches == nil {
			return nil, errors.New("Unable to parse status string")
		}
		r[i] = StatusEntry{
			Status: statusTable[matches[1]],
			Path:   matches[2],
		}
	}

	return r, nil
}

// IndexAll stages all changes in workspace
func (c CLI) IndexAll() error {
	_, err := c.runCommand("add", "-A")
	return err
}

// Commit creates a commit with specified message
func (c CLI) Commit(message string) error {
	_, err := c.runCommand("commit", "-m", message)
	return err
}

// ConfigureUser setup identity for commits etc
func (c CLI) ConfigureUser(username, email string) error {
	_, err := c.runCommand("config", "--local", "user.name", username)
	if err != nil {
		return err
	}
	_, err = c.runCommand("config", "--local", "user.email", email)
	return err
}

// MergeBase returns the merge base of two commits
func (c CLI) MergeBase(first, second string) (string, error) {
	r, err := c.runCommand("merge-base", first, second)
	if err != nil {
		return "", err
	}
	return r[0], nil
}

// Log returns the commit log
func (c CLI) Log() ([]string, error) {
	return c.runCommand("log", "--pretty=\"%H\"")
}

// Checkout checks out the specified commit sha
func (c CLI) Checkout(sha string) error {
	return nil
}

// runCommand implements the driver for running git with specified arguments
// and parsing its output
func (c CLI) runCommand(command string, arg ...string) ([]string, error) {
	cmd := exec.Command("git", append([]string{command}, arg...)...)
	cmd.Dir = c.path
	stdout, err := cmd.Output()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			return nil, err
		}

		errorLines, err := readLines(exitErr.Stderr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.Join(errorLines, ";"))
	}

	return readLines(stdout)
}

func readLines(buf []byte) ([]string, error) {
	out := []string{}
	stdoutScanner := bufio.NewScanner(bytes.NewBuffer(buf))
	for stdoutScanner.Scan() {
		out = append(out, stdoutScanner.Text())
	}
	if err := stdoutScanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
