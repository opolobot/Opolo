package parsers

import "strconv"

// Int parses ints.
type Int struct {}

// Parse parses a string.
func (*Int) Parse(raw string) (interface{}, error) {
	return strconv.Atoi(raw)
}

// ZeroVal returns the zero value.
func (*Int) ZeroVal() interface{} {
	return 0
}