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
	initialDirPath    string
	translator        i18n.II18n
	userConfig        iCommonUserConfig
	pathEdit          *walk.LineEdit
	openPushButton    *walk.PushButton
	historyPushButton *walk.PushButton
	historyDialog     *walk.Dialog
	historyListBox    *walk.ListBox
	prevPaths         []string
	onPathsChanged    func(*controller.ControlWindow, []string)
}

// iCommonUserConfig は履歴保存に使うI/Fを表す。
type iCommonUserConfig interface {
	GetStringSlice(key string) ([]string, error)
	SetStringSlice(key string, values []string, limit int) error
}

// NewFolderPicker はFolderPickerを生成する。
func NewFolderPicker(userConfig iCommonUserConfig, translator i18n.II18n, historyKey string, title string, tooltip string, onPathsChanged func(*controller.ControlWindow, []string)) *FolderPicker {
	return &FolderPicker{
		title:          title,
		tooltip:        tooltip,
		historyKey:     historyKey,
		onPathsChanged: onPathsChanged,
		userConfig:     userConfig,
		translator:     translator,
	}
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

// SetPaths はパス一覧を設定する。
func (fp *FolderPicker) SetPaths(paths []string) {
	fp.applyPaths(paths, true)
}

// Widgets はUI構成を返す。
func (fp *FolderPicker) Widgets() declarative.Composite {
	inputWidgets := []declarative.Widget{
		declarative.LineEdit{
			AssignTo:    &fp.pathEdit,
			ToolTipText: fp.tooltip,
			OnTextChanged: func() {
				fp.handleTextChanged(fp.pathEdit.Text())
			},
			OnEditingFinished: func() {
				fp.handleTextConfirmed(fp.pathEdit.Text())
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
				fp.onOpenClicked()
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
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.TextLabel{
						Text:        fp.title,
						ToolTipText: fp.tooltip,
					},
				},
			},
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: inputWidgets,
			},
		},
	}
}

// onOpenClicked は開くボタンの処理を行う。
func (fp *FolderPicker) onOpenClicked() {
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
	fp.handleTextConfirmed(fd.FilePath)
}

// handleTextChanged は入力途中の変更を処理する。
func (fp *FolderPicker) handleTextChanged(text string) {
	paths := splitPaths(text)
	fp.applyPaths(paths, false)
}

// handleTextConfirmed は確定時の変更を処理する。
func (fp *FolderPicker) handleTextConfirmed(text string) {
	paths := splitPaths(text)
	fp.applyPaths(paths, true)
}

// handleDropFiles はドロップされたフォルダ一覧を反映する。
func (fp *FolderPicker) handleDropFiles(files []string) {
	if fp == nil || len(files) == 0 {
		return
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		cleaned := cleanPath(file)
		if cleaned == "" {
			continue
		}
		info, err := os.Stat(cleaned)
		if err != nil || !info.IsDir() {
			continue
		}
		paths = append(paths, cleaned)
	}
	if len(paths) == 0 {
		return
	}
	sort.Strings(paths)
	fp.applyPaths(paths, true)
}

// applyPaths はパス更新処理を共通化する。
func (fp *FolderPicker) applyPaths(paths []string, allowSame bool) {
	cleaned := normalizePaths(paths)
	if len(cleaned) == 0 {
		return
	}
	if !allowSame && samePaths(cleaned, fp.prevPaths) {
		return
	}

	valid := make([]string, 0, len(cleaned))
	for _, path := range cleaned {
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			continue
		}
		valid = append(valid, path)
	}
	if len(valid) == 0 {
		return
	}

	fp.prevPaths = append([]string{}, valid...)
	if fp.pathEdit != nil {
		fp.pathEdit.SetText(strings.Join(valid, ";"))
	}
	if fp.onPathsChanged != nil {
		fp.onPathsChanged(fp.window, valid)
	}
	fp.saveHistoryIfNeeded(valid)
}

// saveHistoryIfNeeded は履歴保存が可能な場合に保存する。
func (fp *FolderPicker) saveHistoryIfNeeded(paths []string) {
	if fp.historyKey == "" || fp.userConfig == nil || len(paths) == 0 {
		return
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil {
		return
	}
	merged := append([]string{}, paths...)
	merged = append(merged, values...)
	merged = dedupe(merged)
	if err := fp.userConfig.SetStringSlice(fp.historyKey, merged, 50); err != nil {
		logger := logging.DefaultLogger()
		logger.Warn("履歴保存に失敗しました: %s", err.Error())
	}
}

// resolveInitialDir は初期ディレクトリを決定する。
func (fp *FolderPicker) resolveInitialDir() string {
	if fp.pathEdit != nil {
		current := splitPaths(fp.pathEdit.Text())
		if len(current) > 0 {
			return current[0]
		}
	}
	if fp.historyKey == "" || fp.userConfig == nil {
		return fp.initialDirPath
	}
	values, err := fp.userConfig.GetStringSlice(fp.historyKey)
	if err != nil || len(values) == 0 {
		return fp.initialDirPath
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
					fp.handleTextConfirmed(values[idx])
					dlg.Accept()
				},
			},
			declarative.Composite{
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

// t は翻訳済み文言を返す。
func (fp *FolderPicker) t(key string) string {
	if fp == nil || fp.translator == nil || !fp.translator.IsReady() {
		return "●●" + key + "●●"
	}
	return fp.translator.T(key)
}

// splitPaths は入力文字列を分割する。
func splitPaths(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	return strings.FieldsFunc(input, func(r rune) bool {
		switch r {
		case ';', '\n', '\r', '\t':
			return true
		default:
			return false
		}
	})
}

// normalizePaths は入力パスを正規化して重複排除する。
func normalizePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		cleaned := cleanPath(p)
		if cleaned == "" {
			continue
		}
		if _, ok := seen[cleaned]; ok {
			continue
		}
		seen[cleaned] = struct{}{}
		result = append(result, cleaned)
	}
	return result
}

// cleanPath は入力パスを正規化して返す。
func cleanPath(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.Trim(trimmed, "\"")
	trimmed = strings.Trim(trimmed, "'")
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}

// samePaths はパス一覧が同一か判定する。
func samePaths(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}
	return true
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
