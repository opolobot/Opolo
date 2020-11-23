package args

import (
	"fmt"
	"strings"
)

const (
	openingOptionalBracket = '['
	closingOptionalBracket = ']'

	openingRequiredBracket = '<'
	closingRequiredBracket = '>'
)

type bracketsError struct {
	msg     string
	opening rune
	closing rune
}

func (err *bracketsError) Error() string {
	return fmt.Sprintf(err.msg+": opening=%v, closing=%v", err.opening, err.closing)
}

func parseID(ID string) (name string, required, greedy, keyValue bool, err error) {
	if len(ID) < 3 {
		err = fmt.Errorf("can not have an ID that is < 3 characters")
		return
	}

	opening, closing := getBrackets(ID)
	err = validateBrackets(opening, closing)
	if err != nil {
		return
	}

	required = isRequiredBrackets(opening, closing)
	greedy = isGreedyID(ID)
	name = getNameOfID(ID, greedy)

	keyName := keyValueName(name)
	keyValue = keyName != ""
	if keyValue {
		if greedy {
			err = fmt.Errorf("can not have a key value argument that is greedy")
			return
		}

		name = keyName
	}

	return
}

func getBrackets(ID string) (opening, closing rune) {
	opening = rune(ID[0])
	closing = rune(ID[len(ID)-1])

	return
}

func validateBrackets(opening, closing rune) error {
	if opening != openingOptionalBracket &&
		opening != openingRequiredBracket &&
		closing != closingOptionalBracket &&
		closing != closingRequiredBracket {
		return &bracketsError{"unknown opening or closing bracket types", opening, closing}
	}

	if opening == openingOptionalBracket &&
		closing != closingOptionalBracket ||
		opening == openingRequiredBracket &&
			closing != closingRequiredBracket {
		return &bracketsError{"bracket mismatch", opening, closing}
	}

	return nil
}

func isRequiredBrackets(opening, closing rune) bool {
	return opening == openingRequiredBracket && closing == closingRequiredBracket
}

func isGreedyID(ID string) bool {
	return ID[1:3] == "..."
}

func getNameOfID(ID string, greedy bool) string {
	nameIdxStart := 1
	if greedy {
		nameIdxStart += 3
	}

	endIdx := len(ID) - 1
	return ID[nameIdxStart:endIdx]
}

func keyValueName(name string) string {
	keyAndVal := strings.SplitN(name, "=", 2)
	if len(keyAndVal) == 0 {
		return ""
	}

	return keyAndVal[0]
}
