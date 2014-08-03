// Package chefsolo implements the provisioner.Provisoner interface using Chef
// Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/berkshelf"
	"github.com/mlafeldt/chef-runner/cookbook"
	"github.com/mlafeldt/chef-runner/log"
	. "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/rsync"
	"github.com/mlafeldt/chef-runner/util"
)

const (
	DefaultFormat   = "null"
	DefaultLogLevel = "info"
)

type Provisoner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
}

func (p Provisoner) prepareJSON() error {
	log.Debug("Preparing JSON data")
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(SandboxPathTo("dna.json"), []byte(data), 0644)
}

func (p Provisoner) prepareSoloConfig() error {
	log.Debug("Preparing Chef Solo config")
	data := fmt.Sprintf("cookbook_path \"%s\"\n", RootPathTo("cookbooks"))
	return ioutil.WriteFile(SandboxPathTo("solo.rb"), []byte(data), 0644)
}

func (p Provisoner) resolveWithBerkshelf() error {
	log.Info("Installing cookbooks with Berkshelf")
	return berkshelf.Install(SandboxPathTo("cookbooks"))
}

func (p Provisoner) copyThisCookbook() error {
	cb, err := cookbook.NewCookbook(".")
	if err != nil {
		return err
	}
	files, err := cb.Files()
	if err != nil {
		return err
	}
	log.Info("Updating", cb.Name, "cookbook with rsync")
	c := rsync.Client{
		Archive:  true,
		Delete:   true,
		Compress: true,
		Verbose:  true,
	}
	return c.Copy(path.Join(SandboxPathTo("cookbooks"), cb.Name), files...)
}

func (p Provisoner) prepareCookbooks() error {
	if !util.FileExist(SandboxPathTo("cookbooks")) {
		return p.resolveWithBerkshelf()
	}
	return p.copyThisCookbook()
}

func (p Provisoner) CreateSandbox() error {
	if err := CreateSandbox(); err != nil {
		return err
	}
	if err := p.prepareJSON(); err != nil {
		return err
	}
	if err := p.prepareSoloConfig(); err != nil {
		return err
	}
	return p.prepareCookbooks()
}

func (p Provisoner) CleanupSandbox() error {
	return CleanupSandbox()
}

func (p Provisoner) Command() []string {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}
	logLevel := p.LogLevel
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}
	return []string{
		"sudo", "chef-solo",
		"--config", RootPathTo("solo.rb"),
		"--json-attributes", RootPathTo("dna.json"),
		"--override-runlist", strings.Join(p.RunList, ","),
		"--format", format,
		"--log_level", logLevel,
	}
}