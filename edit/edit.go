package main

import (
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	editor.SetText("")
	if len(os.Args) > 1 {
		b, e := ioutil.ReadFile(os.Args[1])
		if e != nil {
			editor.SetText(e.Error())
		} else {
			editor.SetText(string(b))
		}
	}
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}
func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
func loop(w *app.Window) error {
	// th := material.NewTheme(gofont.Collection())
	th := material.NewTheme(collection)
	th.TextSize = unit.Sp(20)
	//	th.Palette.Bg = rgb(0xffffea)         // 0xffffea
	//	th.Palette.ContrastBg = rgb(0xdaf5f4) // 0xeeee9e

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				material.Editor(th, editor, "Hint").Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

var editor = new(widget.Editor)
var collection []text.FontFace

func init() {
	ttf, err := ioutil.ReadFile("font.ttf")
	if err == nil {
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(err)
		}
		collection = []text.FontFace{text.FontFace{Font: text.Font{}, Face: face}}
	} else {
		collection = gofont.Collection()
	}
}
