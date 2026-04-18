package fyne

import (
    "fmt"
    fyne2 "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/widget"
)

type window struct {
    w           fyne2.Window
    temperature *TextWidget
}

func NewW(w fyne2.Window) *window {
    return &window{
        w: w,
    }
}

func (win *window) Resize(size WindowSize) error {
    win.w.Resize(fyne2.NewSize(float32(size.Width()), float32(size.Height())))
    return nil
}

func (win *window) UpdateTemperature(t float32) error {
    if win.temperature != nil {
        win.temperature.SetText(fmt.Sprintf("Температура: %.1f°C", t))
    }
    return nil
}

func (win *window) SetTemperatureWidget(tw *TextWidget) error {
    win.temperature = tw
    win.w.SetContent(tw.Render().(*widget.Label))
    return nil
}

func (win *window) Render() error {
    win.w.Show()
    return nil
}

// WindowSize структура для размера окна
type WindowSize struct {
    width  int
    height int
}

func NewWS(w, h int) WindowSize {
    return WindowSize{width: w, height: h}
}

func (ws WindowSize) Width() int {
    return ws.width
}

func (ws WindowSize) Height() int {
    return ws.height
}
