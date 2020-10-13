package dialog

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalwidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/widget"
)

var _ fyne.Widget = (*colorChannel)(nil)

// colorChannel controls a channel of a color and triggers the callback when changed.
type colorChannel struct {
	widget.BaseWidget
	name      string
	min, max  int
	value     int
	onChanged func(int)
}

// newColorChannel returns a new color channel control for the channel with the given name.
func newColorChannel(name string, min, max, value int, onChanged func(int)) *colorChannel {
	c := &colorChannel{
		name:      name,
		min:       min,
		max:       max,
		value:     clamp(value, min, max),
		onChanged: onChanged,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *colorChannel) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabelWithStyle(c.name, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	entry := &widget.Entry{
		Text: "0",
		OnChanged: func(text string) {
			value, err := strconv.Atoi(text)
			if err != nil {
				fyne.LogError("Couldn't parse: "+text, err)
			} else {
				c.SetValue(value)
			}
		},
		// TODO extend Entry to force MinSize to always be wide enough for 3 chars
		// TODO add number min/max validator
	}
	slider := &widget.Slider{
		Value:       0.0,
		Min:         float64(c.min),
		Max:         float64(c.max),
		Step:        1.0,
		Orientation: widget.Horizontal,
		OnChanged: func(value float64) {
			c.SetValue(int(value))
		},
	}
	r := &colorChannelRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]fyne.CanvasObject{
			label,
			slider,
			entry,
		}),
		control: c,
		label:   label,
		entry:   entry,
		slider:  slider,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *colorChannel) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetValue updates the value in this color widget
func (c *colorChannel) SetValue(value int) {
	value = clamp(value, c.min, c.max)
	if c.value == value {
		return
	}
	c.value = value
	c.Refresh()
	if f := c.onChanged; f != nil {
		f(value)
	}
}

type colorChannelRenderer struct {
	internalwidget.BaseRenderer
	control *colorChannel
	label   *widget.Label
	entry   *widget.Entry
	slider  *widget.Slider
}

func (r *colorChannelRenderer) Layout(size fyne.Size) {
	lMin := r.label.MinSize()
	eMin := r.entry.MinSize()
	r.label.Move(fyne.NewPos(0, (size.Height-lMin.Height)/2))
	r.label.Resize(fyne.NewSize(lMin.Width, lMin.Height))
	r.slider.Move(fyne.NewPos(lMin.Width, 0))
	r.slider.Resize(fyne.NewSize(size.Width-lMin.Width-eMin.Width, size.Height))
	r.entry.Move(fyne.NewPos(size.Width-eMin.Width, 0))
	r.entry.Resize(fyne.NewSize(eMin.Width, size.Height))
}

func (r *colorChannelRenderer) MinSize() fyne.Size {
	lMin := r.label.MinSize()
	sMin := r.slider.MinSize()
	eMin := r.entry.MinSize()
	return fyne.NewSize(
		lMin.Width+sMin.Width+eMin.Width,
		fyne.Max(lMin.Height, fyne.Max(sMin.Height, eMin.Height)),
	)
}

func (r *colorChannelRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.control.Size())
	canvas.Refresh(r.control)
}

func (r *colorChannelRenderer) updateObjects() {
	r.entry.SetText(strconv.Itoa(r.control.value))
	r.slider.Value = float64(r.control.value)
	r.slider.Refresh()
}