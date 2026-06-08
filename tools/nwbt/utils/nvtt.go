package utils

import "os/exec"

type nvtt struct{ command string }

var Nvtt = nvtt{"nvtt_export"}

func (it nvtt) Name() string {
	return it.command
}

func (it nvtt) Check() (string, bool) {
	p, err := exec.LookPath(it.command)
	if err != nil {
		return p, false
	}
	return p, true
}

func (it nvtt) Info() string {
	return "Convert images to/from DDS format (NVIDIA)\n" +
		"download: https://developer.nvidia.com/texture-tools-exporter"
}
