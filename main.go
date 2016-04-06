package main

import (
	"fmt"
	"log"
	"strings"

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

		// View settings
		v.Autoscroll = true
		v.Wrap = true

		// Opening message
		fmt.Fprintln(v, "Welcome to disc-go!")
	}

	// define text entry screen
	if v, err := g.SetView("input", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// View settings
		v.Autoscroll = true
		v.Editable = true
		v.Wrap = true

		// Set focus on input
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
	// Submit a line
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, submitLine); err != nil {
		return err
	}
	// Mouse cursor up needs to select the correct line from input
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(v, -1)
			return nil
		}); err != nil {
		return err
	}
	// Mouse cursor down needs to select the correct line from input
	if err := g.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(v, 1)
			return nil
		}); err != nil {
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

	// Parse if it is an internal command or to be sent to the mud
	if len(line) > 0 && string(line[0]) == "/" {
		// parse internal command
		s, args := strings.Split(line, " "), []string{}
		cmd := s[0]
		if len(s) > 1 {
			args = append(args, s[1:]...)
		} else {
			args = append(args, "")
		}

		switch cmd {
		case "/clear":
			ov.Clear()
		case "/printInputBuffer":
			fmt.Fprintln(ov, "INPUT BUFFER:")
			fmt.Fprintln(ov, v.Buffer())
		case "/clearInputBuffer":
			v.Clear()
		}
	} else {
		fmt.Fprintln(ov, "input:", line)
		// send it to the mud
	}

	return nil
}

func scroll(v *gocui.View, dy int) error {
	if v != nil {
		ox, oy := v.Cursor()
		v.MoveCursor(ox, oy+dy, true)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
