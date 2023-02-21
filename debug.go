package commandlinetoolkit

import "fmt"

type debugHandler struct {
	_verboseColor Color

	_verbose CLICODE
}

func newDebugHandler() *debugHandler {

	d := &debugHandler{
		//_verboseColor: GenColor("3", "", "2"),
		_verboseColor: COLOR_GRAY_debug,
		_verbose:      0,
	}

	return d
}

func (d *debugHandler) printInputBuffer(buffer []Key) {
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, "\n-->shell: Previous parseable input: ")
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, buffer)
	d.printVerboseBuffer(buffer)

}

func (d *debugHandler) printByteBuffer(b []byte) {
	d.printVerbose(CLI_VERBOSE_DEBUG, b)
}

/*
DEPRECATED
Prints with the verbose color overlay

func (d *debugHandler) printVerbose(str interface{}) {
	if !(d._verbose&CLI_VERBOSE_DEBUG > 0) {
		return
	}
	fmt.Print(d._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}
*/

/*
*
Prints with the verbose color overlay
*/
func (d *debugHandler) printVerbose(code CLICODE, str interface{}) {
	if !(d._verbose&code > 0) {
		return
	}
	fmt.Print(d._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}

/*
*
Prints with vibrant red error code
*/
func (d *debugHandler) printError(str interface{}) {
	fmt.Print(COLOR_RED_I, "Error:\n", str, COLOR_RESET)

}

/*
*
Prints the current input buffer with spaces, in debug mode
*/
func (d *debugHandler) printVerboseBuffer(buffer []Key) {
	sr := "["
	for _, i := range buffer {
		if i == KEY_ESC {
			sr += " ESC,"
		} else if i == KEY_DELETE {
			sr += " DEL,"
		} else {
			sr += "" + string(i) + ", "
		}
	}
	l := len(sr)
	if l > 1 {
		sr = sr[:l-2]
	}
	sr += "]"
	d.printVerbose(CLI_VERBOSE_SHELL, sr)
}

func (d *debugHandler) boot() {
	d.printVerbose(CLI_VERBOSE_SHELL, "\n-->shell: Booted shell subroutine")

}

/*
*************************************
debugging with codes
*/
func (d *debugHandler) debugBufferSingle(k Key) {

	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, "\n-->shell: numBytes: ")
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, numBytesAvailable())
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, "\n")
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, "-->shell: inputbyte: ")
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, []byte{byte(k)})
	d.printVerbose(CLI_VERBOSE_DEBUG, "\n")

}

/*
*************************************
debugging with codes
*/
func (d *debugHandler) debugBuffer(k []Key) {

	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, "\n")
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, "-->shell: inputbyte: ")
	d.printVerbose(CLI_VERBOSE_SHELL_BUFFER, k)
	d.printVerbose(CLI_VERBOSE_DEBUG, "\n")

}

func (d *debugHandler) debugReturn() {
	d.printVerbose(CLI_VERBOSE_SHELL, "\n-->shell: Registered RETURN KEY\n")

}

/*
*
General debug function to print debug messages in the printVerbose handle function
*/
func (d *debugHandler) debug(verbose CLICODE, msg string) {

	d.printVerbose(CLI_VERBOSE_SHELL, msg)

}

func (d *debugHandler) debugSuggestions(s *shellHandler) {
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, "\n-->shell: Requesting current-layer suggestions: ")
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, "Layer")
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, s._parseDepth)
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, "; Request: ")
	d.printVerbose(CLI_VERBOSE_SHELL_PARSE, s._requestSuggestions)

}
