//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mpresenter/messages"
)

// NewMenuItems はメニュー項目を生成する。
func NewMenuItems(translator i18n.II18n, logger logging.ILogger) []declarative.MenuItem {
	return controller.BuildMenuItemsWithMessages(translator, logger, []controller.MenuMessageItem{
		{TitleKey: messages.HelpUsageTitle, MessageKey: messages.HelpUsage},
		{TitleKey: messages.LabelFolderPath, MessageKey: messages.LabelFolderPathTip},
		{TitleKey: messages.LabelMotionPath, MessageKey: messages.LabelMotionPathTip},
		{TitleKey: messages.LabelTreeView, MessageKey: messages.LabelTreeViewTip},
		{TitleKey: messages.LabelPathCopy, MessageKey: messages.LabelPathCopyTip},
	})
}
