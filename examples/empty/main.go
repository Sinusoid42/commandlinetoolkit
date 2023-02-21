package main

import (
	"commandlinetoolkit"
)

func main() {

	cmdline := commandlinetoolkit.NewCommandLine()

	//cmdline.Parse(os.Args)

	cmdline.ReadJSON("config.json")

	cmdline.Wait()

}
