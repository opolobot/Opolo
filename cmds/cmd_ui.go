package cmds

import (
	"fmt"
	"strings"

	"github.com/acomagu/trie"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// CommandUI manages commands and dispatches them.
type CommandUI struct {
	cmds     trie.Tree
	cmdRunes [][]rune
	aliases  map[string]string
	rawCmds  map[string]*Command

	// map[category display name]*map[command name]command help string
	// {
	//   "category display name": *{
	//	   "command name": "help string"
	//   }
	// }
	helpOutputs map[string]*map[string]string
}

// AddCategory adds a command category and registers all of its commands.
func (cmdUI *CommandUI) AddCategory(cmdCat *CommandCategory) {
	helpOutputs := make(map[string]string)
	cmdUI.helpOutputs[cmdCat.displayName()] = &helpOutputs
	for _, cmd := range cmdCat.cmds {
		helpOutputs[cmd.name] = generateHelpStr(cmd)
		cmdUI.addCommand(cmd)
	}
}

func (cmdUI *CommandUI) addCommand(cmd *Command) {
	cmdUI.cmdRunes = append(cmdUI.cmdRunes, []rune(cmd.name))
	for _, alias := range cmd.aliases {
		cmdUI.aliases[alias] = cmd.name
	}

	cmdUI.rawCmds[cmd.name] = cmd
}

// LookupCommand looks up a command using either its name or alias.
func (cmdUI *CommandUI) LookupCommand(cmdNameOrAlias string) (*Command, error) {
	cmdInterface, ok := cmdUI.cmds.Trace([]byte(cmdNameOrAlias)).Terminal()
	if !ok {
		cmdName, ok := cmdUI.aliases[cmdNameOrAlias]
		if !ok {
			return nil, nil
		}

		return cmdUI.LookupCommand(cmdName)
	}

	cmd, ok := cmdInterface.(*Command)
	if !ok {
		return nil, fmt.Errorf("Failed to assert cmdInterface type as cmd pointer. cmd name: %v", cmdNameOrAlias)
	}

	return cmd, nil
}

// FindClosestCmdMatch supplies "did you mean" functionality for a command.
func (cmdUI *CommandUI) FindClosestCmdMatch(nonExistentCmd string) (string, int) {
	nonExtCmdRunes := []rune(nonExistentCmd)
	var shortestDistance int
	var bestCmdRunes []rune
	for _, cmdRunes := range cmdUI.cmdRunes {
		dist := levenshtein.DistanceForStrings(nonExtCmdRunes, cmdRunes, levenshtein.DefaultOptions)
		if dist < shortestDistance {
			bestCmdRunes = cmdRunes
			shortestDistance = dist
		}
	}

	if len(bestCmdRunes) == 0 {
		return "", 0
	}

	return string(bestCmdRunes), shortestDistance
}

// Build builds the internal radix tree for faster command lookups.
func (cmdUI *CommandUI) Build() {
	var keys [][]byte
	var vals []interface{}
	for cmdName, cmd := range cmdUI.rawCmds {
		keys = append(keys, []byte(cmdName))
		vals = append(vals, cmd)
	}

	cmdUI.cmds = trie.New(keys, vals)
	cmdUI.rawCmds = nil
}

func (cmdUI *CommandUI) built() bool {
	return cmdUI.rawCmds == nil
}

func generateHelpStr(cmd *Command) string {
	// Prebuild help strings for performance.
	var cmdHelpStrBldr strings.Builder

	// Name and usage.
	cmdHelpStrBldr.WriteString(cmd.name)
	if usage := cmd.Usage(); usage != "" {
		cmdHelpStrBldr.WriteString(" ")
		cmdHelpStrBldr.WriteString(usage)
	}

	// The first **` is missing for dynamic prefix support, as
	// including it and styling the dynamic prefix would ruin the formatting
	cmdHelpStrBldr.WriteString("`**")

	// Description.
	if cmd.description != "" {
		cmdHelpStrBldr.WriteString("\n> **desc ~** ")
		cmdHelpStrBldr.WriteString(cmd.description)
	}

	// Aliases
	if len(cmd.aliases) > 0 {
		cmdHelpStrBldr.WriteString("\n> **aliases ~** ")
		cmdHelpStrBldr.WriteString(strings.Join(cmd.aliases, ", "))
	}

	return cmdHelpStrBldr.String()
}
