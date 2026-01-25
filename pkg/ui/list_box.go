//go:build windows
// +build windows

package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// ListBoxWidget はListBoxのラッパーウィジェットを表す。
type ListBoxWidget struct {
	listBox       *walk.ListBox
	tooltip       string
	minSize       declarative.Size
	maxSize       declarative.Size
	stretchFactor int
}

// NewListBoxWidget はListBoxWidgetを生成する。
func NewListBoxWidget(tooltip string) *ListBoxWidget {
	return &ListBoxWidget{tooltip: tooltip}
}

// SetMinSize は最小サイズを設定する。
func (lb *ListBoxWidget) SetMinSize(size declarative.Size) {
	lb.minSize = size
}

// SetMaxSize は最大サイズを設定する。
func (lb *ListBoxWidget) SetMaxSize(size declarative.Size) {
	lb.maxSize = size
}

// SetStretchFactor は伸長率を設定する。
func (lb *ListBoxWidget) SetStretchFactor(factor int) {
	lb.stretchFactor = factor
}

// SetWindow はウィンドウ参照を設定する（ListBoxは未使用）。
func (lb *ListBoxWidget) SetWindow(_ *controller.ControlWindow) {
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (lb *ListBoxWidget) SetEnabledInPlaying(playing bool) {
	lb.SetEnabled(!playing)
}

// SetEnabled はウィジェットの有効状態を設定する。
func (lb *ListBoxWidget) SetEnabled(enabled bool) {
	if lb == nil || lb.listBox == nil {
		return
	}
	lb.listBox.SetEnabled(enabled)
}

// SetItems はリストの表示内容を更新する。
func (lb *ListBoxWidget) SetItems(items []string) error {
	if lb == nil || lb.listBox == nil {
		return nil
	}
	return lb.listBox.SetModel(items)
}

// Widgets はUI構成を返す。
func (lb *ListBoxWidget) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:      &lb.listBox,
				ToolTipText:   lb.tooltip,
				MinSize:       lb.minSize,
				MaxSize:       lb.maxSize,
				StretchFactor: lb.stretchFactor,
			},
		},
	}
}
