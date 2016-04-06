package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	// Initialise gui
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Set layout
	g.SetLayout(layout)

	// Assign keybindings
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true

	// Call MainLoop to instantiate ui
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// define starting size for main window (text output from server)
	if v, err := g.SetView("main", -1, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
		v.Frame = true
		fmt.Fprintln(v, "our text output appears here")
	}

	// define text entry screen
	if v, err := g.SetView("input", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true

		if err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	// Set quit
	if err := g.SetKeybinding("input", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, submitLine); err != nil {
		return err
	}
	return nil
}

func submitLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy - 1)
	if err != nil {
		line = ""
	}

	ov, err := g.View("main")
	if err != nil {
		log.Panicln(err)
	}

	fmt.Fprintln(ov, "input:", line)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
