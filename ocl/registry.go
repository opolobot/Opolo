package ocl

import (
	"fmt"
	"log"

	"github.com/acomagu/trie"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var regInst *Registry

// Registry holds all the commands and uses a radix trie for lookups.
type Registry struct {
	tree trie.Tree
	cmds []*Command

	Categories map[string]*Category
}

// AddCategory adds a command category and registers all of its commands.
func (reg *Registry) AddCategory(cat *Category) {
	log.Printf("Registering category %v\n", cat.name)

	reg.Categories[cat.name] = cat
	for _, cmd := range cat.Commands {
		reg.addCommand(cmd)
	}
}

func (reg *Registry) addCommand(cmd *Command) {
	log.Printf("Registering command %v\n", cmd.Name)
	reg.cmds = append(reg.cmds, cmd)
}

// LookupCommand looks up a command using either its name or alias.
func (reg *Registry) LookupCommand(cmdNameOrAlias string) (*Command, error) {
	cmdInterface, ok := reg.tree.Trace([]byte(cmdNameOrAlias)).Terminal()
	if !ok {
		return nil, nil
	}

	cmd, ok := cmdInterface.(*Command)
	if !ok {
		return nil, fmt.Errorf("Failed to assert cmdInterface type as cmd pointer. cmd name: %v", cmdNameOrAlias)
	}

	return cmd, nil
}

// FindClosestCmdMatch aids in supplying "did you mean" functionality for a command.
func (reg *Registry) FindClosestCmdMatch(nonExistentCmd string) (string, int) {
	nonExtCmdRunes := []rune(nonExistentCmd)
	shortestDistance := 100
	var bestRunes []rune
	for _, cmd := range reg.cmds {
		runes := []rune(cmd.Name)
		dist := levenshtein.DistanceForStrings(nonExtCmdRunes, runes, levenshtein.DefaultOptions)
		if dist < shortestDistance {
			bestRunes = runes
			shortestDistance = dist
		}
	}

	if len(bestRunes) == 0 {
		return "", 0
	}

	return string(bestRunes), shortestDistance
}

// Populate constructs the radix trie that looks up commands.
func (reg *Registry) Populate() {
	var keys [][]byte
	var vals []interface{}
	for _, cmd := range reg.cmds {
		keys = append(keys, []byte(cmd.Name))
		for _, alias := range cmd.Aliases {
			keys = append(keys, []byte(alias))
			vals = append(vals, cmd)
		}

		vals = append(vals, cmd)
	}

	reg.tree = trie.New(keys, vals)
}

// GetRegistry gets the command registry.
func GetRegistry() *Registry {
	if regInst == nil {
		regInst = &Registry{
			Categories: make(map[string]*Category),
		}
	}

	return regInst
}
