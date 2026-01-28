//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/mu_motion_viewer/pkg/adapter/mpresenter/messages"
)

// NewMenuItems はメニュー項目を生成する。
func NewMenuItems(translator i18n.II18n, logger logging.ILogger) []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text: i18n.TranslateOrMark(translator, messages.HelpUsage),
			OnTriggered: func() {
				logInfoLine(logger, messages.HelpUsage)
			},
		},
	}
}
