package parsers

import "strconv"

// ParseIntArg parses an integer argument.
func ParseIntArg(arg string) (interface{}, error) {
	return strconv.Atoi(arg)
}