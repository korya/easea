package main

import (
	"fmt"
	"os"
)

const stateFileEnvVar = "EASEA_STATE_FILE"

type State struct {
	Filename string
	Parsers  []string
}

func loadState() State {
	return State{
		Filename: os.Getenv(stateFileEnvVar),
	}
}

func (s State) IsDefined() bool {
	return s.Filename != ""
}

func (s State) IsInitialized() bool {
	if !s.IsDefined() {
		return false
	}

	if _, err := os.Stat(s.Filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func (s State) Validate() error {
	fi, err := os.Stat(s.Filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %q does not exist", s.Filename)
	}

	if m := fi.Mode(); !m.IsRegular() {
		return fmt.Errorf("file %q is not a regular file (%s)", s.Filename, m)
	}

	f, err := os.OpenFile(s.Filename, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("file %q: cannot write to file: %s", s.Filename, err)
	}
	defer f.Close()

	return nil
}
