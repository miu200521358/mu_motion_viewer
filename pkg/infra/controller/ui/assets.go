//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"sync"

	_ "embed"

	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/walk"
)

//go:embed assets/content_copy_24.png
var copyIconData []byte

var (
	copyIconOnce   sync.Once
	copyIconBitmap *walk.Bitmap
	copyIconErr    error
)

// loadCopyIcon はコピーアイコンのビットマップを取得する。
func loadCopyIcon(logger logging.ILogger) *walk.Bitmap {
	copyIconOnce.Do(func() {
		img, err := decodePng(copyIconData)
		if err != nil {
			copyIconErr = err
			return
		}
		copyIconBitmap, copyIconErr = walk.NewBitmapFromImageForDPI(img, 96)
	})
	if copyIconErr != nil {
		if logger == nil {
			logger = logging.DefaultLogger()
		}
		logger.Error("コピーアイコンの読み込みに失敗しました: %s", copyIconErr.Error())
	}
	return copyIconBitmap
}

// decodePng はPNGバイト列をデコードする。
func decodePng(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, errors.New("png data is empty")
	}
	return png.Decode(bytes.NewReader(data))
}
