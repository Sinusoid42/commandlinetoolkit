package main

import (
	"commandlinetoolkit"
	"os"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	//cmdline.Set(commandlinetoolkit.SHELL, commandlinetoolkit.CLI_TRUE)

	//cmdline.PrintTitle(true)

	cmdline.ReadJSON("config.json")

	cmdline.Parse(os.Args)

	cmdline.Wait()

}
