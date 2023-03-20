package main

import (
	"fmt"
	"github.com/Sinusoid42/commandlinetoolkit"
	"os"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	//cmdline.Set(commandlinetoolkit.SHELL, commandlinetoolkit.CLI_TRUE)

	//	cmdline.PrintTitle(true)

	//cmdline.StyleTitle(true)

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

	test := map[string]interface{}{}
	test[commandlinetoolkit.TYPEKEY] = commandlinetoolkit.COMMANDSTRING
	test[commandlinetoolkit.RUNKEY] = "test"
	test[commandlinetoolkit.LONGFLAGKEY] = "test"
	arg, _ := commandlinetoolkit.NewArgument(test)

	cmdline.AddArgument(arg, func(parameters []*commandlinetoolkit.Argument, arguments []*commandlinetoolkit.Argument, cmdline *commandlinetoolkit.CommandLine) commandlinetoolkit.CLICODE {

		if len(parameters) > 0 {
			fmt.Println(parameters)
		}

		if len(arguments) > 0 {
			fmt.Println(arguments)
			fmt.Println(arguments[0].GetValue())
		}

		fmt.Println("\nDo Something")

		fmt.Println(cmdline.Program())

		return commandlinetoolkit.CLI_SUCCESS
	})

	cmdline.Wait()

}
