package main

import (
	"fmt"
	"github.com/Sinusoid42/commandlinetoolkit"
	"os"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	//cmdline.Set(commandlinetoolkit.SHELL, commandlinetoolkit.CLI_TRUE)

	//cmdline.PrintTitle(true)

	cmdline.ReadJSON("config.json")
	/*
		program := cmdline.Program()

		fmt.Println(program.String())
	*/

	cmdline.Parse(os.Args)

	parseTree := cmdline.ParseTree()

	myoption := parseTree.Get("port")

	fmt.Println(myoption.Argument().GetValue())

	fmt.Println(myoption.Next().Next().Argument().GetValue())

	cmdline.Wait()

}
