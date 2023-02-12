# Commandline Toolkit

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



