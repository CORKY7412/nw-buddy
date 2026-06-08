package utils

import (
	"os/exec"
)

type ktx struct{ command string }

var Ktx = ktx{"ktx"}

func (it ktx) Name() string {
	return it.command
}

func (it ktx) Check() (string, bool) {
	p, err := exec.LookPath(it.command)
	if err != nil {
		return p, false
	}
	return p, true
}

func (it ktx) Info() string {
	return "Create, validate and transcode KTX texture files\n" +
		"docs:     https://github.com/KhronosGroup/KTX-Software\n" +
		"download: https://github.com/KhronosGroup/KTX-Software/releases"
}

func (it ktx) PresetNormalMap() []string {
	return []string{
		"--generate-mipmap",
		"--encode", "uastc",
		"--uastc-quality", "2",
		"--zstd", "18",
		"--assign-oetf", "linear",
		"--assign-primaries", "none",
		"--format", "R8G8B8_UNORM",
	}
}

func (it ktx) PresetMaskMap() []string {
	return []string{
		"--generate-mipmap",
		"--encode", "uastc",
		"--uastc-quality", "2",
		"--zstd", "18",
		"--assign-oetf", "linear",
		"--assign-primaries", "none",
		"--format", "R8G8B8A8_UNORM",
	}
}

func (it ktx) PresetSpecularMap() []string {
	return []string{
		"--generate-mipmap",
		"--encode", "uastc",
		"--uastc-quality", "2",
		"--zstd", "18",
		"--assign-oetf", "srgb",
		"--assign-primaries", "bt709",
		"--format", "R8G8B8A8_UNORM",
	}
}
func (it ktx) PresetDiffuseMap() []string {
	return []string{
		"--generate-mipmap",
		"--encode", "basis-lz",
		"--assign-oetf", "srgb",
		"--assign-primaries", "bt709",
		"--format", "R8G8B8A8_SRGB",
	}
}

func (it ktx) Run(args ...string) error {
	cmd := exec.Command(it.command, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
