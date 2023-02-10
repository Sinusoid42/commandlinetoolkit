package commandlinetoolkit

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
)

func numBytesAvailable() (int, error) {
	cmd := exec.Command("sysctl", "-n", "kern.ipc.pts_nread")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(&out)
	if scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, scanner.Err()
}
