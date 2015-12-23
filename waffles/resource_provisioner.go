package waffles

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/armon/circbuf"
	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/go-linereader"
	"github.com/mitchellh/mapstructure"
)

const (
	// maxBufSize limits how much output we collect from a local
	// invocation. This is to prevent TF memory usage from growing
	// to an enormous amount due to a faulty process.
	maxBufSize = 8 * 1024

	// wafflesExec is the default location to find waffles.
	wafflesExec = "/etc/waffles/waffles.sh"
)

// Waffles represents a Waffles configuration
type Waffles struct {
	Debug         bool   `mapstructure:"debug"`
	Host          string `mapstructure:"host"`
	PrivateKey    string `mapstructure:"private_key"`
	RemoteDir     string `mapstructure:"remote_dir"`
	Retry         int    `mapstructure:"retry"`
	Role          string `mapstricture:"role"`
	SiteDirectory string `mapstructure:"site_directory"`
	Sudo          bool   `mapstructure:"sudo"`
	User          string `mapstructure:"user"`
	WafflesExec   string `mapstructure:"waffles_exec"`
	Wait          int    `mapstructure:"wait"`
}

// ResourceProvisioner represents a generic Waffles provisioner
type ResourceProvisioner struct{}

// Apply executes waffles
func (p *ResourceProvisioner) Apply(
	o terraform.UIOutput,
	s *terraform.InstanceState,
	c *terraform.ResourceConfig) error {

	w, err := p.decodeConfig(c)
	if err != nil {
		return err
	}

	// Execute waffles via a shell
	if runtime.GOOS == "windows" {
		return fmt.Errorf("Waffles is not supported on Windows at this time.")
	}

	pr, pw := io.Pipe()
	copyDoneCh := make(chan struct{})
	go p.copyOutput(o, pr, copyDoneCh)

	// Setup the command

	// required flags
	flags := []string{"-s", w.Host, "-r", w.Role}

	// optional flags
	if w.Debug {
		flags = append(flags, "-d")
	}

	if w.PrivateKey != "" {
		flags = append(flags, "-k", w.PrivateKey)
	}

	if w.RemoteDir != "" {
		flags = append(flags, "-z", w.RemoteDir)
	}

	if w.Retry != 0 {
		flags = append(flags, "-c", strconv.Itoa(w.Retry))
	}

	if w.Sudo {
		flags = append(flags, "-y")
	}

	if w.User != "" {
		flags = append(flags, "-u", w.User)
	}

	if w.Wait != 0 {
		flags = append(flags, "-w", strconv.Itoa(w.Wait))
	}

	// run the command
	cmd := exec.Command(w.WafflesExec, flags...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("WAFFLES_SITE_DIR=%s", w.SiteDirectory))
	output, _ := circbuf.NewBuffer(maxBufSize)
	cmd.Stderr = io.MultiWriter(output, pw)
	cmd.Stdout = io.MultiWriter(output, pw)

	// Output what we're about to run
	fullCommand := fmt.Sprintf("WAFFLES_SITE_DIR=%s %s %s", w.SiteDirectory, w.WafflesExec, strings.Join(flags, " "))
	o.Output(fmt.Sprintf("Executing: %s", fullCommand))

	// Run the command to completion
	err = cmd.Run()

	// Close the write-end of the pipe so that the goroutine mirroring output
	// ends properly.
	pw.Close()
	<-copyDoneCh

	if err != nil {
		return fmt.Errorf("Error running command '%s': %v. Output: %s", fullCommand, err, output.Bytes())
	}

	return nil

}

// Validate checks if the required arguments are configured
func (p *ResourceProvisioner) Validate(c *terraform.ResourceConfig) ([]string, []error) {
	validator := config.Validator{
		Required: []string{"host", "role", "site_directory"},
		Optional: []string{"debug", "remote_dir", "private_key", "retry", "sudo", "user", "waffles_exec", "wait"},
	}
	return validator.Validate(c)
}

func (p *ResourceProvisioner) decodeConfig(c *terraform.ResourceConfig) (*Waffles, error) {
	w := new(Waffles)

	decConf := &mapstructure.DecoderConfig{
		ErrorUnused:      true,
		WeaklyTypedInput: true,
		Result:           w,
	}
	dec, err := mapstructure.NewDecoder(decConf)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for k, v := range c.Raw {
		m[k] = v
	}

	for k, v := range c.Config {
		m[k] = v
	}

	if err := dec.Decode(m); err != nil {
		return nil, err
	}

	// fill in the blanks
	if w.WafflesExec == "" {
		w.WafflesExec = wafflesExec
	}

	// Expand home on possible areas
	w.SiteDirectory, err = homedir.Expand(w.SiteDirectory)
	if err != nil {
		return nil, err
	}

	w.PrivateKey, err = homedir.Expand(w.PrivateKey)
	if err != nil {
		return nil, err
	}

	return w, nil

}

func (p *ResourceProvisioner) copyOutput(
	o terraform.UIOutput, r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)
	lr := linereader.New(r)
	for line := range lr.Ch {
		o.Output(line)
	}
}
