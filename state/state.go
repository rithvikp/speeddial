package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/exp/slices"
)

const (
	dumpVersion1 = 1

	// This path is relative to the user's home directory.
	primaryStatePath = ".config/speeddial/state.json"
)

// Command is a fundamental unit that is some string that can be run in a shell along with
// additional metadata.
type Command struct {
	Invocation  string `json:"i"`
	Description string `json:"d"`

	// This field is lazily set during search.
	state *state
}

// state is made up primarily of a set of commands. A state is the module that is stored persistently
// and can be shared.
type state struct {
	primary  bool
	path     string
	Commands []*Command `json:"c"`
}

// dump is a wrapper around state that is persisted and saved to a file for use across invocations
// of the program. It is used to provide extensibility (and backwards compatible) in case the
// structure of state changes.
type dump struct {
	Version int    `json:"v"`
	Data    *state `json:"d"`
}

// Container encapsulates the various states loaded.
type Container struct {
	states []*state
}

// initFile creates a new speeddial state file at the given path.
func initFile(path string) error {
	// TODO: Handle the case where something else already uses ~/.config/speeddial
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d := dump{
		Version: dumpVersion1,
		Data:    &state{},
	}
	return json.NewEncoder(f).Encode(&d)
}

// Init initializes the state container, also loading in the primary state file.
func Init() (*Container, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch your home directory: %v", err)
	}
	return initialize(filepath.Join(u.HomeDir, primaryStatePath))
}

func initialize(statePath string) (*Container, error) {
	var c Container

	err := c.Load(statePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load your primary speeddial state: %v", err)
	}

	// TODO: Fix this hack
	c.states[0].primary = true

	return &c, nil
}

func (c *Container) List() []*Command {
	var commands []*Command
	for _, s := range c.states {
		commands = append(commands, s.Commands...)
	}
	return commands
}

// Load loads the speeddial state at the given path into the provided container, creating a new one
// if one does not exist.
func (c *Container) Load(path string) error {
	path = filepath.Clean(path)

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := initFile(path); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var d dump
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return err
	}

	if d.Version < dumpVersion1 {
		return fmt.Errorf("%d is an unsupported version for state at %s", d.Version, path)
	}

	s := d.Data
	if s == nil {
		return fmt.Errorf("dump at %s does not have any state", path)
	}

	s.path = path

	c.states = append(c.states, s)

	return nil
}

// Dump stores the contents of every state to disk.
func (c *Container) Dump() {
	for _, s := range c.states {
		f, err := os.Create(s.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open the file at %s to dump state: %v\n", s.path, err)
			continue
		}

		d := dump{
			Version: dumpVersion1,
			Data:    s,
		}
		if err := json.NewEncoder(f).Encode(&d); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to dump the state at %s: %v\n", s.path, err)
			f.Close()
			continue
		}
		f.Close()
	}
}

// NewCommand creates a new command in the primary state with the given invocation string and
// description.
func (c *Container) NewCommand(invocation, desc string) error {
	for _, s := range c.states {
		if !s.primary {
			continue
		}
		s.newCommand(invocation, desc)
		return nil
	}

	return errors.New("there is no main state to which the new command should be added")
}

// DeleteCommand deletes the given command from the container.
func (c *Container) DeleteCommand(command *Command) error {
	if command.state == nil {
		return fmt.Errorf("command %q did not have a corresponding state", command.Invocation)
	}

	found := false
	s := command.state
	for i, sc := range s.Commands {
		if sc == command {
			s.Commands = slices.Delete(s.Commands, i, i+1)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("command %q was not found in the state", command.Invocation)
	}

	return nil
}

func (s *state) newCommand(invocation, desc string) {
	var c Command
	c.Invocation = invocation
	c.Description = desc

	// TODO: Check for duplicates
	s.Commands = append(s.Commands, &c)
}
