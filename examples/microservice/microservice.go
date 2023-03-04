package main

import (
	"fmt"
	"github.com/Sinusoid42/commandlinetoolkit"
)

func main() {

	cmdLine := commandlinetoolkit.NewCommandLine()

	cmdLine.ReadJSON("commandlineconfig.json")

	//e := cmdLine.Parse(os.Args)

	//if e == commandlinetoolkit.CLI_SUCCESS {

	//}

	cmdLine.Set(commandlinetoolkit.SHELL|
		commandlinetoolkit.HISTORY|
		commandlinetoolkit.HISTORYFILE|
		commandlinetoolkit.PREDICTIONS|
		commandlinetoolkit.SUGGESTIONS,
		commandlinetoolkit.CLI_TRUE)

	fmt.Println("DO STH.... ")

	cmdLine.Wait()
}
