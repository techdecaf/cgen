package app

import (
	"bufio"
	"fmt"
	"os/exec"
)

func execute(name string, arguments ...string) (err error) {
	cmd := exec.Command(name, arguments...)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanErr := bufio.NewScanner(stderr)
	for scanErr.Scan() {
		fmt.Println(scanErr.Text())
	}

	scanOut := bufio.NewScanner(stdout)
	for scanOut.Scan() {
		fmt.Println(scanOut.Text())
	}

	return cmd.Wait()
}
