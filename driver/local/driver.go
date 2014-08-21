// Package local implements a driver for the local system.
package local

import (
	"fmt"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/rsync"
)

// Driver is a driver for the local system.
type Driver struct {
}

// NewDriver creates a new local driver.
func NewDriver() (*Driver, error) {
	return &Driver{}, nil
}

// RunCommand runs the specified command on the local system.
func (drv Driver) RunCommand(args []string) error {
	return exec.RunCommand(args)
}

// Upload copies files to the right place on the local system.
func (drv Driver) Upload(dst string, src ...string) error {
	return rsync.MirrorClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("Local driver")
}
