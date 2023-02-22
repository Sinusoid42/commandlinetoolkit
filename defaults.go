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
	m[RUNKEY] = "configFile, bootonly"
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

	return m
}
