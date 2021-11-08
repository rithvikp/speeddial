package cmd

import (
	"fmt"
	"os"

	"github.com/rithvikp/speeddial/state"
	"github.com/rithvikp/speeddial/term"
	"github.com/spf13/cobra"
)

const (
	maxDisplayedSearchResults = 10
)

var (
	rootCmd = &cobra.Command{
		Use:   "speeddial",
		Short: "Commands at your fingertips",

		Run: runSearch,
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

// Execute starts the program.
func Execute() error {
	return rootCmd.Execute()
}

func setup() *state.Container {
	c, err := state.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize speeddial state: %v\n", err)
		os.Exit(1)
	}

	return c
}

func cleanup(c *state.Container) {
	c.Dump()
}

func runSearch(cmd *cobra.Command, args []string) {
	c := setup()
	defer cleanup(c)

	rawCommand, err := term.List(c.Searcher(), maxDisplayedSearchResults)
	if err == term.ErrUserQuit {
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to select a new command: %v", err)
		os.Exit(1)
	}

	command, ok := rawCommand.(*state.Command)
	if !ok {
		fmt.Fprintf(os.Stderr, "Ran into an unexpected error")
		os.Exit(1)
	}

	fmt.Println(command.Invocation)
}
