package parsers

import (
	"strconv"

	"github.com/zorbyte/whiskey/arg"
)

// ParseInt parses an integer argument.
func ParseInt(_ *arg.Argument, rawArg string) (interface{}, error) {
	return strconv.Atoi(rawArg)
}
