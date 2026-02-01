//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// FolderPicker はフォルダ選択ウィジェットを表す。
type FolderPicker struct {
	window            *controller.ControlWindow
	title             string
	tooltip           string
	historyKey        string
	translator        i18n.II18n
	userConfig        config.IUserConfig
	pathEdit          *walk.LineEdit
	openPushButton    *walk.PushButton
	historyPushButton *walk.PushButton
	historyDialog     *walk.Dialog
	historyListBox    *walk.ListBox
	prevPaths         []string
	minSize           declarative.Size
	maxSize           declarative.Size
	stretchFactor     int
	onPathsChanged    func(*controller.ControlWindow, []string)
}

// NewFolderPicker はFolderPickerを生成する。
func NewFolderPicker(userConfig config.IUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathsChanged func(*controller.ControlWindow, []string)) *FolderPicker {
	return &FolderPicker{
		title:          title,
		tooltip:        tooltip,
		historyKey:     historyKey,
		translator:     translator,
		userConfig:     userConfig,
		onPathsChanged: onPathsChanged,
	}
}

// SetMinSize は最小サイズを設定する。
func (fp *FolderPicker) SetMinSize(size declarative.Size) {
	fp.minSize = size
}

// SetMaxSize は最大サイズを設定する。
func (fp *FolderPicker) SetMaxSize(size declarative.Size) {
	fp.maxSize = size
}

// SetStretchFactor は伸長率を設定する。
func (fp *FolderPicker) SetStretchFactor(factor int) {
	fp.stretchFactor = factor
}

// SetWindow はウィンドウ参照を設定する。
func (fp *FolderPicker) SetWindow(window *controller.ControlWindow) {
	fp.window = window
}

// SetEnabledInPlaying は再生中の有効状態を設定する。
func (fp *FolderPicker) SetEnabledInPlaying(playing bool) {
	if fp == nil {
		return
	}
	enabled := !playing
	if fp.pathEdit != nil {
		fp.pathEdit.SetEnabled(enabled)
	}
	if fp.openPushButton != nil {
		fp.openPushButton.SetEnabled(enabled)
	}
	if fp.historyPushButton != nil {
		fp.historyPushButton.SetEnabled(enabled)
	}
}

// SetPath は単一フォルダパスを設定する。
func (fp *FolderPicker) SetPath(path string) {
	fp.SetPaths([]string{path})
}

// SetPaths は複数フォルダパスを設定する。
func (fp *FolderPicker) SetPaths(paths []string) {
	fp.applyPaths(paths, true)
}

// Widgets はUI構成を返す。
func (fp *FolderPicker) Widgets() declarative.Composite {
	titleWidgets := []declarative.Widget{
		declarative.TextLabel{
			Text:        fp.title,
			ToolTipText: fp.tooltip,
		},
	}

	inputWidgets := []declarative.Widget{
		declarative.LineEdit{
			AssignTo:    &fp.pathEdit,
			ToolTipText: fp.tooltip,
			OnTextChanged: func() {
				fp.handlePathChanged(fp.pathEdit.Text())
			},
			OnEditingFinished: func() {
				fp.handlePathConfirmed(fp.pathEdit.Text())
			},
			OnDropFiles: func(files []string) {
				fp.handleDropFiles(files)
			},
		},
		declarative.PushButton{
			AssignTo:    &fp.openPushButton,
			Text:        fp.t("開く"),
			ToolTipText: fp.tooltip,
			OnClicked: func() {
				fp.showOpenDialog()
			},
			MinSize: declarative.Size{Width: 70, Height: 20},
			MaxSize: declarative.Size{Width: 70, Height: 20},
		},
	}

	if fp.historyKey != "" {
		inputWidgets = append(inputWidgets, declarative.PushButton{
			AssignTo:    &fp.historyPushButton,
			Text:        fp.t("履歴"),
			ToolTipText: fp.tooltip,
			OnClicked: func() {
				fp.openHistoryDialog()
			},
			MinSize: declarative.Size{Width: 70, Height: 20},
			MaxSize: declarative.Size{Width: 70, Height: 20},
		})
	}

	return declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: titleWidgets,
			},
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: inputWidgets,
			},
		},
	}
}

// handlePathChanged はパス変更時の処理を行う。
func (fp *FolderPicker) handlePathChanged(path string) {
	// パス変更時は同一パスを再適用しない。
	fp.applyPaths([]string{path}, false)
}

// handlePathConfirmed はEnter確定時にパス変更処理を実行する。
func (fp *FolderPicker) handlePathConfirmed(path string) {
	// Enter確定時は同一パスでも再適用する。
	fp.applyPaths([]string{path}, true)
}

// handleDropFiles はドロップされたフォルダ一覧を反映する。
func (fp *FolderPicker) handleDropFiles(files []string) {
	if fp == nil || len(files) == 0 {
		return
	}
	// D&Dはフォルダのみ採用し、複数件を同時反映する。
	paths := cleanPaths(files)
	if len(paths) == 0 {
		return
	}
	fp.applyPaths(paths, true)
}

// showOpenDialog はフォルダ選択ダイアログを表示する。
func (fp *FolderPicker) showOpenDialog() {
	fd := new(walk.FileDialog)
	fd.Title = fp.title
	fd.InitialDirPath = fp.resolveInitialDir()
	ok, err := fd.ShowBrowseFolder(fp.window)
	if err != nil {
		walk.MsgBox(fp.window, fp.t("読み込み失敗"), err.Error(), walk.MsgBoxIconError)
		return
	}
	if !ok {
		return
	}
	fp.applyPaths([]string{fd.FilePath}, true)
}

// resolveInitialDir は初期ディレクトリを決定する。
func (fp *FolderPicker) resolveInitialDir() string {
	if fp.pathEdit != nil {
		current := cleanPath(fp.pathEdit.Text())
		if current != "" {
			return current
		}
	}
	if fp.historyKey == "" || fp.userConfig == nil {
		return ""
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil || len(values) == 0 {
		return ""
	}
	return values[0]
}

// openHistoryDialog は履歴ダイアログを表示する。
func (fp *FolderPicker) openHistoryDialog() {
	if fp.historyKey == "" {
		return
	}
	values := []string{}
	if fp.userConfig != nil {
		var err error
		values, err = fp.userConfig.GetStringSlice(fp.historyKey)
		if err != nil {
			logger := logging.DefaultLogger()
			logger.Warn("履歴読込に失敗しました")
			values = []string{}
		}
	}

	if fp.historyDialog != nil {
		if fp.historyDialog.IsDisposed() {
			fp.historyDialog = nil
			fp.historyListBox = nil
		} else {
			fp.historyListBox.SetModel(values)
			fp.historyDialog.Show()
			return
		}
	}

	dlg := new(walk.Dialog)
	lb := new(walk.ListBox)
	push := new(walk.PushButton)
	var parent walk.Form
	if fp.window != nil {
		parent = fp.window
	} else {
		parent = walk.App().ActiveForm()
	}
	if parent == nil {
		return
	}
	if err := (declarative.Dialog{
		AssignTo: &dlg,
		Title:    fp.t("履歴"),
		MinSize:  declarative.Size{Width: 800, Height: 400},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo: &lb,
				Model:    values,
				MinSize:  declarative.Size{Width: 800, Height: 400},
				OnItemActivated: func() {
					idx := lb.CurrentIndex()
					if idx < 0 || idx >= len(values) {
						return
					}
					push.SetEnabled(true)
					fp.applyPaths([]string{values[idx]}, true)
					dlg.Accept()
				},
			}, declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &push,
						Text:     fp.t("OK"),
						Enabled:  true,
						OnClicked: func() {
							dlg.Accept()
						},
					},
					declarative.PushButton{
						Text: fp.t("キャンセル"),
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}).Create(parent); err != nil {
		return
	}

	fp.historyDialog = dlg
	fp.historyListBox = lb
	fp.historyDialog.Disposing().Attach(func() {
		fp.historyDialog = nil
		fp.historyListBox = nil
	})
	push.SetEnabled(true)
	fp.historyDialog.Show()
}

// applyPaths はパス更新処理を共通化する。
func (fp *FolderPicker) applyPaths(paths []string, allowSame bool) {
	cleaned := cleanPaths(paths)
	if len(cleaned) == 0 {
		return
	}
	if !allowSame && sameStringSlice(cleaned, fp.prevPaths) {
		return
	}
	fp.prevPaths = cleaned
	if fp.pathEdit != nil {
		fp.pathEdit.SetText(cleaned[0])
	}
	if fp.onPathsChanged != nil {
		fp.onPathsChanged(fp.window, cleaned)
	}
	fp.saveHistoryIfNeeded(cleaned)
}

// saveHistoryIfNeeded は履歴保存が可能な場合に保存する。
func (fp *FolderPicker) saveHistoryIfNeeded(paths []string) {
	if fp.historyKey == "" || fp.userConfig == nil {
		return
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil {
		return
	}
	values = append(append([]string{}, paths...), values...)
	values = dedupe(values)
	if err := fp.userConfig.SetStringSlice(fp.historyKey, values, 50); err != nil {
		logger := logging.DefaultLogger()
		logger.Warn("履歴保存に失敗しました: %s", err.Error())
	}
}

// t は翻訳文字列を返す。
func (fp *FolderPicker) t(text string) string {
	return i18n.TranslateOrMark(fp.translator, text)
}

// cleanPaths はパス一覧を正規化し、フォルダのみを返す。
func cleanPaths(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		cleaned := cleanPath(path)
		if cleaned == "" {
			continue
		}
		info, err := os.Stat(cleaned)
		if err != nil || info == nil || !info.IsDir() {
			continue
		}
		result = append(result, cleaned)
	}
	result = dedupe(result)
	sortPaths(result)
	return result
}

// cleanPath は入力パスを正規化する。
func cleanPath(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.Clean(path)
	path = strings.Trim(path, "\"")
	path = strings.Trim(path, "'")
	path = strings.TrimSpace(path)
	path = strings.Trim(path, ".")
	return path
}

// dedupe は重複を排除したスライスを返す。
func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// sameStringSlice は同一内容か判定する。
func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// sortPaths はパス配列を昇順で整列する。
func sortPaths(values []string) {
	if len(values) == 0 {
		return
	}
	sort.SliceStable(values, func(i, j int) bool {
		return strings.ToLower(values[i]) < strings.ToLower(values[j])
	})
}
