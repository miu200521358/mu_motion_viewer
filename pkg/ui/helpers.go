//go:build windows
// +build windows

package ui

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// translate は翻訳済み文言を返す。
func translate(translator i18n.II18n, key string) string {
	if translator == nil || !translator.IsReady() {
		return "●●" + key + "●●"
	}
	return translator.T(key)
}

// logInfoLine は情報ログを1行として出力する。
func logInfoLine(logger logging.ILogger, message string, params ...any) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if lineLogger, ok := logger.(interface {
		InfoLine(msg string, params ...any)
	}); ok {
		lineLogger.InfoLine(message, params...)
		return
	}
	logger.Info(message, params...)
}

// logErrorWithTitle はタイトル付きのエラーログを出力する。
func logErrorWithTitle(logger logging.ILogger, title string, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if err == nil {
		if titled, ok := logger.(interface {
			ErrorTitle(title string, err error, msg string, params ...any)
		}); ok {
			titled.ErrorTitle(title, nil, "")
			return
		}
		logger.Error(title)
		return
	}
	if titled, ok := logger.(interface {
		ErrorTitle(title string, err error, msg string, params ...any)
	}); ok {
		titled.ErrorTitle(title, err, "")
		return
	}
	logger.Error("%s: %s", title, err.Error())
}
