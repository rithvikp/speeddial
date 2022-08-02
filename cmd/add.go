package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new command to speeddial",
		Args:  cobra.ExactArgs(1),

		Run: runAdd,
	}
)

func runAdd(cmd *cobra.Command, args []string) {
	c := setup()
	defer dump(c)

	command := args[0]

	printCommand := os.Getenv(addPrintCommandEnvVar) != ""
	if printCommand {
		fmt.Fprintf(os.Stderr, "Adding command: %s\n", pterm.Bold.Sprint(command))
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stderr, "Please input a description if desired: ")
	scanner.Scan()
	desc := scanner.Text()

	err := c.NewCommand(command, desc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add the new command: %v\n", err)
	}
}
