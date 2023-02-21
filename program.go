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
	
	m := make(map[string]interface{})
	
	_, _ = _programFile.Read(data)
	
	text := string(data)
	
	err := json.Unmarshal([]byte(text), &m)
	
	if err != nil {
		if len(text) > 0 {
			
			saveErrorFile(p._programFile, text)
			
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
	
	var argArr []interface{}
	
	if len(argArr0) > len(argArr1) {
		argArr = argArr0
	} else {
	
	}
	
	if err {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "Could not read Arguments, writing defaults")
	}
	
	// check factory defaults
	hasInteractiveOption := false
	hasLoggingOption := false
	hasHistoryOption := false
	hasHelpOption := false
	hasHistoryFileOption := false
	hasConfigurationFileOption := false
	hasVerbosityOption := false
	
	//switch over top level options/args
	for index, arg := range argArr {
		
		if arg, ok := arg.(map[string]interface{}); ok {
			
			if arg["type"] == "OPTION" {
				
				if strings.Compare(arg["flag"].(string), _defaultInteractiveOption) == 0 {
					hasInteractiveOption = true
					//just in case replacement
					if !theInteractiveOption["muteable"].(bool) {
						argArr[index] = theInteractiveOption
					}
				}
				if strings.Compare(arg["flag"].(string), _defaultLoggingOption) == 0 {
					
					hasLoggingOption = true
					if !theLoggingOption["muteable"].(bool) {
						argArr[index] = theLoggingOption
					}
				}
				if strings.Compare(arg["flag"].(string), _defaultHistoryOption) == 0 {
					
					hasHistoryOption = true
					if !theHistoryOption["muteable"].(bool) {
						argArr[index] = theHistoryOption
					}
					
				}
				if strings.Compare(arg["flag"].(string), _defaultHelpOption) == 0 {
					
					hasHelpOption = true
					if !theHelpOption["muteable"].(bool) {
						argArr[index] = theHelpOption
					}
				}
				if strings.Compare(arg["flag"].(string), _ddefaultHistoryFileOption) == 0 {
					
					hasHistoryFileOption = true
					if !theHistoryFileOption["muteable"].(bool) {
						argArr[index] = theHistoryFileOption
					}
				}
				if strings.Compare(arg["flag"].(string), _defaultConfigurationFileOption) == 0 {
					
					//just replace the option anyways, so we dont have to do recursive or deep layer checks
					hasConfigurationFileOption = true
					if !theConfigFileOption["muteable"].(bool) {
						argArr[index] = theConfigFileOption
					}
				}
				if strings.Compare(arg["flag"].(string), _defaultVerbosityOption) == 0 {
					
					hasVerbosityOption = true
					if !theVerbosityOption["muteable"].(bool) {
						argArr[index] = theVerbosityOption
					}
				}
			}
		}
	}
	
	//handle nonmuteables
	
	if !hasInteractiveOption {
		argArr = append(argArr, theInteractiveOption)
	}
	if !hasLoggingOption {
		argArr = append(argArr, theLoggingOption)
	}
	if !hasHistoryOption {
		argArr = append(argArr, theHistoryOption)
	}
	if !hasHelpOption {
		argArr = append(argArr, theHelpOption)
	}
	if !hasHistoryFileOption {
		argArr = append(argArr, theHistoryFileOption)
	}
	if !hasConfigurationFileOption {
		argArr = append(argArr, theConfigFileOption)
	}
	if !hasVerbosityOption {
		argArr = append(argArr, theVerbosityOption)
	}
	
	p._program._theProgramJsonMap["arguments"] = argArr
	
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
