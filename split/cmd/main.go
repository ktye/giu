package main

import (
	"image"
	"image/color"
	"os"

	"github.com/ktye/giu/split"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow()
		if e := loop(w); e != nil {
			panic(e)
		}
		os.Exit(0)
	}()
	app.Main()
}
func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	s1 := split.Split{Axis: layout.Horizontal}
	s2 := split.Split{Axis: layout.Vertical}
	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				split1(gtx, th, &s1, &s2)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func split1(gtx layout.Context, th *material.Theme, s1, s2 *split.Split) layout.Dimensions {
	return s1.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions { return split2(gtx, th, s2) })
}
func split2(gtx layout.Context, th *material.Theme, s *split.Split) layout.Dimensions {
	return s.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Top", green)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Bottom", blue)
	})
}

/*
func exampleSplit(gtx layout.Context, th *material.Theme, s *split.Split) layout.Dimensions {
	return s.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Right", blue)
	})
}
*/

func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
}

var (
	background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

// ColorBox creates a widget with the specified dimensions and color.
func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
