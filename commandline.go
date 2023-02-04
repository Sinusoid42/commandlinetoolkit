package commandlinetoolkit

//base file for building a commandline
//here struct definition and runtime for a commandline is given
//for a command line to work we need the possibility of parsing a non cyclic directed, n-dimensional tree from arguments that follow conditionally
//either by setting the paramters required or not,

type CommandLine struct {
	root_argument Argument
	
	size    int32
	methods int32
	options int32
	verbose int32
}
