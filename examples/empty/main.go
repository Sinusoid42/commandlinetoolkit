package main

import (
	"commandlinetoolkit"
	"os"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	cmdline.Parse(os.Args)

	cmdline.ReadJSON("config.json")

	cmdline.Wait()

}
