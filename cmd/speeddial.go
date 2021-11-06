package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rithvikp/speeddial/state"
	"github.com/rithvikp/speeddial/term"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                "speeddial",
		Short:              "Commands at your fingertips",
		Args:               cobra.MinimumNArgs(0),
		DisableFlagParsing: true,

		Run: run,
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
		fmt.Printf("Unable to initialize speeddial state: %v\n", err)
		os.Exit(1)
	}

	return c
}

func cleanup(c *state.Container) {
	c.Dump()
}

func run(cmd *cobra.Command, args []string) {
	c := setup()
	defer cleanup(c)

	term.List(c)
}

func runAdd(cmd *cobra.Command, args []string) {
	c := setup()
	defer cleanup(c)

	// TODO: Validate args accordingly
	command := strings.Join(args, " ")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please input a description if desired: ")
	scanner.Scan()
	desc := scanner.Text()

	err := c.NewCommand(command, desc)
	if err != nil {
		fmt.Printf("Unable to add the new command: %v\n", err)
	}
}
