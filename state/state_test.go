package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStateOperations(t *testing.T) {
	commands := []*Command{
		{"git push", "Push changes", nil},
		{"git commit", "Commit changes", nil},
		{"git add", "Stage changes", nil},
	}

	less := func(a, b *Command) bool {
		if a.Invocation < b.Invocation {
			return true
		} else if a.Invocation > b.Invocation {
			return false
		}
		return a.Description <= b.Description
	}
	checkNoErr := func(err error) {
		t.Helper()
		if err != nil {
			t.Error(err)
		}
	}
	checkCommands := func(got, want []*Command) {
		t.Helper()
		if diff := cmp.Diff(got, want, cmpopts.IgnoreUnexported(Command{}), cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("Unexpected commands loaded (-got, +want):\n%s", diff)
		}
	}

	// ================================================
	// Start Test
	// ================================================
	dir := t.TempDir()
	statePath := filepath.Join(dir, "speeddial.json")

	// Initialize c1 with the (not-yet-created) config, add a few commands, and write it out
	c1, err := initialize(statePath)
	if err != nil {
		t.Errorf("Unable to load the first config file into c1: %v", err)
	}

	if _, err = os.Stat(statePath); os.IsNotExist(err) {
		t.Errorf("A config file was not created by Load")
	} else if err != nil {
		t.Errorf("Unable to check if a config file was created: %v", err)
	}

	checkNoErr(c1.NewCommand(commands[0].Invocation, commands[0].Description))
	checkNoErr(c1.NewCommand(commands[1].Invocation, commands[1].Description))

	c1.Dump()

	// Initialize c2 with the (already-created) config, ensure the previously-added commands
	// were persisted, add another command, and write it out
	c2, err := initialize(statePath)
	if err != nil {
		t.Errorf("Unable to load the first config file into c2: %v", err)
	}

	checkCommands(c2.List(), commands[:2])

	checkNoErr(c2.NewCommand(commands[2].Invocation, commands[2].Description))
	checkCommands(c2.List(), commands[:3])

	c2.Dump()

	// TODO: Check for duplicates once duplicate detection is added. Also test loading multiple
	// configs
}
