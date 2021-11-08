package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rithvikp/speeddial/state"
	"github.com/rithvikp/speeddial/term"
	"github.com/spf13/cobra"
)

const (
	addPrintCommandEnvVar = "SPEEDDIAL_ADD_PRINT_COMMAND"

	maxDisplayedSearchResults = 10
)

var (
	rootCmd = &cobra.Command{
		Use:                "speeddial",
		Short:              "Commands at your fingertips",
		Args:               cobra.MinimumNArgs(0),
		DisableFlagParsing: true,

		Run: runSearch,
	}

	addCmd = &cobra.Command{
		Use:                "add",
		Short:              "Add a new command to speeddial",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,

		Run: runAdd,
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
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

func runAdd(cmd *cobra.Command, args []string) {
	c := setup()
	defer cleanup(c)

	// TODO: Validate args accordingly
	command := strings.Join(args, " ")

	printCommand := os.Getenv(addPrintCommandEnvVar) != ""
	if printCommand {
		fmt.Printf("Adding command: %s\n", pterm.Bold.Sprint(command))
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please input a description if desired: ")
	scanner.Scan()
	desc := scanner.Text()

	err := c.NewCommand(command, desc)
	if err != nil {
		fmt.Printf("Unable to add the new command: %v\n", err)
	}
}
