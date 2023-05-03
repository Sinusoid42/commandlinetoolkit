package main

import (
	"fmt"
	"github.com/Sinusoid42/commandlinetoolkit"
	"os"
)

func main() {

	cmdLine := commandlinetoolkit.NewCommandLine()
	cmdLine.PrintTitle(false)

	cmdLine.ReadJSON("commandlineconfig.json")

	cmdLine.Parse(os.Args)

	cmdLine.ParseTree().Get("port").Value()

	//if e == commandlinetoolkit.CLI_SUCCESS {

	//}

	fmt.Println(cmdLine.ParseTree())

	cmdLine.Set(
		commandlinetoolkit.HISTORY|
			commandlinetoolkit.HISTORYFILE|
			commandlinetoolkit.PREDICTIONS|
			commandlinetoolkit.SUGGESTIONS,
		commandlinetoolkit.CLI_TRUE)

	fmt.Println("> DO STH <")
	for i := 0; i < 3; i++ {
		fmt.Println("Do something? y/n")
		a := cmdLine.YesNoConfirm()
		fmt.Println("DONE > ", a)

	}
	cmdLine.Wait()
}
