package commandlinetoolkit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	_programTitle string
	_styleTitle   bool
	_titleLength  int
	_programDepth int
	_program      *commandlinetemplate
	_verbose      CLICODE
	_debugHandler *debugHandler
}

/*
****************************************************************************************************************************************

Builds a new Program
*/
func newprogram(filename string) *program {

	p := &program{
		_programName:  "Command Line: " + VERSION,
		_programFile:  filename,
		_program:      DefaultCommandLineTemplate(),
		_debugHandler: newDebugHandler(),
		_titleLength:  36,
	}

	p._debugHandler._verbose = 0

	return p
}

/*
****************************************************************************************************************************************

Read a json cli configuration file
*/
func (p *program) readJsonProgram(filename string) string {

	//checkInputProgram for file forst
	if _, err := os.OpenFile(filename, os.O_RDONLY, 0666); err != nil {

		fmt.Println(filename)
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

	p._program._theProgramJsonMap = make(map[string]interface{})

	_, _ = _programFile.Read(data)

	text := string(data)

	err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)

	if err != nil {
		if len(text) > 0 {

			saveErrorFile(p._programFile, text)
			p._program._theProgramJsonMap = DefaultTemplate()
			p._debugHandler.printError("-->readJson: Error while reading the JsonProgram: writing Default\n-->readJson: savefile created!\n")
			p.writeJsonProgram()
		}

	} else {
		p._programFile = filename

		p._program._theProgramJson = text

		p._program._theProgramJsonMap = make(map[string]interface{})

		err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)

		if err != nil {
			p._program._theProgramJsonMap = DefaultTemplate()

			p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Error Reading: writing Default\n")

			p.writeJsonProgram()
		}
	}
	_programFile.Close()

	return text
}

func (p *program) storeJsonProfile(prstree *parsetree) CLICODE {

	return CLI_SUCCESS
}

/*
****************************************************************************************************************************************

Write the json to disc
*/
func (p *program) writeJsonProgram() {

	b, _ := json.MarshalIndent(p._program._theProgramJsonMap, " ", "   ")

	p._program._theProgramJson = string(b)

	_programFile, _ := os.OpenFile(p._programFile, os.O_RDWR|os.O_TRUNC, 0666)

	writer := bufio.NewWriter(_programFile)

	writer.WriteString(p._program._theProgramJson)

	writer.Flush()

	_programFile.Close()
}

/*
****************************************************************************************************************************************

Generate a Programtitle that is colorful and looks cool based on given paramters of the program
*/

func (p *program) genTitle() string {

	name := p._programName

	if len(p._programName) > 0 {
		name = p._programName
	}

	p._titleLength = 38

	l := p._titleLength
	if len(name) < p._titleLength {
		l = len(name)
	}
	borderColor := string(COLOR_GRAY_I)
	titleColor := "\033[0;0;0m\033[49;96;22m"
	border := "" + borderColor + strings.Repeat("*", p._titleLength+8) + "\n"

	title := border
	title += "*" + strings.Repeat(" ", p._titleLength+6) + "*\n"
	if len(name) < p._titleLength {

		title += "*   " + strings.Repeat(" ", (p._titleLength-l)/2) + titleColor + name[0:l] + string(COLOR_RESET) + borderColor
	} else {
		title += "*   " + titleColor + name[0:l] + string(COLOR_RESET) + borderColor
	}
	if len(name)%2 == 1 && len(name) < p._titleLength && p._titleLength%2 == 0 {
		title += " "
	}
	if len(name) < p._titleLength {
		title += "   " + strings.Repeat(" ", (p._titleLength-l)/2) + "*\n"
	} else {
		title += "   " + "*\n"
	}
	title += "*" + strings.Repeat(" ", p._titleLength+6) + "*\n"
	title += border

	title += string(COLOR_RESET)

	if p._styleTitle {

		if len(p._programName) < 12 {

			l_0 := ""
			l_1 := ""
			l_2 := ""
			l_3 := ""
			l_4 := ""

			title = ""
			pn := strings.ToLower(p._programName)
			for _, v := range pn {

				a := int(v) - 97

				//color := Color()

				//#anilop
				color := GenColor("1", INTENSITY_COLORTYPE, BLUE_COLOR)

				if a >= 0 && a <= 26 {

					l_0 += string(color) + " " + l0[(a*13+1):(a*13+13)] + string(COLOR_RESET) + ""
					//l_0 += string(color) + "" + "             " + string(COLOR_RESET) + ""
					l_1 += string(color) + "" + l1[(a*13):(a*13+13)] + string(COLOR_RESET) + ""
					l_2 += string(color) + "" + l2[(a*13):(a*13+13)] + string(COLOR_RESET) + ""
					l_3 += string(color) + "" + l3[(a*13):(a*13+13)] + string(COLOR_RESET) + ""
					l_4 += string(color) + "" + l4[(a*13):(a*13+13)] + string(COLOR_RESET) + ""

				} else {
					l_1 += string(color) + "" + "         " + string(COLOR_RESET) + ""
					l_0 += string(color) + "" + "         " + string(COLOR_RESET) + ""
					l_2 += string(color) + "" + "         " + string(COLOR_RESET) + ""
					l_3 += string(color) + "" + "         " + string(COLOR_RESET) + ""
					l_4 += string(color) + "" + "         " + string(COLOR_RESET) + ""

				}
			}

			title = l_0 + "\n" + l_1 + "\n" + l_2 + "\n" + l_3 + "\n" + l_4 + "\n"
		}

	}

	p._programTitle = string(title)

	return string(title)
}

const l_ string = `:::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: :::::::::::: `
const l0 string = ` ______       ______       ______       _____        ______       ______       ______       __  __         __           __         __  __       __           __    __     __   __      ______       ______       ______       ______       ______       ______       __  __       __   __      __     __    __  __       __  __       ______`
const l0_ string = `______       ______       ______       _____        ______       ______       ______       __  __         __           __         __  __       __           __    __     __   __      ______       ______       ______       ______       ______       ______       __  __       __   __      __     __    __  __       __  __       ______`
const l1 string = `/\  __ \     /\  == \     /\  ___\     /\  __-.     /\  ___\     /\  ___\     /\  ___\     /\ \_\ \       /\ \         /\ \       /\ \/ /      /\ \         /\ "-./  \   /\ "-.\ \    /\  __ \     /\  == \     /\  __ \     /\  == \     /\  ___\     /\__  _\     /\ \/\ \     /\ \ / /     /\ \  _ \ \  /\_\_\_\     /\ \_\ \     /\___  \`
const l2 string = `\ \  __ \    \ \  __<     \ \ \____    \ \ \/\ \    \ \  __\     \ \  __\     \ \ \__ \    \ \ \_\ \      \ \ \        \_\ \      \ \  _"-.    \ \ \____    \ \ \-./\ \  \ \ \-.  \   \ \ \/\ \    \ \  _-/     \ \ \/\_\    \ \  __<     \ \___  \    \/_/\ \/     \ \ \_\ \    \ \ \'/      \ \ \/ ".\ \ \/_/\_\/_    \ \____ \    \/_/  /__`
const l3 string = ` \ \_\ \_\    \ \_____\    \ \_____\    \ \____-     \ \_____\    \ \_\        \ \_____\    \ \_\ \_\      \ \_\     /\_____\      \ \_\ \_\    \ \_____\    \ \_\ \ \_\  \ \_\\"\_\   \ \_____\    \ \_\        \ \___\_\    \ \_\ \_\    \/\_____\      \ \_\      \ \_____\    \ \__|       \ \__/".~\_\  /\_\/\_\    \/\_____\     /\_____\`
const l4 string = `  \/_/\/_/     \/_____/     \/_____/     \/____/      \/_____/     \/_/         \/_____/     \/_/\/_/       \/_/     \/_____/       \/_/\/_/     \/_____/     \/_/  \/_/   \/_/ \/_/    \/_____/     \/_/         \/___/_/     \/_/ /_/     \/_____/       \/_/       \/_____/     \/_/         \/_/   \/_/  \/_/\/_/     \/_____/     \/_____/ `
