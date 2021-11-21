package tui

import (
	"errors"

	"github.com/tauraamui/maildew/pkg/log"
	"github.com/tauraamui/maildew/pkg/tui/widget"

	"github.com/jroimartin/gocui"
)

type ui struct {
	goCui *gocui.Gui
	login widget.Widget
}

func newUi() (*ui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	i := ui{goCui: g}
	i.init()
	return &i, nil
}

func (u *ui) init() {
	u.registerWidgets()

	if err := u.goCui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, u.quit); err != nil {
		log.Fatal(err.Error())
	}

	if err := u.goCui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err.Error())
	}
}

func (u *ui) registerWidgets() {
	maxX, maxY := u.goCui.Size()
	u.login = widget.NewLogin("login", (maxX/2)-5, maxY/2, 10)

	u.goCui.SetManager(u, u.login)
}

func (u *ui) Layout(g *gocui.Gui) error {
	if u.goCui != g {
		return errors.New("root ui layout callback given incorrect gocui.Gui instance pointer")
	}
	return nil
}

func (u *ui) quit(g *gocui.Gui, v *gocui.View) error {
	u.Close()
	return gocui.ErrQuit
}

func (u *ui) Close() {
	if u.goCui != nil {
		u.goCui.Close()
	}
}
