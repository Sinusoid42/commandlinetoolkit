package commandlinetoolkit

type commandlinetemplate struct {
	_theProgramJson string

	_theProgramJsonMap map[string]interface{}
}

func DefaultCommandLineTemplate() *commandlinetemplate {

	return &commandlinetemplate{_theProgramJson: "",
		_theProgramJsonMap: make(map[string]interface{})}
}

func D() *commandlinetemplate {
	return &commandlinetemplate{
		_theProgramJson:    "",
		_theProgramJsonMap: make(map[string]interface{}),
	}
}

func DefaultTemplate() map[string]interface{} {

	m := make(map[string]interface{})

	m[VERSION] = "0.0.1"
	m[AUTHORSTRING] = "ben"
	m[EXECUTEABLESTRING] = ""
	m[DESCRIPTIONSTRING] = "The description of the application"
	m[MANUALSTRING] = "The Man Page for this application"

	args := []map[string]interface{}{}

	args = append(args, defaultHelpOption())

	args = append(args, defaultInteractiveOption())

	args = append(args, defaultLoggingOption())

	args = append(args, defaultHistoryOption())

	args = append(args, defaultHistoryFileOption())

	args = append(args, defaultVerbosityOption())

	args = append(args, defaultConfigurationFileOption())

	m[ARGUMENTSSTRING] = args

	return m

}
