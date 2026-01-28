//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

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
