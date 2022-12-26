package config

import (
	"github.com/tauraamui/maildew/internal/configdef"
)

func DefaultResolver() configdef.Resolver {
	return defaultResolver{}
}

type defaultResolver struct{}

func (d defaultResolver) Resolve() (configdef.Values, error) {
	return load()
}
