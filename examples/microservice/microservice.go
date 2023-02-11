package main

import (
	"commandlinetoolkit"
	"os"
)

func main() {
	
	cmdLine := commandlinetoolkit.NewCommandLine()
	
	cmdLine.ReadJSON("commandlineconfig.json")
	
	e := cmdLine.Parse(os.Args)
	
	if e == commandlinetoolkit.CLI_SUCCESS {
	
	}
	
	cmdLine.Wait()
}
