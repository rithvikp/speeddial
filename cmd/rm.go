package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a command from speeddial",

		Run: runRm,
	}
)

func runRm(cmd *cobra.Command, args []string) {
	c := setup()
	defer dump(c)

	command := search(c)
	err := c.DeleteCommand(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete the command: %v", err)
	}
}
