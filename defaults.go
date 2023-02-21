package commandlinetoolkit

// when checking a given command|option template, we check in the program, if any of the required or optional commands are present
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

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "interactive"
	m[HELPSTRING] = ""
	m[MUTEABLESTRING] = false
	m[SHORTHELPSTRING] = "Interactive shell mode"
	m[RUNSTRING] = "shell"
	m[DATATYPESTRING] = "bool"

	return m
}

func defaultLoggingOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "logging"
	m[MUTEABLESTRING] = false
	m[HELPSTRING] = "In case we want to create a history file in which we store all previously executed commands"
	m[RUNSTRING] = "logging"
	m[DATATYPESTRING] = "bool"

	return m
}

func defaultHistoryOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "history"
	m[MUTEABLESTRING] = false
	m[HELPSTRING] = "Use and allow a current history when in interactive mode. Commands need to be rerunable"
	m[RUNSTRING] = "history"
	m[DATATYPESTRING] = "bool"

	return m
}

func defaultHelpOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = HELPSTRING
	m[SHORTFLAGSTRING] = "h"
	m[MUTEABLESTRING] = true
	m[HELPSTRING] = "This is the default help message for the command line demo Execting the example binary, with ./mybinary --help this callback is executed. Can be overwritten"
	m[SHORTHELPSTRING] = "Only the short Help Menu.Use '--help' for more info."
	m[RUNSTRING] = "exit"

	return m
}

func defaultHistoryFileOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "historyfile"
	m[MUTEABLESTRING] = false
	m[HELPSTRING] = "Use and allow a historyfile for multiple executions of the shell, Commands can be reentered etc.."
	m[RUNSTRING] = "historyfile"
	m[DATATYPESTRING] = "bool"

	return m
}

func defaultConfigurationFileOption() map[string]interface{} {

	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "config"
	m[MUTEABLESTRING] = false
	m[HELPSTRING] = "Specifies a configuration file, from which commands will be parsed, can only be executed on booting the application"
	m[RUNSTRING] = "configFile, bootonly"
	m[DATATYPESTRING] = "file"
	m[ARGUMENTSSTRING] = []map[string]interface{}{}

	param := make(map[string]interface{})

	param[TYPESTRING] = "PARAM"
	param[DATATYPESTRING] = "file"
	param[REQUIREDSTRING] = true

	m[ARGUMENTSSTRING] = append(m[ARGUMENTSSTRING].([]map[string]interface{}), param)

	return m
}

func defaultVerbosityOption() map[string]interface{} {
	m := make(map[string]interface{})

	m[TYPESTRING] = "OPTION"
	m[LONGFLAGSTRING] = "verbose"
	m[MUTEABLESTRING] = true
	m[HELPSTRING] = "Run the shell or the entire program in verbose mode. If a number is provided, different verbosities will be used. Can be overwritten."
	m[RUNSTRING] = "verbose"

	return m
}
