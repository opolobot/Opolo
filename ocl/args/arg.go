package args

import (
	"strings"

	"github.com/opolobot/Opolo/pieces/parsers"
)

// Argument for commands.
type Argument struct {
	ID       string
	name     string
	required bool
	greedy   bool
	info     string

	parser Parser
}

// New creates an argument.
func New(ID string, parser Parser, info ...string) *Argument {
	name, required, greedy, keyValue, err := parseID(ID)
	if err != nil {
		panic(err)
	}

	finalInfo := getInfo(info)
	userFriendlyID := makeUserFriendlyID(ID, finalInfo)

	if keyValue {
		parser = &parsers.KeyValue{Key: name, Parser: parser}
	}

	return &Argument{userFriendlyID, name, required, greedy, finalInfo, parser}
}

func getInfo(info []string) (finalInfo string) {
	switch len(info) {
	case 0:
		break
	case 1:
		finalInfo = info[0]
	default:
		finalInfo = strings.Join(info, ", ")
	}

	return
}

func makeUserFriendlyID(originalID, finalInfo string) (ID string) {
	if finalInfo != "" {
		ID = string(originalID[:len(originalID)-2]) + "(" + finalInfo + ")" + originalID[len(originalID)-1:]
	} else {
		ID = originalID
	}

	return
}
