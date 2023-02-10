package commandlinetoolkit

type program struct {
	_programName string
	
	_program *commandlinetemplate
}

func newprogram() *program {
	
	return &program{
		_programName: "Command Line: " + VERSION,
		_program:     DefaultCommandLineTemplate(),
	}
	
}
