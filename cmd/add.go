package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	addPrintCommandEnvVar = "SPEEDDIAL_ADD_PRINT_COMMAND"
)

var (
	addCmd = &cobra.Command{
		Use:                "add",
		Short:              "Add a new command to speeddial",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,

		Run: runAdd,
	}
)

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
