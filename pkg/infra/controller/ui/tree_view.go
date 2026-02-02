//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"path/filepath"
	"unsafe"

	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mtree"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

const (
	copyIconSize     = 16
	copyButtonPad    = 4
	copyButtonMargin = 4
)

// TreeViewWidget はツリービュー表示を担当する。
type TreeViewWidget struct {
	window           *controller.ControlWindow
	translator       i18n.II18n
	logger           logging.ILogger
	container        *walk.Composite
	treeView         *walk.TreeView
	copyButton       *walk.PushButton
	copyBitmap       *walk.Bitmap
	copyPath         string
	stretchFactor    int
	onFileSelected   func(*controller.ControlWindow, string)
	onFoldersDropped func(*controller.ControlWindow, []string)
	model            *treeViewModel
}

// NewTreeViewWidget はTreeViewWidgetを生成する。
func NewTreeViewWidget(translator i18n.II18n, logger logging.ILogger) *TreeViewWidget {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	bitmap, err := loadCopyBitmap()
	if err != nil {
		logger.Warn("コピーアイコンの読み込みに失敗しました: %s", err.Error())
	}
	return &TreeViewWidget{
		translator: translator,
		logger:     logger,
		copyBitmap: bitmap,
	}
}

// SetWindow はウィンドウ参照を設定する。
func (tv *TreeViewWidget) SetWindow(window *controller.ControlWindow) {
	tv.window = window
	if tv.treeView != nil {
		tv.treeView.DropFiles().Attach(func(files []string) {
			if tv.onFoldersDropped != nil {
				tv.onFoldersDropped(tv.window, files)
			}
		})
	}
	if window != nil {
		window.MouseMove().Attach(func(_, _ int, _ walk.MouseButton) {
			if !tv.isCursorInsideTreeView() {
				tv.hideCopy()
			}
		})
	}
	if tv.copyButton != nil {
		_ = tv.copyButton.SetImage(tv.copyBitmap)
		tv.copyButton.SetVisible(false)
		tv.copyButton.SetText("")
	}
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (tv *TreeViewWidget) SetEnabledInPlaying(playing bool) {
	if tv == nil || tv.treeView == nil {
		return
	}
	tv.treeView.SetEnabled(!playing)
	if playing {
		tv.hideCopy()
	}
}

// SetStretchFactor は伸長率を設定する。
func (tv *TreeViewWidget) SetStretchFactor(factor int) {
	tv.stretchFactor = factor
}

// SetOnFileSelected はファイル選択時のコールバックを設定する。
func (tv *TreeViewWidget) SetOnFileSelected(handler func(*controller.ControlWindow, string)) {
	tv.onFileSelected = handler
}

// SetOnFoldersDropped はフォルダD&D時のコールバックを設定する。
func (tv *TreeViewWidget) SetOnFoldersDropped(handler func(*controller.ControlWindow, []string)) {
	tv.onFoldersDropped = handler
}

// SetRoots はルートノードを設定する。
func (tv *TreeViewWidget) SetRoots(nodes []*mtree.Node) {
	model := newTreeViewModel(nodes)
	tv.model = model
	if tv.treeView == nil {
		return
	}
	if err := tv.treeView.SetModel(model); err != nil {
		if tv.logger != nil {
			tv.logger.Warn("ツリービューの更新に失敗しました: %s", err.Error())
		}
	}
	tv.hideCopy()
}

// SelectPath は指定パスのノードを選択する。
func (tv *TreeViewWidget) SelectPath(path string) bool {
	if tv == nil || tv.treeView == nil || tv.model == nil {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	item := tv.model.findByPath(absPath)
	if item == nil {
		return false
	}
	if err := tv.treeView.SetCurrentItem(item); err != nil {
		return false
	}
	return true
}

// Widgets はUI構成を返す。
func (tv *TreeViewWidget) Widgets() declarative.Composite {
	buttonSize := copyIconSize + copyButtonPad*2
	return declarative.Composite{
		AssignTo:      &tv.container,
		StretchFactor: tv.stretchFactor,
		OnBoundsChanged: func() {
			tv.layoutChildren()
		},
		Children: []declarative.Widget{
			declarative.TreeView{
				AssignTo: &tv.treeView,
				OnCurrentItemChanged: func() {
					tv.handleCurrentItemChanged()
				},
				OnMouseMove: func(x, y int, _ walk.MouseButton) {
					tv.handleMouseMove(x, y)
				},
			},
			declarative.PushButton{
				AssignTo: &tv.copyButton,
				Text:     "",
				MinSize:  declarative.Size{Width: buttonSize, Height: buttonSize},
				MaxSize:  declarative.Size{Width: buttonSize, Height: buttonSize},
				Visible:  false,
				OnClicked: func() {
					tv.handleCopyClicked()
				},
			},
		},
	}
}

// handleCurrentItemChanged は選択変更を処理する。
func (tv *TreeViewWidget) handleCurrentItemChanged() {
	if tv == nil || tv.treeView == nil || tv.onFileSelected == nil {
		return
	}
	item, ok := tv.treeView.CurrentItem().(*treeViewItem)
	if !ok || item == nil {
		return
	}
	if !item.isFile() {
		return
	}
	tv.onFileSelected(tv.window, item.node.Path)
}

// handleMouseMove はホバー時にコピーアイコンを制御する。
func (tv *TreeViewWidget) handleMouseMove(x, y int) {
	item, rect, ok := tv.hitTestItem(x, y)
	if !ok || item == nil {
		tv.hideCopy()
		return
	}
	if !item.isFile() {
		tv.hideCopy()
		return
	}
	tv.showCopy(rect, item.node.Path)
}

// handleCopyClicked はコピーアイコン押下を処理する。
func (tv *TreeViewWidget) handleCopyClicked() {
	if tv == nil || tv.copyPath == "" {
		return
	}
	if err := walk.Clipboard().SetText(tv.copyPath); err != nil {
		logErrorWithTitle(tv.logger, i18n.TranslateOrMark(tv.translator, "クリップボードコピー失敗"), err)
	}
}

// hitTestItem は座標からツリーアイテムと矩形を取得する。
func (tv *TreeViewWidget) hitTestItem(x, y int) (*treeViewItem, win.RECT, bool) {
	if tv == nil || tv.treeView == nil {
		return nil, win.RECT{}, false
	}
	hti := win.TVHITTESTINFO{Pt: win.POINT{X: int32(x), Y: int32(y)}}
	tv.treeView.SendMessage(win.TVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))
	if hti.HItem == 0 || hti.Flags&win.TVHT_ONITEM == 0 {
		return nil, win.RECT{}, false
	}
	item, ok := tv.treeView.ItemAt(x, y).(*treeViewItem)
	if !ok {
		return nil, win.RECT{}, false
	}
	var rect win.RECT
	if tv.treeView.SendMessage(win.TVM_GETITEMRECT, uintptr(hti.HItem), uintptr(unsafe.Pointer(&rect))) == 0 {
		return nil, win.RECT{}, false
	}
	return item, rect, true
}

// showCopy はコピーアイコンを表示する。
func (tv *TreeViewWidget) showCopy(rect win.RECT, path string) {
	if tv.copyButton == nil || tv.treeView == nil {
		return
	}
	bounds := tv.treeView.ClientBoundsPixels()
	buttonSize := copyIconSize + copyButtonPad*2
	textRight := int(rect.Right) + copyButtonMargin
	maxX := bounds.Width - buttonSize - copyButtonMargin
	x := textRight
	if x > maxX {
		x = maxX
	}
	if x < copyButtonMargin {
		x = copyButtonMargin
	}
	rowHeight := int(rect.Bottom - rect.Top)
	if rowHeight < buttonSize {
		rowHeight = buttonSize
	}
	y := int(rect.Top) + (rowHeight-buttonSize)/2
	_ = tv.copyButton.SetBoundsPixels(walk.Rectangle{X: x, Y: y, Width: buttonSize, Height: buttonSize})
	if !tv.copyButton.Visible() {
		tv.copyButton.SetVisible(true)
	}
	tv.copyPath = path
}

// hideCopy はコピーアイコンを非表示にする。
func (tv *TreeViewWidget) hideCopy() {
	if tv == nil || tv.copyButton == nil {
		return
	}
	tv.copyPath = ""
	if tv.copyButton.Visible() {
		tv.copyButton.SetVisible(false)
	}
}

// layoutChildren は子ウィジェットの配置を行う。
func (tv *TreeViewWidget) layoutChildren() {
	if tv.container == nil || tv.treeView == nil {
		return
	}
	bounds := tv.container.ClientBoundsPixels()
	_ = tv.treeView.SetBoundsPixels(bounds)
	tv.hideCopy()
}

// isCursorInsideTreeView はカーソルがツリービュー内にあるか判定する。
func (tv *TreeViewWidget) isCursorInsideTreeView() bool {
	if tv == nil || tv.treeView == nil {
		return false
	}
	var pt win.POINT
	if !win.GetCursorPos(&pt) {
		return false
	}
	if !win.ScreenToClient(tv.treeView.Handle(), &pt) {
		return false
	}
	bounds := tv.treeView.ClientBoundsPixels()
	return int(pt.X) >= 0 && int(pt.Y) >= 0 && int(pt.X) < bounds.Width && int(pt.Y) < bounds.Height
}

// treeViewModel はTreeView用のモデルを表す。
type treeViewModel struct {
	walk.TreeModelBase
	roots []*treeViewItem
}

// newTreeViewModel はTreeViewモデルを生成する。
func newTreeViewModel(nodes []*mtree.Node) *treeViewModel {
	roots := make([]*treeViewItem, 0, len(nodes))
	for _, node := range nodes {
		roots = append(roots, newTreeViewItem(node, nil))
	}
	return &treeViewModel{roots: roots}
}

// RootCount はルート件数を返す。
func (m *treeViewModel) RootCount() int {
	return len(m.roots)
}

// RootAt はルートを取得する。
func (m *treeViewModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}

// findByPath はフルパス一致のノードを探索する。
func (m *treeViewModel) findByPath(path string) *treeViewItem {
	for _, root := range m.roots {
		if found := root.findByPath(path); found != nil {
			return found
		}
	}
	return nil
}

// treeViewItem はTreeView用のツリーアイテムを表す。
type treeViewItem struct {
	node     *mtree.Node
	parent   *treeViewItem
	children []*treeViewItem
}

// newTreeViewItem はTreeView用アイテムを生成する。
func newTreeViewItem(node *mtree.Node, parent *treeViewItem) *treeViewItem {
	item := &treeViewItem{node: node, parent: parent}
	children := make([]*treeViewItem, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, newTreeViewItem(child, item))
	}
	item.children = children
	return item
}

// Text は表示文字列を返す。
func (i *treeViewItem) Text() string {
	if i == nil || i.node == nil {
		return ""
	}
	return i.node.Name
}

// Parent は親アイテムを返す。
func (i *treeViewItem) Parent() walk.TreeItem {
	if i == nil {
		return nil
	}
	return i.parent
}

// ChildCount は子要素数を返す。
func (i *treeViewItem) ChildCount() int {
	if i == nil {
		return 0
	}
	return len(i.children)
}

// ChildAt は指定インデックスの子要素を返す。
func (i *treeViewItem) ChildAt(index int) walk.TreeItem {
	return i.children[index]
}

// isFile はファイルノードか判定する。
func (i *treeViewItem) isFile() bool {
	return i != nil && i.node != nil && i.node.IsFile
}

// findByPath は指定パスのノードを探索する。
func (i *treeViewItem) findByPath(path string) *treeViewItem {
	if i == nil || i.node == nil {
		return nil
	}
	if i.node.Path == path {
		return i
	}
	for _, child := range i.children {
		if found := child.findByPath(path); found != nil {
			return found
		}
	}
	return nil
}
