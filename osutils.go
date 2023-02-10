package commandlinetoolkit

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
)

func numBytesAvailable() int {
	cmd := exec.Command("sysctl", "-n", "kern.ipc.pts_nread")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0
	}
	scanner := bufio.NewScanner(&out)
	if scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0
		}
		return n
	}
	return 0
}
