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

		Run: run,
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(rmCmd)
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

func dump(c *state.Container) {
	c.Dump()
}

func run(cmd *cobra.Command, args []string) {
	c := setup()
	fmt.Println(search(c).Invocation)
}

func search(c *state.Container) *state.Command {
	rawCommand, err := term.List(c.Searcher(), maxDisplayedSearchResults, true)
	if err == term.ErrUserQuit {
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to select a new command: %v", err)
		os.Exit(1)
	}

	command, ok := rawCommand.(*state.Command)
	if !ok {
		fmt.Fprintf(os.Stderr, "Ran into an unexpected error")
		os.Exit(1)
	}

	return command
}
