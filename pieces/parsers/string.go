package parsers

// String parses strings.
type String struct {}

// Parse parses a string.
func (*String) Parse(raw string) (interface{}, error) {
	return raw, nil
}

// ZeroVal returns the zero value.
func (*String) ZeroVal() interface{} {
	return ""
}