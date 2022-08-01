package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	zshShell              = "zsh"
	fishShell             = "fish"
	addPrintCommandEnvVar = "SPEEDDIAL_ADD_PRINT_COMMAND"
	initializedEnvVar     = "SPEEDDIAL_INITIALIZED"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Output initialization code for shells",
		Long: `Setup the shell wrapper for speeddial.

For zsh, add the following to your rc file: eval "$(speeddial init zsh)"
For fish, add the following to your rc file: speeddial init fish | source

Check https://github.com/rithvikp/speeddial for more information.`,

		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{zshShell, fishShell},
		Run:       runInit,
	}
)

func runInit(cmd *cobra.Command, args []string) {
	switch args[0] {
	case zshShell:
		fmt.Println(zshInitialization)
	case fishShell:
		fmt.Println(fishInitialization)
	}
}

const (
	zshInitialization = `
spd() {
    if [ "$1" = "add" ] && [ "$#" = 1 ]; then
		SPEEDDIAL_INITIALIZED=1 SPEEDDIAL_ADD_PRINT_COMMAND=1 speeddial add "$(fc -ln -1)"
    elif [ "$1" = "" ]; then
        print -z $(SPEEDDIAL_INITIALIZED=1 speeddial)
    else
        SPEEDDIAL_INITIALIZED=1 speeddial $@
    fi
}`

	fishInitialization = `
function spd
    if test "$argv[1]" = "add"; and test (count $argv) = 1
        SPEEDDIAL_INITIALIZED=1 SPEEDDIAL_ADD_PRINT_COMMAND=1 speeddial add "$history[1]"
    else if test (count $argv) = 0
        commandline (SPEEDDIAL_INITIALIZED=1 speeddial)
    else
        SPEEDDIAL_INITIALIZED=1 speeddial $argv
    end
end`
)
