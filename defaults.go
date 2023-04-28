package commandlinetoolkit

import "fmt"

// when checking a given command|option template, we checkInputProgram in the program, if any of the required or optional commands are present
// if not, they will be overwritten in the config file
// if optional and present, they will not be overwritten

const (
	_defaultInteractiveOption       string = "interactive"
	_defaultLoggingOption                  = "logging"
	_defaultHistoryOption                  = "history"
	_defaultHelpOption                     = "help"
	_ddefaultHistoryFileOption             = "historyfile"
	_defaultConfigurationFileOption        = "config"
	_defaultVerbosityOption                = "verbose"
)

const (
	_defaultVerbosityCommand   string = "verbose"
	_defaultHistoryCommand            = "history"
	_defaultHistoryFileCommand        = "historyfile"
	_defaultShellCommand              = "shell"
	_defaultExitCommand               = "exit"
	_defaultLOGGINGCommand            = "logging"
	_defaultConfigFileCommand         = "config"
	_defaultBootOnlyCommand           = "bootonly"
)

type RUNPROPERTY int32

// runs the command in verbose
const VERBOSITYPROPERTY RUNPROPERTY = 0b0000000000000000001
const EXITPROPERTY RUNPROPERTY = 0b00000000000000000010
const BOOTONLYPROPERTY RUNPROPERTY = 0b000000000000000000100

type RUNCOMMAND int32

const RUNINTERACTIVE RUNCOMMAND = 0b0000000000000000001
const RUNLOGGING RUNCOMMAND = 0b0000000000000000001
const RUNHISTORY RUNCOMMAND = 0b0000000000000000001
const RUNHISTORYFILE RUNCOMMAND = 0b0000000000000000001
const RUNHELP RUNCOMMAND = 0b0000000000000000001
const RUNCONFIG RUNCOMMAND = 0b0000000000000000001

// turns the entire shell verbose
const RUNVERBOSE RUNCOMMAND = 0b0000000000000000001

// store the default arguments in global scope buffer, when rereading or rebuilding the program during runtime
var (
	theInteractiveOption = defaultInteractiveOption()
	theLoggingOption     = defaultLoggingOption()
	theHistoryOption     = defaultHistoryOption()
	theHistoryFileOption = defaultHistoryFileOption()
	theHelpOption        = defaultHelpOption()
	theConfigFileOption  = defaultConfigurationFileOption()
	theVerbosityOption   = defaultVerbosityOption()
)

// the defaultInteractiveOption
// run the shell when --interactive is given to the shell
// accept the "run" command and perfom the shell callback
func defaultInteractiveOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "interactive"
	m[HELPKEY] = ""
	m[MUTEABLEKEY] = false
	m[SHORTHELPKEY] = "Interactive shell mode"
	m[RUNKEY] = "shell"
	m[DATATYPEKEY] = "bool"

	return m
}

func defaultLoggingOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "logging"
	m[MUTEABLEKEY] = false
	m[HELPKEY] = "In case we want to create a history file in which we store all previously executed commands"
	m[RUNKEY] = "logging"
	m[DATATYPEKEY] = "bool"

	return m
}

func defaultHistoryOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "history"
	m[MUTEABLEKEY] = false
	m[HELPKEY] = "Use and allow a current history when in interactive mode. Commands need to be rerunable"
	m[RUNKEY] = "history"
	m[DATATYPEKEY] = "bool"

	return m
}

func defaultHelpOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = HELPKEY
	m[SHORTFLAGKEY] = "h"
	m[MUTEABLEKEY] = true
	m[HELPKEY] = "This is the default help message for the command line demo Execting the example binary, with ./mybinary --help this callback is executed. Can be overwritten"
	m[SHORTHELPKEY] = "Only the short Help Menu.Use '--help' for more info."
	m[RUNKEY] = "exit"

	return m
}

func defaultHistoryFileOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "historyfile"
	m[MUTEABLEKEY] = false
	m[HELPKEY] = "Use and allow a historyfile for multiple executions of the shell, Commands can be reentered etc.."
	m[RUNKEY] = "historyfile"
	m[DATATYPEKEY] = "bool"

	return m
}

func defaultConfigurationFileOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "config"
	m[MUTEABLEKEY] = false
	m[HELPKEY] = "Specifies a configuration file, from which commands will be parsed, can only be executed on booting the application"
	m[RUNKEY] = "config, bootonly"
	m[DATATYPEKEY] = "file"
	m[ARGUMENTSKEY] = []map[string]interface{}{}

	param := make(map[string]interface{})

	param[TYPEKEY] = "PARAM"
	param[DATATYPEKEY] = "file"
	param[REQUIREDKEY] = true

	m[ARGUMENTSKEY] = append(m[ARGUMENTSKEY].([]map[string]interface{}), param)

	return m
}

func defaultVerbosityOption() map[string]interface{} {
	m := make(map[string]interface{})

	m[TYPEKEY] = "OPTION"
	m[LONGFLAGKEY] = "verbose"
	m[MUTEABLEKEY] = true
	m[HELPKEY] = "Run the shell or the entire program in verbose mode. If a number is provided, different verbosities will be used. Can be overwritten."
	m[RUNKEY] = "verbose"
	m[DATATYPEKEY] = "bool"
	return m
}

func isLibCommand(str string) bool {

	switch str {
	case _defaultVerbosityCommand:
		{
			return true
		}
	case _defaultHistoryCommand:
		{
			return true
		}
	case _defaultHistoryFileCommand:
		{
			return true
		}
	case _defaultShellCommand:
		{
			return true
		}
	case _defaultExitCommand:
		{
			return true
		}
	case _defaultLOGGINGCommand:
		{
			return true
		}
	case _defaultConfigFileCommand:
		{
			return true
		}
	case _defaultBootOnlyCommand:
		{
			return true
		}
	}
	return false
}

func getRunCommand(str string) func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

	switch str {
	case _defaultVerbosityCommand:
		{
			return func(parameters []*Argument, opargumentstions []*Argument, cmdline *CommandLine) CLICODE {
				if len(parameters) == 0 {
					//full
					cmdline.SetVerbosity(-1)
					return CLI_TRUE
				}
				return CLI_FALSE
			}
		}

	case _defaultHistoryCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {
				if len(parameters) == 0 {

					cmdline.Set(HISTORY, CLI_TRUE)
					return CLI_TRUE
				}
				return CLI_FALSE
			}
		}
	case _defaultHistoryFileCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {
				if len(parameters) == 0 {
					cmdline.Set(HISTORYFILE, CLI_TRUE)
					return CLI_TRUE
				}
				return CLI_FALSE
			}
		}
	case _defaultShellCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {
				if len(parameters) == 0 {
					cmdline.Set(SHELL, CLI_TRUE)
					return CLI_TRUE
				}

				if parameters[0].GetValue() == true {
					cmdline.Set(SHELL, CLI_TRUE)
					return CLI_TRUE
				}

				return CLI_FALSE
			}
		}
	case _defaultExitCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

				cmdline.Exit()

				//TODO

				return CLI_TRUE
			}
		}
	case _defaultLOGGINGCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

				if len(parameters) == 0 {
					cmdline.Set(LOGGING, CLI_TRUE)
					return CLI_TRUE
				}
				return CLI_FALSE
			}
		}
	case _defaultConfigFileCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

				if !cmdline._booted {
					return CLI_FALSE
				}
				fmt.Println(parameters)

				if len(parameters) == 1 {
					fmt.Println("\n" + string(COLOR_CYAN) + "Reading config file..." + string(COLOR_RESET))
					cmdline.ReadJSON(parameters[0].GetValue().([]string)[0])
				}
				return CLI_FALSE
			}
		}
	case _defaultBootOnlyCommand:
		{
			return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {
				if !cmdline._booted {
					return CLI_FALSE
				}
				return CLI_FALSE
			}
		}
	}
	return nil
}

func getRunCommands(cmds []string) func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

	runProperties := retrieveRunProperties(cmds)
	runCmd := retrieveRunCommand(cmds)

	return func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE {

		if runProperties&BOOTONLYPROPERTY > 0 {
			if cmdline._booted {
				return CLI_FALSE
			}
		}

		if runProperties&VERBOSITYPROPERTY > 0 {
			if len(parameters) == 0 {
				//full

				cmdline.SetVerbosity(-1)
			}
		}

		//TODO
		if runCmd&RUNCONFIG > 0 {
			if len(parameters) == 1 {
				fmt.Println("\n" + string(COLOR_CYAN) + "Reading config file..." + string(COLOR_RESET))
				fmt.Println(parameters[0].GetValue().([]string)[0])
				cmdline.ReadJSON(parameters[0].GetValue().([]string)[0])
			}
		}
		if runProperties&VERBOSITYPROPERTY > 0 {
			if len(parameters) == 0 {

				cmdline.SetVerbosity(0)
			}
		}
		if runProperties&EXITPROPERTY > 0 {
			cmdline.Exit()
		}

		return CLI_TRUE
	}
}

func retrieveRunProperties(cmds []string) RUNPROPERTY {
	property := RUNPROPERTY(0)
	for _, s := range cmds {

		if s == _defaultBootOnlyCommand {
			property |= BOOTONLYPROPERTY
		}
		if s == _defaultVerbosityCommand {
			property |= VERBOSITYPROPERTY
		}
		if s == _defaultExitCommand {
			property |= EXITPROPERTY
		}
	}
	return property
}

func retrieveRunCommand(cmds []string) RUNCOMMAND {
	r := RUNCOMMAND(0)
	for _, s := range cmds {
		if s == _defaultInteractiveOption {
			r |= RUNINTERACTIVE
		}
		if s == _defaultShellCommand {
			r |= RUNINTERACTIVE
		}
		if s == _defaultHistoryCommand {
			r |= RUNHISTORY
		}
		if s == _defaultHistoryFileCommand {
			r |= RUNHISTORYFILE
		}
		if s == _defaultLOGGINGCommand {
			r |= RUNLOGGING
		}
		if s == _defaultConfigFileCommand {
			r |= RUNCONFIG
		}
	}
	return r
}
