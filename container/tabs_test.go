package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/neflyte/fyne/v2"
	"github.com/neflyte/fyne/v2/canvas"
	"github.com/neflyte/fyne/v2/internal/cache"
	"github.com/neflyte/fyne/v2/test"
	"github.com/neflyte/fyne/v2/theme"
	"github.com/neflyte/fyne/v2/widget"
)

func TestTabButton_Icon_Change(t *testing.T) {
	b := &tabButton{icon: theme.CancelIcon()}
	r := cache.Renderer(b)
	icon := r.Objects()[3].(*canvas.Image)
	oldResource := icon.Resource

	b.icon = theme.ConfirmIcon()
	b.Refresh()
	assert.NotEqual(t, oldResource, icon.Resource)
}

func TestTab_ThemeChange(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	a.Settings().SetTheme(theme.LightTheme())

	tabs := NewAppTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")))
	w := test.NewWindow(tabs)
	w.Resize(fyne.NewSize(180, 120))

	initial := w.Canvas().Capture()

	a.Settings().SetTheme(theme.DarkTheme())
	tabs.SelectIndex(1)
	second := w.Canvas().Capture()
	assert.NotEqual(t, initial, second)

	a.Settings().SetTheme(theme.LightTheme())
	tabs.SelectIndex(0)
	assert.Equal(t, initial, w.Canvas().Capture())
}
