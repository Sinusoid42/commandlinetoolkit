package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

var originalSttyState bytes.Buffer

func getSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func doSth(i int32) {
	err := getSttyState(&originalSttyState)
	if err != nil {
		log.Fatal(err)
	}
	defer setSttyState(&originalSttyState)

	setSttyState(bytes.NewBufferString("cbreak"))
	setSttyState(bytes.NewBufferString("-echo"))

	wg := sync.WaitGroup{}
	wg.Add(1)

	var b []byte = make([]byte, 1)
	t := func() {
		for {
			os.Stdin.Read(b)
			fmt.Printf("Read character: %s\n", b[0])

			if b[0] == byte('q') {
				wg.Add(-1)
			}
		}
	}

	go t()

	wg.Wait()
}

func main() {

	for i := 0; i < 16; i++ {

		for j := 0; j < 16; j++ {

			fmt.Println("\033[48;5m")
			fmt.Println("Test")
		}

	}

	fmt.Println("\033[0m")

}
