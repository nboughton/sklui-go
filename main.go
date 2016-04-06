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
		fmt.Fprintln(v, "Welcome to disc-go!")
	}

	// define text entry screen
	if v, err := g.SetView("input", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// View settings
		//v.Autoscroll = true
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
	// Mouse cursor up needs to select the correct line from input
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollHistory(g, v, -1)
			return nil
		}); err != nil {
		return err
	}
	// Mouse cursor down needs to select the correct line from input
	if err := g.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollHistory(g, v, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func inputLine(g *gocui.Gui, v *gocui.View) error {
	line, err := v.Line(-1)
	if err != nil {
		line = ""
	}

	if line != "" {
		cmdBuffer = append(cmdBuffer, line)
		cmdIdx = len(cmdBuffer)
	}

	ov, _ := g.View("main")

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
			fmt.Fprintln(ov, cmdBuffer)
		case "/clearInputBuffer":
			cmdBuffer = nil
		}
	} else {
		// print to output and send it to the mud
		fmt.Fprintln(ov, "cmd:", line)
	}

	v.Clear()
	return nil
}

func scrollHistory(g *gocui.Gui, v *gocui.View, dy int) error {
	if v != nil {
		v.Clear()
		if i := cmdIdx + dy; i >= 0 && i < len(cmdBuffer) {
			cmdIdx = i
			fmt.Fprintf(v, "%v", cmdBuffer[cmdIdx])
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
