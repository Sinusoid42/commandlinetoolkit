package main

import (
	"commandlinetoolkit"
	"fmt"
	"os"
)

func main() {
	
	cmdLine := commandlinetoolkit.NewCommandLine()
	
	cmdLine.ReadJSON("commandlineconfig.json")
	
	fmt.Println(cmdLine)
	
	e := cmdLine.Parse(os.Args)
	
	if e == commandlinetoolkit.CLI_SUCCESS {
	
	}
	
	cmdLine.Wait()
}
