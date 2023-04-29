package commandlinetoolkit

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// operating system helper struct for sys signals and callbacks
type osHandler struct {

	//os reading state of the previous host shell
	_originalSttyState *bytes.Buffer

	_sysCall syscall.Signal

	_sysCallInterrupt syscall.Signal

	_wg sync.WaitGroup

	_sysSignal chan os.Signal
}

func newOSHandler() *osHandler {
	o := &osHandler{

		_originalSttyState: &bytes.Buffer{},

		_sysCallInterrupt: syscall.SIGINT,
		//most recent syscall input
		_sysCall: 0,
	}
	return o
}

func (o *osHandler) registerSystemSignalCallbacks(s *shellHandler) {

	s._osHandler._sysSignal = make(chan os.Signal, 1)
	signal.Notify(s._osHandler._sysSignal, os.Interrupt)
	//the shell callback
	f := func() {

		//need double loop as syscall can be happening in different scope at certain times

		if s._debugHandler._verbose&CLI_VERBOSE_OS_SIG > 0 {
			fmt.Println("osHandler: Booted os signal handling subroutine")
		}
		index := 0
		for sig := range s._osHandler._sysSignal {
			// sig is a ^C, handle it

			if sig == nil {
				continue
			}

			if sig == syscall.SIGINT {
				index++
				if s._debugHandler._verbose&CLI_VERBOSE_OS_SIG > 0 {
					fmt.Println("\n-->osHandler: syscall.SIGINT")
				}
				fmt.Println(index)
				fmt.Println("Keyboard Interrupt")
				fmt.Println("Exit? y/n")

				s.printPrefix()

				s._osHandler._sysCall = syscall.SIGINT //sysExit
				//reset buffers
				//s._currInput = s._lastInput
				s._inputDisplayBuffer = []Key{}

				for {

					if s._exit == 1 {
						if s._debugHandler._verbose&CLI_VERBOSE_OS_SIG > 0 {
							fmt.Println("\n-->osHandler: Exiting out of os handling subroutine")
						}

						//run at the end once
						s._osHandler._wg.Add(-1)

						return
					}

					if s._osHandler._sysCall == 0 || s._exit == 0 {
						break
					}
				}
			}
		}
	}
	go f()
	s._osHandler._wg.Add(1)
}

func (o *osHandler) exit() {
	setSttyState(o._originalSttyState)
	//reset raw
	setSttyState(bytes.NewBufferString("-raw"))
	setSttyState(bytes.NewBufferString("-icanon"))
}

func getSttyState(state *bytes.Buffer) (err error) {
	//https://gist.github.com/mrnugget/9582788
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setSttyState(state *bytes.Buffer) (err error) {
	//https://gist.github.com/mrnugget/9582788
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (o *osHandler) removeTerminalBuffering() {
	//https://gist.github.com/mrnugget/9582788
	err := getSttyState(o._originalSttyState)
	if err != nil {
		log.Fatal(err)
	}
	//run in this function!
	//"/dev/tty", "raw", "-echo", "cbreak", "-g"
	//setSttyState(bytes.NewBufferString("icanon raw cbreak -g"))
	//setSttyState(bytes.NewBufferString(" cbreak -echo"))
	//setSttyState(bytes.NewBufferString("min 1"))
	//setSttyState(bytes.NewBufferString("-raw"))
	//setSttyState(bytes.NewBufferString("cbreak"))
	//setSttyState(bytes.NewBufferString("-echo"))
	setSttyState(bytes.NewBufferString("-raw"))
	setSttyState(bytes.NewBufferString("cbreak"))
	setSttyState(bytes.NewBufferString("-echo"))
}

func (o *osHandler) reset() {
	setSttyState(bytes.NewBufferString("-raw"))
	setSttyState(bytes.NewBufferString("echo"))
}

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
