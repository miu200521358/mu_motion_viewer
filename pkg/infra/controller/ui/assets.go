//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"bytes"
	"embed"
	"image/png"

	"github.com/miu200521358/walk/pkg/walk"
)

//go:embed assets/content_copy.png
var contentCopyPNG embed.FS

// loadCopyBitmap はコピーアイコンのビットマップを生成する。
func loadCopyBitmap() (*walk.Bitmap, error) {
	data, err := contentCopyPNG.ReadFile("assets/content_copy.png")
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return walk.NewBitmapFromImage(img)
}
