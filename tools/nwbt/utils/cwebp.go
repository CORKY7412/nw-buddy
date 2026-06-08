package utils

import (
	"os/exec"
)

type cwebp struct{ command string }

var Cwebp = cwebp{"cwebp"}

func (it cwebp) Name() string {
	return it.command
}

func (it cwebp) Check() (string, bool) {
	p, err := exec.LookPath(it.command)
	if err != nil {
		return "p", false
	}
	return p, true
}

func (it cwebp) Info() string {
	return "Convert images to WebP format\n" +
		"docs:     https://developers.google.com/speed/webp/docs/cwebp\n" +
		"download: https://developers.google.com/speed/webp/docs/precompiled"
}

func (it cwebp) Run(args ...string) error {
	cmd := exec.Command(it.command, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
