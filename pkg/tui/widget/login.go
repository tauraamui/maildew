package widget

import "github.com/jroimartin/gocui"

type Widget interface {
	Layout(g *gocui.Gui) error
}

type login struct {
	name    string
	x, y, w int
}

func NewLogin(name string, x, y, w int) Widget {
	return &login{name, x, y, w}
}

func (w *login) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	return nil
}
