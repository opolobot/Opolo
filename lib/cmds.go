package lib

import "log"

// CmdRunner is the func to run for the cmd
type CmdRunner func(ctx *Ctx) (string, error)

// Cmd is information about the command as well as its runner
type Cmd struct {
	Runner      CmdRunner
	Name        string
	Usage       string
	Aliases     []string
	Description string
	Category    string
}

// RegCmd registers a command to whiskey.
func (w *Whiskey) RegCmd(cmd *Cmd) {
	_, ok := w.Cmds[cmd.Name]
	if ok {
		log.Fatal("Duplicate command key " + cmd.Name + " was registered")
	}

	w.Cmds[cmd.Name] = cmd

	for _, alias := range cmd.Aliases {
		_, ok = w.aliases[alias]
		if ok {
			log.Fatal("Duplicate command alias " + alias + " was registered")
		}

		w.aliases[alias] = cmd.Name
	}
}

// FindCmd finds the cmd with either the name or the alias
func (w *Whiskey) FindCmd(name string) *Cmd {
	cmd, ok := w.Cmds[name]
	if !ok {
		name, ok = w.aliases[name]
		if ok {
			return w.FindCmd(name)
		}

		return nil
	}

	return cmd
}
