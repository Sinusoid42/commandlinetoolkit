# Commandline Toolkit

### v0.3.0

-baseline option parsing is done, works from shell and from commandline
-callbackfunctions + parameterized options available
-todo: parameterized binary execution of 3rd. party binaries for full dynamic runtime possibilities


### v0.2.0

- checking of program file
- parsing into parsetree and building parse api



### v0.1.1 dev version

- jsonfile reading and building of parsetree from json, if no json provided will build default
  - building of default commands, rejecting certain user commands/options as they interfere with the API of this software
  - rejecting:
    - --interactive : Allow the shell to be run, otherwise the commandline interface will only parse ARGS and wont boot into a secondary goroutine
    - --history : Commands can be rerun, arrow keys are captured by the shell
    - --historyfile :  (depends on available commands from the user), here we have a text file build, that contains previously run commands
  - overwriteable commands:
    - --help : builds autohelp information about the entire available parsetree




### v0.1.0 dev version
- now has interactive shell, with colored prefix and seperates debug and gpio
- historyfile can be read, prev. commands are stored and buffered
- CTRL+C prompts for exiting, confirm with y or Y
- TAB Completion is also implemented, single word tabcompletion can be proposed or multiple (TODO)

- TODO
- parsetree implementation and callback API for outside communication
- proper code cleanup (started with seperating shellHandler .. see@ commit 6dddef41..)
- read jsonfile and create default json file for the command tree is required

- output and callback implementation will be callbackdriven, registering API TODO



