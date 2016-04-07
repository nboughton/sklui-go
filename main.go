package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	cmdBuffer = []string{}
	cmdIdx    = 0
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
		fmt.Fprintln(v, "Welcome!")
	}

	// define text entry screen
	if v, err := g.SetView("input", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// View settings
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
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, inputLine); err != nil {
		return err
	}

	// Arrow up/down scrolls cmd history
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollHistory(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollHistory(v, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func inputLine(g *gocui.Gui, v *gocui.View) error {
	line := strings.TrimSpace(v.Buffer())
	if line != "" {
		cmdBuffer = append(cmdBuffer, line)
		cmdIdx = len(cmdBuffer)
	} else {
		// it's an empty line, return nil and be done with it
		return nil
	}

	// Get a main view obvject to print output to
	ov, _ := g.View("main")

	// Parse if it is an internal command or otherwise
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
		case "/quit":
			return gocui.ErrQuit
		case "/clear":
			ov.Clear()
		case "/printInputBuffer":
			fmt.Fprintln(ov, "INPUT BUFFER:")
			fmt.Fprintln(ov, cmdBuffer)
		case "/clearInputBuffer":
			cmdBuffer = nil
		}
	} else {
		// print to output and do whatever with it
		fmt.Fprintln(ov, "cmd:", line)
	}

	// Clear the input buffer now that the line has been dealt with
	v.Clear()
	return nil
}

func scrollHistory(v *gocui.View, dy int) {
	if v != nil {
		v.Clear()
		if i := cmdIdx + dy; i >= 0 && i < len(cmdBuffer) {
			cmdIdx = i
			fmt.Fprintf(v, "%v", cmdBuffer[cmdIdx])
			v.SetOrigin(0, 0)
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
