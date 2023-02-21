package commandlinetoolkit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
Helper file that represents a runnable program,

	_programName        The name of the executeable file, is parsed when getting os.Args from the first run
	_programFile        The name of the configuration json file, that holds all the available Arguments for the CLI Parser
	_program            The command line template struct, that holds the interface storage for the actual program
	_verbose            A given debugging verbosity that is passed to the debug handler
	_debugHandler       The debug handler that manages the printing and output
*/

type program struct {
	_programName  string
	_programFile  string
	_program      *commandlinetemplate
	_verbose      CLICODE
	_debugHandler *debugHandler
}

// build a new program
func newprogram(filename string) *program {

	p := &program{
		_programName:  "Command Line: " + VERSION,
		_programFile:  filename,
		_program:      DefaultCommandLineTemplate(),
		_debugHandler: newDebugHandler(),
	}

	p._debugHandler._verbose = 0

	return p
}

// read a program from a json
func (p *program) readJsonProgram(filename string) string {

	//check for file forst
	if _, err := os.OpenFile(filename, os.O_RDONLY, 0666); err != nil {
		filename = p._programFile
		p._debugHandler.printError("-->readJson: File not available, falling back to default\n")
	}

	p._programFile = filename

	_programFile, _error := os.OpenFile(p._programFile, os.O_CREATE, 0666)

	if _error != nil {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Error\n")

	} else {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Read Success\n")
	}

	fileInfo, _ := _programFile.Stat()

	fileSize := fileInfo.Size()

	data := make([]byte, fileSize)

	//m := make(map[string]interface{})

	p._program._theProgramJsonMap = make(map[string]interface{})

	_, _ = _programFile.Read(data)

	text := string(data)

	err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)

	if err != nil {
		if len(text) > 0 {

			saveErrorFile(p._programFile, text)
			fmt.Println("124135412362346234623")
			p._program._theProgramJsonMap = DefaultTemplate()

			p._debugHandler.printError("-->readJson: Error while reading the JsonProgram: writing Default\n-->readJson: savefile created!\n")

			p.write()
		}

	} else {
		p._programFile = filename

		p._program._theProgramJson = text

		p._program._theProgramJsonMap = make(map[string]interface{})

		err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)

		if err != nil {
			fmt.Println("HGAGHAHAGAGAOhivOVZ")
			p._program._theProgramJsonMap = DefaultTemplate()

			p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Error Reading: writing Default\n")

			p.write()
		}
	}
	_programFile.Close()

	return text
}

// write the program to the provided configuration file
func (p *program) write() {

	b, _ := json.MarshalIndent(p._program._theProgramJsonMap, " ", "   ")

	p._program._theProgramJson = string(b)

	_programFile, _ := os.OpenFile(p._programFile, os.O_RDWR|os.O_TRUNC, 0666)

	writer := bufio.NewWriter(_programFile)

	writer.WriteString(p._program._theProgramJson)

	writer.Flush()

	_programFile.Close()
}

// check the program for mutable eg immutable options in the top level of the tree
func (p *program) check() {

	//check for required defaults, that are required
	allArgs := p._program._theProgramJsonMap["arguments"]

	if allArgs == nil {

		allArgs = []map[string]interface{}{}
	}

	argArr0, err := allArgs.([]interface{})
	argArr1, err := allArgs.([]map[string]interface{})

	//TODO Deeper parsing, here problems with type casing in go and maps

	hasInteractiveOption := false
	hasLoggingOption := false
	hasHistoryOption := false
	hasHelpOption := false
	hasHistoryFileOption := false
	hasConfigurationFileOption := false
	hasVerbosityOption := false

	if len(argArr0) > len(argArr1) {

		for index, v := range argArr0 {

			if m, ok := v.(map[string]interface{}); ok {

				argArr0[index] = p.checkTopLevelArg(m)

				if argArr0[index].(map[string]interface{})[TYPESTRING] == OPTIONSTRING {

					switch argArr0[index].(map[string]interface{})[LONGFLAGSTRING].(string) {
					case _defaultInteractiveOption:
						{
							hasInteractiveOption = true
						}
					case _defaultLoggingOption:
						{
							hasLoggingOption = true
						}
					case _defaultHistoryOption:
						{
							hasHistoryOption = true
						}
					case _defaultHelpOption:
						{
							hasHelpOption = true
						}
					case _ddefaultHistoryFileOption:
						{
							hasHistoryFileOption = true
						}
					case _defaultConfigurationFileOption:
						{
							hasConfigurationFileOption = true
						}
					case _defaultVerbosityOption:
						{
							hasVerbosityOption = true
						}
					}
				}
			}
		}

		if !hasInteractiveOption {
			argArr0 = append(argArr0, theInteractiveOption)
		}
		if !hasLoggingOption {
			argArr0 = append(argArr0, theLoggingOption)
		}
		if !hasHistoryOption {
			argArr0 = append(argArr0, theHistoryOption)
		}
		if !hasHelpOption {
			argArr0 = append(argArr0, theHelpOption)
		}
		if !hasHistoryFileOption {
			argArr0 = append(argArr0, theHistoryFileOption)
		}
		if !hasConfigurationFileOption {
			argArr0 = append(argArr0, theConfigFileOption)
		}
		if !hasVerbosityOption {
			argArr0 = append(argArr0, theVerbosityOption)
		}

		allArgs = argArr0

		//p.checkTopLevelArgs(argArr0)
	} else {

		for index, m := range argArr1 {
			argArr0[index] = p.checkTopLevelArg(m)
			if argArr0[index].(map[string]interface{})[TYPESTRING] == OPTIONSTRING {

				switch argArr0[index].(map[string]interface{})[LONGFLAGSTRING].(string) {
				case _defaultInteractiveOption:
					{
						hasInteractiveOption = true
					}
				case _defaultLoggingOption:
					{
						hasLoggingOption = true
					}
				case _defaultHistoryOption:
					{
						hasHistoryOption = true
					}
				case _defaultHelpOption:
					{
						hasHelpOption = true
					}
				case _ddefaultHistoryFileOption:
					{
						hasHistoryFileOption = true
					}
				case _defaultConfigurationFileOption:
					{
						hasConfigurationFileOption = true
					}
				case _defaultVerbosityOption:
					{
						hasVerbosityOption = true
					}

				}
			}

		}

		if !hasInteractiveOption {
			argArr1 = append(argArr1, theInteractiveOption)
		}
		if !hasLoggingOption {
			argArr1 = append(argArr1, theLoggingOption)
		}
		if !hasHistoryOption {
			argArr1 = append(argArr1, theHistoryOption)
		}
		if !hasHelpOption {
			argArr1 = append(argArr1, theHelpOption)
		}
		if !hasHistoryFileOption {
			argArr1 = append(argArr1, theHistoryFileOption)
		}
		if !hasConfigurationFileOption {
			argArr1 = append(argArr1, theConfigFileOption)
		}
		if !hasVerbosityOption {
			argArr1 = append(argArr1, theVerbosityOption)
		}

		allArgs = argArr1
	}

	if err {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "Could not read Arguments, writing defaults")
	}

	p._program._theProgramJsonMap["arguments"] = allArgs

	p.write()
}

func saveErrorFile(fileName string, content string) {

	yr, month, day := time.Now().Date()
	hr := time.Now().Hour()
	min := time.Now().Minute()

	theTime := strconv.Itoa(yr) + "_" + strconv.Itoa(int(month)) + "_" + strconv.Itoa(day) + "_" + strconv.Itoa(hr) + "_" + strconv.Itoa(min) + "_"

	fmt.Println("THETIME", strconv.Itoa(yr))

	file, _ := os.OpenFile(theTime+"errsave_"+fileName, os.O_WRONLY|os.O_CREATE, 0666)

	writer := bufio.NewWriter(file)

	writer.WriteString(content)

	writer.Flush()

	file.Close()
}

func (p *program) checkTopLevelArg(arg map[string]interface{}) map[string]interface{} {

	if arg[TYPESTRING] == OPTIONSTRING {

		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultInteractiveOption) == 0 {
			//hasInteractiveOption = true
			//just in case replacement
			if !theInteractiveOption[MUTEABLESTRING].(bool) {
				arg = theInteractiveOption
			}
		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultLoggingOption) == 0 {

			//hasLoggingOption = true
			if !theLoggingOption[MUTEABLESTRING].(bool) {
				arg = theLoggingOption
			}
		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultHistoryOption) == 0 {

			//hasHistoryOption = true
			if !theHistoryOption[MUTEABLESTRING].(bool) {
				arg = theHistoryOption
			}

		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultHelpOption) == 0 {

			//hasHelpOption = true
			if !theHelpOption[MUTEABLESTRING].(bool) {
				arg = theHelpOption
			}
		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _ddefaultHistoryFileOption) == 0 {

			//hasHistoryFileOption = true
			if !theHistoryFileOption[MUTEABLESTRING].(bool) {
				arg = theHistoryFileOption
			}
		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultConfigurationFileOption) == 0 {

			//just replace the option anyways, so we dont have to do recursive or deep layer checks
			//hasConfigurationFileOption = true
			if !theConfigFileOption[MUTEABLESTRING].(bool) {
				arg = theConfigFileOption
			}
		}
		if strings.Compare(arg[LONGFLAGSTRING].(string), _defaultVerbosityOption) == 0 {

			//hasVerbosityOption = true
			if !theVerbosityOption[MUTEABLESTRING].(bool) {
				arg = theVerbosityOption
			}
		}
	}

	return arg
}
