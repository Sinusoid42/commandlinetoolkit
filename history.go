package commandlinetoolkit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const HISTORY_FILENAME = ".history"

type history struct {
	_historyFileName string
	_programName     string
	_theHistoryFile  *os.File
	_writer          *bufio.Writer

	_bufferedLines [][]byte
	_keyLines      [][]Key

	_currHistoryIndex int

	_enabledHistoryFile bool

	_verboseColor Color

	_debugPrintPrefix string

	_rwFlag int

	_error   error
	_verbose CLICODE
}

func (h *history) create() {

	h._theHistoryFile, h._error = os.Create(h._historyFileName)

	if h._error != nil {

		h._rwFlag = -2
		if h._verbose&CLI_VERBOSE_FILE > 0 {
			h.printVerbose("\n--> hF-Handler: createFile: ERROR >> " + h._error.Error() + "\n")
		}
	}
}

func (h *history) open() {

	//h._theHistoryFile, h._error = os.OpenFile(h._historyFileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)

	h._theHistoryFile, h._error = os.OpenFile(".history", os.O_CREATE|os.O_RDWR, 0666)

	if h._verbose&CLI_VERBOSE_FILE > 0 {
		h.printVerbose("\n--> hF-Handler: openFile: " + h._theHistoryFile.Name())
	}
	if h._error != nil {
		//file probably already existing
		h._rwFlag = -1
		if h._verbose&CLI_VERBOSE_FILE > 0 {
			h.printVerbose("\n--> hF-Handler: openFile: ERROR >> " + h._error.Error() + "\n")
		}
		//now we have to create file
		//h.create()
	}

	//read
	h.read()

	//overwrite so we can keep appending
	h.create()

	h._writer = bufio.NewWriter(h._theHistoryFile)

	for _, line := range h._bufferedLines {
		_, err := h._writer.WriteString(string(line) + "\n")

		if err != nil {
			h._error = err
			if h._verbose&CLI_VERBOSE_FILE > 0 {
				h.printVerbose("\n--> hF-Handler: writeFile: ERROR >> " + h._error.Error() + "\n")
			}
		}
	}

	h._writer.Flush()

}

func (h *history) append(newLine string) {
	if len(newLine) < 1 {
		return
	}
	_, err := h._writer.WriteString(newLine + "\n")
	if err != nil {
		h._error = err

		if h._verbose&CLI_VERBOSE_FILE > 0 {
			h.printVerbose("\n--> hF-Handler: appendFile: ERROR >> " + h._error.Error() + "\n")
		}

	}

	h._writer.Flush()
}

/*
*
close the file, needs to be done in the end of the program, when the shell is closing
*/
func (h *history) close() {

	if h._verbose&CLI_VERBOSE_FILE > 0 {
		h.printVerbose("\n--> hF-Handler: closeFile: ERROR ")
	}

	h._writer.Flush()

	defer h._theHistoryFile.Close()
}

func (h *history) read() [][]Key {

	//reader := bufio.NewReader(h._theHistoryFile)

	reader := bufio.NewScanner(h._theHistoryFile)
	//for scanner.Scan() {
	//	line := scanner.Text()
	//	println(line)

	keylines := [][]Key{}

	for reader.Scan() {

		line := reader.Bytes()

		h._bufferedLines = append(h._bufferedLines, line)

		if !strings.ContainsAny(string(line), h._debugPrintPrefix) {
			k := []Key{}
			for _, b := range line {
				k = append(k, Key(b))
			}
			keylines = append(keylines, k)
		}
	}
	//once here, we have to revert the

	h._keyLines = [][]Key{}
	//l := len(keylines)
	for i, _ := range keylines {

		h._keyLines = append(h._keyLines, keylines[i])

	}

	return h._keyLines
}

func newHistoryFileHandler(_programName string) *history {

	h := &history{
		_historyFileName:    HISTORY_FILENAME,
		_programName:        _programName,
		_enabledHistoryFile: true,
		_error:              nil,
		_verbose:            0,
		_verboseColor:       COLOR_PINK_I,
		_debugPrintPrefix:   "#",
		_rwFlag:             -1,
		_bufferedLines:      [][]byte{},
		_currHistoryIndex:   -1,
	}

	h.open()

	return h
}

/*
*
Prints with the verbose color overlay
*/
func (h *history) printVerbose(str interface{}) {
	fmt.Print(h._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}

/*
*
Iterate the previous history in the present shell
*/
func (h *history) iterateHistory(s *shellHandler) {
	if s._enabledHistory && (s._arrowAction == 2 || s._arrowAction == 3) {
		linputs := len(s._previnputs)

		if h._currHistoryIndex >= 0 && linputs > h._currHistoryIndex {

			s._inputDisplayBuffer = []Key{}

			for i, _ := range s._previnputs[linputs-1-h._currHistoryIndex] {
				s._inputDisplayBuffer = append(s._inputDisplayBuffer, s._previnputs[linputs-1-h._currHistoryIndex][i])
			}

			s._rtAction = 0
		} else {
			if -1 >= h._currHistoryIndex {
				s._inputDisplayBuffer = []Key{}
				s.clearCurrentLine()
				s.printPrefix()
			}
			if h._currHistoryIndex < -1 {
				h._currHistoryIndex = -1
				s._alert = true
			}
			if h._currHistoryIndex >= linputs-1 {
				h._currHistoryIndex = linputs - 1
				s._alert = true
			}
		}

		if s._alert {
			s._alert = false
			if s._playAlert {
				fmt.Print("\a")
			}
		}
		s.reprintCurrentLine()
	}
}

func (h *history) up() bool {
	h._currHistoryIndex--

	return h._currHistoryIndex > 0
}

func (h *history) down() bool {
	h._currHistoryIndex++

	return h._currHistoryIndex > 0
}

func (h *history) reset() {
	h._currHistoryIndex = -1
}
