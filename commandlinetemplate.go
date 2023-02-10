package commandlinetoolkit

type commandlinetemplate struct {
	_theProgramJson string
}

func DefaultCommandLineTemplate() *commandlinetemplate {
	return &commandlinetemplate{_theProgramJson: "{" +
		"\"version\" : \"0.0.4\"" +
		"\"author\" : \"benwinter\"" +
		"\"description\" : \"This is a default command line demo example, how i want a  json to be parsed into a callback tree with contitional execution, runtime callback optioons and exeption handling\"" +
		"\"arguments\": [\n   " +
		"{\n      " +
		"\"type\" : \"OPTION\",\n      " +
		"\"flag\" : \"help\",\n      " +
		"\"sflag\" : \"h\",\n      " +
		"\"help\" : \"This is the default help message for the command line demo\n\nExecting the example binary, with ./mybinary --help this callback is executed, if custom code is given, the code can just be added by using the api for this software\n\",\n" +
		"\"shelp\" : \"Only the short Help Menu.\nUse '--help' for more info.\" ,\n" +
		"\"method\" : \"exit()\"\n    " +
		"},\n    " +
		"{\n      " +
		"\"type\" : \"OPTION\",\n      \"flag\" : \"interactive\",\n      \"sflag\" : \"i\",\n      \"help\" : \"Interactive shell mode\"\n    },\n    {\n      \"type\": \"OPTION\",\n      \"flag\": \"_logging\",\n      \"help\": \"In case we want to create a history file in which we store all previously executed commands\"\n    }\n  ]" +
		"}"}
}
