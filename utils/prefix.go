package utils

// StubPrefix is a stub prefix to be used until we have guild configs.
func StubPrefix() string {
	return GetConfig().Prefix
}
