package tui

import "github.com/tauraamui/maildew/pkg/log"

type UI interface {
	Close()
}

func New() UI {
	inst, err := newUi()
	if err != nil {
		log.Fatal(err.Error())
	}
	return inst
}
