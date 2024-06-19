package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mazznoer/csscolorparser"

	jsoniter "github.com/json-iterator/go"
)

type colorEntry struct {
	Name       string
	Transforms []func(csscolorparser.Color) csscolorparser.Color
}

func Entry(name string, xforms ...func(csscolorparser.Color) csscolorparser.Color) colorEntry {
	return colorEntry{name, xforms}
}

type colorMapEntry struct {
	Key    string
	Colors []colorEntry
}

var backgroundEntries = []colorEntry{Entry("terminal.background"), Entry("editor.background")}

var colorMap = []colorMapEntry{
	{Key: "foreground", Colors: []colorEntry{Entry("terminal.foreground")}},
	{Key: "background", Colors: backgroundEntries},
	{Key: "selection_foreground", Colors: []colorEntry{Entry("terminal.selectionForeground"), Entry("terminal.selectionBackground", invert)}},
	{Key: "selection_background", Colors: []colorEntry{Entry("terminal.selectionBackground")}},
	{Key: "cursor", Colors: []colorEntry{Entry("terminalCursor.foreground")}},
	{Key: "cursor_text_color", Colors: []colorEntry{Entry("editorCursor.foreground")}},
	{Key: "color0", Colors: []colorEntry{Entry("terminal.ansiBlack")}},
	{Key: "color1", Colors: []colorEntry{Entry("terminal.ansiRed")}},
	{Key: "color2", Colors: []colorEntry{Entry("terminal.ansiGreen")}},
	{Key: "color3", Colors: []colorEntry{Entry("terminal.ansiYellow")}},
	{Key: "color4", Colors: []colorEntry{Entry("terminal.ansiBlue")}},
	{Key: "color5", Colors: []colorEntry{Entry("terminal.ansiMagenta")}},
	{Key: "color6", Colors: []colorEntry{Entry("terminal.ansiCyan")}},
	{Key: "color7", Colors: []colorEntry{Entry("terminal.ansiWhite")}},
	{Key: "color8", Colors: []colorEntry{Entry("terminal.ansiBrightBlack")}},
	{Key: "color9", Colors: []colorEntry{Entry("terminal.ansiBrightRed")}},
	{Key: "color10", Colors: []colorEntry{Entry("terminal.ansiBrightGreen")}},
	{Key: "color11", Colors: []colorEntry{Entry("terminal.ansiBrightYellow")}},
	{Key: "color12", Colors: []colorEntry{Entry("terminal.ansiBrightBlue")}},
	{Key: "color13", Colors: []colorEntry{Entry("terminal.ansiBrightMagenta")}},
	{Key: "color14", Colors: []colorEntry{Entry("terminal.ansiBrightCyan")}},
	{Key: "color15", Colors: []colorEntry{Entry("terminal.ansiBrightWhite")}},
}

type VSCodeTheme struct {
	Colors map[string]string `json:"colors"`
}

func (t *VSCodeTheme) Color(ces ...colorEntry) (csscolorparser.Color, bool) {
	for _, ce := range ces {
		v, ok := t.Colors[ce.Name]
		if !ok {
			continue
		}
		// remove alpha
		col, err := csscolorparser.Parse(v)
		if err != nil {
			panic("failed to parse " + v + ": " + err.Error())
		}

		if col.A < 1 {
			// alpha blend with background
			bg, _ := t.Color(backgroundEntries...)
			col.R = (col.R * col.A) + (bg.R * (1 - col.A))
			col.G = (col.R * col.G) + (bg.G * (1 - col.A))
			col.B = (col.R * col.B) + (bg.B * (1 - col.A))
			col.A = 1
		}
		for _, xform := range ce.Transforms {
			col = xform(col)
		}
		return col, true
	}
	// return black
	return csscolorparser.Color{A: 1}, false
}

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainE() error {
	if len(os.Args) != 2 {
		return errors.New("source theme required")
	}

	var theme VSCodeTheme
	if err := Get(os.Args[1], &theme); err != nil {
		return err
	}

	for _, pair := range colorMap {
		// skip unknown values to use default
		c, ok := theme.Color(pair.Colors...)
		if ok {
			fmt.Println(pair.Key, c.HexString())
		}
	}

	return nil
}

func Get(src string, theme *VSCodeTheme) error {
	URL, err := url.Parse(os.Args[1])
	if err != nil {
		return fmt.Errorf("failed to parse source: %w", err)
	}
	if len(URL.Scheme) == 0 {
		return GetFile(URL.String(), theme)
	}
	return GetURL(URL, theme)
}

func GetFile(fname string, theme *VSCodeTheme) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", os.Args[1], err)
	}
	defer f.Close()
	jsoniter.Parse(jsoniter.ConfigDefault, f, 2048).ReadVal(theme)
	return nil
}

func GetURL(URL *url.URL, theme *VSCodeTheme) error {
	resp, err := http.DefaultClient.Do(&http.Request{URL: URL})
	if err != nil {
		return fmt.Errorf("failed to get: %w", err)
	}
	defer resp.Body.Close()
	jsoniter.Parse(jsoniter.ConfigDefault, resp.Body, 2048).ReadVal(theme)
	return nil
}

func invert(c csscolorparser.Color) csscolorparser.Color {
	c.R = 1 - c.R
	c.G = 1 - c.G
	c.B = 1 - c.B
	return c
}
