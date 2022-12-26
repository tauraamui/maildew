package config

import "github.com/tauraamui/maildew/internal/configdef"

func DefaultCreator() configdef.Creator {
	return defaultCreator{}
}

type defaultCreator struct{}

func (d defaultCreator) Create() error {
	return create()
}
