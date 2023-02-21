package main

import (
	"commandlinetoolkit"
	"os"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	cmdline.ReadJSON("")

	cmdline.Parse(os.Args)

}
