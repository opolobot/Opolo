package cmds

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/acomagu/trie"
	"github.com/bwmarrin/discordgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"github.com/zorbyte/whiskey/args"
	"github.com/zorbyte/whiskey/utils"
)

var stringSplitter = regexp.MustCompile(" +")

var cmdUIInstance *CommandUI

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
	HelpOutputs map[string]*map[string]string
}

// AddCategory adds a command category and registers all of its commands.
func (cmdUI *CommandUI) AddCategory(cmdCat *CommandCategory) {
	log.Printf("Registering category %v\n", cmdCat.name)

	helpOutputs := make(map[string]string)
	cmdUI.HelpOutputs[cmdCat.DisplayName()] = &helpOutputs
	for _, cmd := range cmdCat.cmds {
		helpOutputs[cmd.Name] = generateHelpStr(cmd)
		cmdUI.addCommand(cmd)
	}
}

func (cmdUI *CommandUI) addCommand(cmd *Command) {
	log.Printf("Registering command %v\n", cmd.Name)
	cmdUI.cmdRunes = append(cmdUI.cmdRunes, []rune(cmd.Name))
	for _, alias := range cmd.Aliases {
		cmdUI.aliases[alias] = cmd.Name
	}

	cmdUI.rawCmds[cmd.Name] = cmd
}

// Dispatch dispatches a command.
func (cmdUI *CommandUI) Dispatch(session *discordgo.Session, msg *discordgo.Message) NextFunc {
	startTime := time.Now()

	config := utils.GetConfig()
	if !strings.HasPrefix(msg.Content, config.Prefix) {
		return nil
	}

	cmdStr := strings.TrimSpace(msg.Content)[len(config.Prefix):]

	cmdSegs := stringSplitter.Split(cmdStr, -1)
	cmdCallKey := cmdSegs[0]
	cmdArgs := cmdSegs[1:]
	log.Printf("cmdName -> %v, cmdArgs -> %v\n", cmdCallKey, cmdArgs)

	ctx := &Context{
		Session: session,
		Msg:     msg,

		CmdCallKey: cmdCallKey,

		StartTime: startTime,
	}

	cmd, err := cmdUI.LookupCommand(cmdCallKey)
	if err != nil {
		ctx.SendError(err)
	}

	if cmd == nil {
		closest, distance := cmdUI.FindClosestCmdMatch(cmdCallKey)
		if distance <= 2 && distance != 0 {
			ctx.Send(fmt.Sprintf("**:question: ~ Did you mean `%v`?**", config.Prefix + closest))
		}
	} else {
		ctx.Cmd = cmd

		ctx.ArgsCtx, err = args.NewArgumentsContext(session, cmd.argCodecs, cmdArgs)
		if err != nil {
			if pErr, ok := err.(*args.ParsingError); ok {
				err = errors.Unwrap(pErr)
				if err != nil {
					ctx.SendError(err)
				} else {
					ctx.Send(err.Error())
					return nil
				}
			}
		}

		var nextFunc NextFunc
		idx := -1
		nextFunc = func(err ...error) {
			if idx == -1 {
				defer (func() {
					if r := recover(); r != nil {
						ctx.SendError(fmt.Errorf("%v", r))
					}
				})()
			}

			if len(err) > 0 {
				ctx.SendError(err[0])
				return
			} else if idx++; idx <= len(cmd.stack)-1 {
				cmd.stack[idx](ctx, nextFunc)
			}

			log.Printf("Time for command execution: %v", time.Since(startTime))

			// TODO(@zorbyte): Does this need a goroutine?
			// Consider enabling this.
			// go ctx.CleanUp()
		}

		return nextFunc
	}

	return nil
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
	shortestDistance := 100
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
	cmdHelpStrBldr.WriteString(cmd.Name)
	if usage := cmd.Usage(); usage != "" {
		cmdHelpStrBldr.WriteString(" ")
		cmdHelpStrBldr.WriteString(usage)
	}

	// The first **` is missing for dynamic prefix support, as
	// including it and styling the dynamic prefix would ruin the formatting
	cmdHelpStrBldr.WriteString("`**")

	// Description.
	if cmd.Description != "" {
		cmdHelpStrBldr.WriteString("\n> **desc ~** ")
		cmdHelpStrBldr.WriteString(cmd.Description)
	}

	// Aliases
	if len(cmd.Aliases) > 0 {
		cmdHelpStrBldr.WriteString("\n> **aliases ~** ")
		cmdHelpStrBldr.WriteString(strings.Join(cmd.Aliases, ", "))
	}

	return cmdHelpStrBldr.String()
}

// GetCommandUI gets the command user interface manager.
func GetCommandUI() *CommandUI {
	if cmdUIInstance == nil {
		cmdUIInstance = newCommandUI()
	}

	return cmdUIInstance
}

func newCommandUI() *CommandUI {
	return &CommandUI{
		aliases: make(map[string]string),
		rawCmds: make(map[string]*Command),

		HelpOutputs: make(map[string]*map[string]string),
	}
}
