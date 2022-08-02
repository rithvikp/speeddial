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
		Short: "Shell commands at your fingertips",
		Long:  `After starting this command, type and use the arrow keys to search for the entry you desire. Press "enter" to select the entry: it will be loaded into the subsequent terminal prompt.`,

		Run: run,
	}

	rootRegexArg bool
)

func init() {
	rootCmd.AddCommand(addCmd, initCmd, rmCmd)
	rootCmd.Flags().BoolVarP(&rootRegexArg, "regex", "r", false, "Use regex instead of fuzzy search")
}

// Text output is printed to stderr instead of stdout as what is sent to stderr is printed right
// away whereas what is sent to stdout is buffered and then shown in the next prompt by the shell
// wrapper.

// Execute starts the program.
func Execute() error {
	// If not running the init command, the shell wrapper should be used
	if os.Getenv(initializedEnvVar) == "" && (len(os.Args) < 2 || os.Args[1] != "init") {
		fmt.Fprintln(os.Stderr, "Please use the shell wrapper to call speeddial. Check `speeddial init --help` for more information")
	}
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
	fmt.Println(search(c, rootRegexArg).Invocation)
}

func search(c *state.Container, useRegex bool) *state.Command {
	searcher := term.QueryableList[*state.Command](c.Searcher(useRegex))
	command, err := term.List(searcher, maxDisplayedSearchResults, true)
	if err == term.ErrUserQuit {
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to select a new command: %v\n", err)
		os.Exit(1)
	}

	return command
}
