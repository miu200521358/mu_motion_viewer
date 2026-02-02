//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mpresenter/messages"
)

// TreeViewWidget はツリービュー表示のウィジェットを表す。
type TreeViewWidget struct {
	container      *walk.Composite
	treeView       *walk.TreeView
	copyButton     *walk.PushButton
	copyIcon       *walk.Bitmap
	model          *TreeModel
	translator     i18n.II18n
	logger         logging.ILogger
	minSize        declarative.Size
	maxSize        declarative.Size
	stretchFactor  int
	hoverPath      string
	onFileSelected func(string)
	onCopyPath     func(string)
}

// NewTreeViewWidget はTreeViewWidgetを生成する。
func NewTreeViewWidget(translator i18n.II18n, logger logging.ILogger, onFileSelected func(string), onCopyPath func(string)) *TreeViewWidget {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	return &TreeViewWidget{
		translator:     translator,
		logger:         logger,
		model:          NewTreeModel(),
		onFileSelected: onFileSelected,
		onCopyPath:     onCopyPath,
		copyIcon:       loadCopyIcon(logger),
	}
}

// SetMinSize は最小サイズを設定する。
func (tw *TreeViewWidget) SetMinSize(size declarative.Size) {
	if tw == nil {
		return
	}
	tw.minSize = size
}

// SetMaxSize は最大サイズを設定する。
func (tw *TreeViewWidget) SetMaxSize(size declarative.Size) {
	if tw == nil {
		return
	}
	tw.maxSize = size
}

// SetStretchFactor は伸長率を設定する。
func (tw *TreeViewWidget) SetStretchFactor(factor int) {
	if tw == nil {
		return
	}
	tw.stretchFactor = factor
}

// SetWindow はウィンドウ参照を設定する（TreeViewは未使用）。
func (tw *TreeViewWidget) SetWindow(_ *controller.ControlWindow) {
	if tw == nil {
		return
	}
	if tw.container != nil {
		tw.container.Synchronize(func() {
			tw.updateLayout()
		})
	}
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (tw *TreeViewWidget) SetEnabledInPlaying(playing bool) {
	if tw == nil || tw.treeView == nil {
		return
	}
	// 再生中も選択変更を許可する。
	tw.treeView.SetEnabled(true)
}

// SetModelPaths はルートパス一覧からツリーを再構築する。
func (tw *TreeViewWidget) SetModelPaths(paths []string) error {
	if tw == nil {
		return nil
	}
	tw.hideCopyButton()
	if tw.model == nil {
		tw.model = NewTreeModel()
		if tw.treeView != nil {
			if err := tw.treeView.SetModel(tw.model); err != nil {
				return err
			}
		}
	}
	if err := tw.model.SetRoots(paths); err != nil {
		tw.updateLayout()
		return err
	}
	tw.updateLayout()
	return nil
}

// Widgets はUI構成を返す。
func (tw *TreeViewWidget) Widgets() declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			// 内側でオーバーレイ配置を行うため、外側はサイズ確保用に分離する。
			declarative.Composite{
				AssignTo:      &tw.container,
				Layout:        declarative.VBox{MarginsZero: true, SpacingZero: true},
				StretchFactor: tw.stretchFactor,
				MinSize:       tw.minSize,
				MaxSize:       tw.maxSize,
				OnSizeChanged: func() {
					tw.updateLayout()
				},
				OnBoundsChanged: func() {
					tw.updateLayout()
				},
				Children: []declarative.Widget{
					declarative.TreeView{
						AssignTo:             &tw.treeView,
						Model:                tw.model,
						StretchFactor:        1,
						ToolTipText:          i18n.TranslateOrMark(tw.translator, messages.LabelTreeViewTip),
						OnCurrentItemChanged: tw.handleCurrentItemChanged,
						OnExpandedChanged: func(_ walk.TreeItem) {
							tw.hideCopyButton()
						},
						OnMouseMove: func(x, y int, _ walk.MouseButton) {
							tw.handleMouseMove(x, y)
						},
						OnMouseDown: func(x, y int, _ walk.MouseButton) {
							tw.handleMouseMove(x, y)
						},
						OnMouseUp: func(x, y int, _ walk.MouseButton) {
							tw.handleMouseMove(x, y)
						},
					},
					declarative.PushButton{
						AssignTo:    &tw.copyButton,
						Text:        "",
						Image:       tw.copyIcon,
						MinSize:     declarative.Size{Width: 1, Height: 1},
						MaxSize:     declarative.Size{Width: 1, Height: 1},
						Visible:     false,
						ToolTipText: i18n.TranslateOrMark(tw.translator, messages.LabelPathCopyTip),
						OnClicked:   tw.handleCopyClicked,
					},
				},
			},
		},
	}
}

// updateLayout は内部ウィジェットのサイズを調整する。
func (tw *TreeViewWidget) updateLayout() {
	if tw == nil || tw.container == nil || tw.treeView == nil {
		return
	}
	bounds := tw.container.ClientBounds()
	tw.treeView.SetBounds(bounds)
	tw.hideCopyButton()
}

// handleCurrentItemChanged は選択変更時の処理を行う。
func (tw *TreeViewWidget) handleCurrentItemChanged() {
	if tw == nil || tw.treeView == nil {
		return
	}
	item := tw.treeView.CurrentItem()
	node, ok := item.(*TreeNode)
	if !ok || node == nil || node.IsDir() {
		return
	}
	if tw.onFileSelected != nil {
		tw.onFileSelected(node.Path())
	}
}

// handleMouseMove はホバー状態に応じてコピーアイコンを表示する。
func (tw *TreeViewWidget) handleMouseMove(x, y int) {
	if tw == nil || tw.treeView == nil || tw.copyButton == nil {
		return
	}
	item := tw.treeView.ItemAt(x, y)
	node, ok := item.(*TreeNode)
	if !ok || node == nil || node.IsDir() {
		tw.hideCopyButton()
		return
	}
	path := node.Path()
	if path == "" {
		tw.hideCopyButton()
		return
	}
	tw.hoverPath = path
	tw.updateCopyButtonBounds(x, y)
}

// updateCopyButtonBounds はコピーアイコンの位置を更新する。
func (tw *TreeViewWidget) updateCopyButtonBounds(mouseX, mouseY int) {
	if tw == nil || tw.treeView == nil || tw.copyButton == nil {
		return
	}
	bounds := tw.treeView.ClientBounds()
	if bounds.Width == 0 || bounds.Height == 0 {
		return
	}
	iconSize := tw.treeView.IntFrom96DPI(16)
	margin := tw.treeView.IntFrom96DPI(4)
	itemHeight := tw.treeView.ItemHeight()
	if itemHeight <= 0 {
		itemHeight = tw.treeView.IntFrom96DPI(20)
	}
	rowTop := 0
	if itemHeight > 0 {
		rowTop = (mouseY / itemHeight) * itemHeight
	}
	x := bounds.Width - iconSize - margin
	y := rowTop + (itemHeight-iconSize)/2
	if x < 0 || y < 0 {
		tw.hideCopyButton()
		return
	}
	tw.copyButton.SetBounds(walk.Rectangle{X: x, Y: y, Width: iconSize, Height: iconSize})
	tw.copyButton.SetVisible(true)
}

// handleCopyClicked はコピーアイコン押下時の処理を行う。
func (tw *TreeViewWidget) handleCopyClicked() {
	if tw == nil {
		return
	}
	path := tw.hoverPath
	if path == "" {
		return
	}
	if tw.onCopyPath != nil {
		tw.onCopyPath(path)
	}
}

// hideCopyButton はコピーアイコンを非表示にする。
func (tw *TreeViewWidget) hideCopyButton() {
	if tw == nil || tw.copyButton == nil {
		return
	}
	tw.hoverPath = ""
	tw.copyButton.SetVisible(false)
}
