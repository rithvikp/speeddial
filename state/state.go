package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// command is a fundamental unit that is some string that can be run in a shell along with
// additional metadata.
type command struct {
	Invocation  string `json:"i"`
	Description string `json:"d"`
}

// state is made up primarily of a set of commands. A state is the module that is stored persistently
// and can be shared.
type state struct {
	home     bool
	path     string
	Commands []*command `json:"c"`
}

// Container encapsulates the various states loaded.
type Container struct {
	states []*state
}

// initFile creates a new speeddial state file at the given path.
func initFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var s state
	return json.NewEncoder(f).Encode(&s)
}

// Init initializes the state container, also loading in the home state file.
func Init() (*Container, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch your home directory: %v", err)
	}

	var c Container

	err = c.Load(filepath.Join(u.HomeDir, ".speeddial"))
	if err != nil {
		return nil, fmt.Errorf("unable to load your home speeddial state: %v", err)
	}

	// TODO: Fix this hack
	c.states[0].home = true

	return &c, nil
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

	var s state
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return err
	}
	s.path = path

	c.states = append(c.states, &s)

	return nil
}

// Dump stores the contents of every state to disk.
func (c *Container) Dump() {
	for _, s := range c.states {
		f, err := os.Create(s.path)
		if err != nil {
			fmt.Printf("Unable to open the file at %s to dump state: %v\n", s.path, err)
			continue
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(s); err != nil {
			fmt.Printf("Unable to dump the state at %s: %v\n", s.path, err)
			continue
		}
	}
}

// NewCommand creates a new command in the home state with the given invocation string and
// description.
func (c *Container) NewCommand(invocation, desc string) error {
	for _, s := range c.states {
		if !s.home {
			continue
		}
		s.newCommand(invocation, desc)
		return nil
	}

	return errors.New("there is no main state to which the new command should be added")
}

func (s *state) newCommand(invocation, desc string) {
	var c command
	c.Invocation = invocation
	c.Description = desc

	// TODO: Check for duplicates
	s.Commands = append(s.Commands, &c)
}
