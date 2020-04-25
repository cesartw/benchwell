// +build linux

package clipboard

import (
	"os"
	"os/exec"
)

func Copy(txt string) {
	goExecutable, err := exec.LookPath("xsel")
	if err != nil {
		return
	}

	cmd := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, "-i", "-b"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	closer, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	closer.Write([]byte(txt))
	closer.Close()

	cmd.Process.Release()
}
