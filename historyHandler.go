package commandlinetoolkit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const HISTORY_FILENAME = ".history"

type historyHandler struct {
	_historyFileName string
	_programName     string
	_theHistoryFile  *os.File
	_writer          *bufio.Writer

	_bufferedLines [][]byte
	_keyLines      [][]Key

	_enabledHistoryFile bool

	_verboseColor Color

	_debugPrintPrefix string

	_rwFlag int

	_error   error
	_verbose int
}

func (h *historyHandler) create() {

	h._theHistoryFile, h._error = os.Create(h._historyFileName)

	if h._error != nil {

		h._rwFlag = -2
		if h._verbose&CLI_VERBOSE_FILE > 0 {
			h.printVerbose("\n--> hF-Handler: createFile: ERROR >> " + h._error.Error() + "\n")
		}
	}
}

func (h *historyHandler) open() {

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

	h.read()

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

func (h *historyHandler) append(newLine string) {
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
func (h *historyHandler) close() {

	if h._verbose&CLI_VERBOSE_FILE > 0 {
		h.printVerbose("\n--> hF-Handler: closeFile: ERROR ")
	}

	h._writer.Flush()

	defer h._theHistoryFile.Close()
}

func (h *historyHandler) read() [][]Key {

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
	h._keyLines = keylines
	return keylines
}

func newHistoryFileHandler(_programName string) *historyHandler {

	h := &historyHandler{
		_historyFileName:    HISTORY_FILENAME,
		_programName:        _programName,
		_enabledHistoryFile: true,
		_error:              nil,
		_verbose:            0,
		_verboseColor:       COLOR_PINK_I,
		_debugPrintPrefix:   "#",
		_rwFlag:             -1,
		_bufferedLines:      [][]byte{},
	}
	h.open()
	return h
}

/*
*
Prints with the verbose color overlay
*/
func (h *historyHandler) printVerbose(str interface{}) {
	fmt.Print(h._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}
