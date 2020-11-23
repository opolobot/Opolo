package parsers

import "strings"

// parser is here because we can't do a circular import
// so we'll just remake the interface here.
type parser interface {
	Parse(raw string) (interface{}, error)
	ZeroVal() interface{}
}

// KeyValue is a key=value parser.
type KeyValue struct {
	Key    string
	Parser parser
}

// Parse parses the value from a key=value string if the key matches.
func (kv *KeyValue) Parse(raw string) (interface{}, error) {
	keyAndVal := strings.SplitN(raw, "=", 2)
	if len(keyAndVal) == 0 || keyAndVal[0] != kv.Key {
		return "", nil
	}

	return kv.Parser.Parse(keyAndVal[1])
}

// ZeroVal returns the zero value for the key value parser.
func (kv *KeyValue) ZeroVal() interface{} {
	return kv.Parser.ZeroVal()
}
