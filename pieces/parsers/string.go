package parsers

// String parses strings.
type String struct {}

// Parse parses a string.
func (strParser *String) Parse(raw string) (interface{}, error) {
	return raw, nil
}

// ZeroVal returns the zero value.
func (strParser *String) ZeroVal() interface{} {
	return ""
}