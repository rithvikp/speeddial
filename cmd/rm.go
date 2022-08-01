package cmd

import (
	"fmt"
	"os"

	"github.com/rithvikp/speeddial/term"
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a command from speeddial",
		Long:  `Use the search menu and the arrow keys to select an entry and then press \"enter\" to delete it. Confirm this deletion by pressing "y"`,

		Run: runRm,
	}
)

func runRm(cmd *cobra.Command, args []string) {
	c := setup()
	defer dump(c)

	command := search(c)

	confirm, err := term.Confirmation(fmt.Sprintf("Are you sure you want to delete command `%s`?", command.Invocation), true)
	if err == term.ErrUserQuit {
		os.Exit(0) // nolint:gocritic // It is ok that the deferred dump does not run since there was no state update.
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to confirm deletion: %v\n", err)
		os.Exit(1)
	} else if !confirm {
		fmt.Fprintln(os.Stderr, "Deletion cancelled")
		os.Exit(0)
	}

	err = c.DeleteCommand(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete the command: %v\n", err)
	}
}
