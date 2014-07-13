package rsync_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/stretchr/testify/assert"
)

var copyTests = []struct {
	client rsync.Client
	src    []string
	dst    string
	cmd    []string
}{
	{
		rsync.Client{},
		[]string{"a"}, "b",
		[]string{"rsync", "a", "b"},
	},
	{
		rsync.Client{},
		[]string{"a", "b"}, "c",
		[]string{"rsync", "a", "b", "c"},
	},
	{
		rsync.Client{Archive: true},
		[]string{"a"}, "b",
		[]string{"rsync", "--archive", "a", "b"},
	},
	{
		rsync.Client{Delete: true},
		[]string{"a"}, "b",
		[]string{"rsync", "--delete", "a", "b"},
	},
	{
		rsync.Client{Verbose: true},
		[]string{"a"}, "b",
		[]string{"rsync", "--verbose", "a", "b"},
	},
	{
		rsync.Client{Exclude: []string{"x", "y"}},
		[]string{"a"}, "b",
		[]string{"rsync", "--exclude", "x", "--exclude", "y", "a", "b"},
	},
	{
		rsync.Client{Archive: true, Delete: true, Exclude: []string{"x"}},
		[]string{"a"}, "b",
		[]string{"rsync", "--archive", "--delete", "--exclude", "x", "a", "b"},
	},
}

func TestCommand(t *testing.T) {
	for _, test := range copyTests {
		cmd, err := test.client.Command(test.src, test.dst)
		if assert.NoError(t, err) {
			assert.Equal(t, test.cmd, cmd)
		}
	}
}

func TestCopy_MissingSource(t *testing.T) {
	err := rsync.DefaultClient.Copy([]string{}, "a/b")
	assert.EqualError(t, err, "No source given")
}

func TestCopy_MissingDestination(t *testing.T) {
	err := rsync.DefaultClient.Copy([]string{"a"}, "")
	assert.EqualError(t, err, "No destination given")
}
