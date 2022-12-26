package config

import "github.com/tauraamui/maildew/internal/configdef"

func DefaultDestoryer() configdef.Destroyer {
	return defaultDestroyer{}
}

type defaultDestroyer struct{}

func (d defaultDestroyer) Destory() error {
	return destroy()
}
