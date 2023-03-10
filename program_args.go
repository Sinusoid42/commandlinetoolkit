package commandlinetoolkit

/*
****************************************************************************************************************************************

A commandline program argument
*/
type PROGRAM_ARGUMENT int32

const SHELL PROGRAM_ARGUMENT = 0b00000000001
const HISTORY PROGRAM_ARGUMENT = 0b00000000010
const HISTORYFILE PROGRAM_ARGUMENT = 0b00000000100
const SUGGESTIONS PROGRAM_ARGUMENT = 0b00000001000
const PREDICTIONS PROGRAM_ARGUMENT = 0b00000010000
const LOGGING PROGRAM_ARGUMENT = 0b00000100000

type programArgs struct {
	hasInteractiveOption       bool
	hasLoggingOption           bool
	hasHistoryOption           bool
	hasHelpOption              bool
	hasHistoryFileOption       bool
	hasConfigurationFileOption bool
	hasVerbosityOption         bool
}

/*
****************************************************************************************************************************************

Generate a Programtitle
*/
func (p *programArgs) success() bool {
	return p.hasInteractiveOption &&
		p.hasLoggingOption &&
		p.hasHistoryFileOption &&
		p.hasHistoryOption &&
		p.hasConfigurationFileOption
}
