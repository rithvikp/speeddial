package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	zshShell = "zsh"
)

var (
	initCmd = &cobra.Command{
		Use:       "init",
		Short:     "Output initialization code for shells",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{zshShell},
		Run:       runInit,
	}
)

func runInit(cmd *cobra.Command, args []string) {
	switch args[0] {
	case zshShell:
		fmt.Println(zshInitialization)
	}
}

const (
	zshInitialization = `
spd() {
    if [ "$1" = "add" ] && [ "$#" = 1 ]; then
        SPEEDDIAL_ADD_PRINT_COMMAND=1 ./speeddial add $(fc -ln -1)
    elif [ "$1" = "" ]; then
        print -z $(./speeddial)
    else
        ./speeddial $@
    fi
}
`
)
