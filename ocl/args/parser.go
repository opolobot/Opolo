package args

// Parser parses a raw string and outputs data for use.
type Parser interface {
	Parse(raw string) (interface{}, error)
	ZeroVal() interface{}
}
