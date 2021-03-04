package split

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

// Split is a container widget with two children, that has a draggable separator in between.
// Axis is the orientation of the children.
// The axis is flipped with right mouse button clicks on the separator.
type Split struct {
	layout.Axis
	Ratio   float32    // Ratio keeps the current layout. 0 is center, -1 completely to the left, 1 completely to the right.
	Bar     unit.Value // Bar is the width for resizing the layout
	drag    bool
	dragID  pointer.ID
	dragPos float32
}

var defaultBarWidth = unit.Dp(10)

func (s *Split) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	bar := gtx.Px(s.Bar)
	if bar <= 1 {
		bar = gtx.Px(defaultBarWidth)
	}
	maxval := gtx.Constraints.Max.X
	vertical := s.Axis == layout.Vertical
	if vertical {
		maxval = gtx.Constraints.Max.Y
	}
	proportion := (s.Ratio + 1) / 2
	leftsize := int(proportion*float32(maxval) - float32(bar)/2)
	rightoffset := leftsize + bar
	rightsize := maxval - rightoffset
	{
		stack := op.Save(gtx.Ops)
		for _, ev := range gtx.Events(s) {
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}
			switch e.Type {
			case pointer.Press:
				if s.drag {
					break
				}
				if e.Buttons == pointer.ButtonSecondary { // flip orientation
					if vertical {
						s.Axis = layout.Horizontal
					} else {
						s.Axis = layout.Vertical
					}
					break
				}
				s.dragID = e.PointerID
				if vertical {
					s.dragPos = e.Position.Y
				} else {
					s.dragPos = e.Position.X
				}
			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}
				var delta float32
				if vertical {
					delta = e.Position.Y - s.dragPos
					s.dragPos = e.Position.Y
				} else {
					delta = e.Position.X - s.dragPos
					s.dragPos = e.Position.X
				}
				deltaRatio := delta * 2 / float32(maxval)
				s.Ratio += deltaRatio
			case pointer.Release:
				if s.Ratio < -1 {
					s.Ratio = -1
				} else if s.Ratio > 1 {
					s.Ratio = 1.0
				}
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}
		var barRect image.Rectangle
		if vertical {
			barRect = image.Rect(0, leftsize, gtx.Constraints.Max.X, rightoffset)
		} else {
			barRect = image.Rect(leftsize, 0, rightoffset, gtx.Constraints.Max.Y)
		}
		pointer.Rect(barRect).Add(gtx.Ops)
		pointer.InputOp{Tag: s,
			Types: pointer.Press | pointer.Drag | pointer.Release,
			Grab:  s.drag,
		}.Add(gtx.Ops)
		stack.Load()
	}
	{
		stack := op.Save(gtx.Ops)
		gtx := gtx
		if vertical {
			gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, leftsize))
		} else {
			gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		}
		left(gtx)
		stack.Load()
	}
	{
		stack := op.Save(gtx.Ops)
		if vertical {
			op.Offset(f32.Pt(0, float32(rightoffset))).Add(gtx.Ops)
		} else {
			op.Offset(f32.Pt(float32(rightoffset), 0)).Add(gtx.Ops)
		}
		gtx := gtx
		if vertical {
			gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, rightsize))
		} else {
			gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		}
		right(gtx)
		stack.Load()
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
