package commandlinetoolkit

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
