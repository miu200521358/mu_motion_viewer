// 指示: miu200521358
package usecase

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	portio "github.com/miu200521358/mu_motion_viewer/pkg/usecase/port/io"

	"github.com/miu200521358/mu_motion_viewer/pkg/workflow"
)

// SafeMotionSaveRequest は安全モーション保存の入力を表す。
type SafeMotionSaveRequest struct {
	Motion       *motion.VmdMotion
	FallbackPath string
	Writer       portio.IFileWriter
	SaveOptions  portio.SaveOptions
}

// SafeMotionSaveResult は安全モーション保存の結果を表す。
type SafeMotionSaveResult struct {
	BasePath string
	SafePath string
}

// BuildSafeMotion はIKフレームを空にしたモーションを複製する。
func BuildSafeMotion(source *motion.VmdMotion) (*motion.VmdMotion, error) {
	return workflow.BuildSafeMotion(source)
}

// SaveSafeMotion は安全モーションを生成して保存する。
func SaveSafeMotion(request SafeMotionSaveRequest) (*SafeMotionSaveResult, error) {
	result := &SafeMotionSaveResult{}
	if request.Motion == nil {
		return result, nil
	}
	basePath := request.Motion.Path()
	if basePath == "" {
		basePath = request.FallbackPath
	}
	result.BasePath = basePath
	if basePath == "" {
		return result, nil
	}
	if request.Writer == nil {
		return result, fmt.Errorf("保存リポジトリがありません")
	}

	safeMotion, err := BuildSafeMotion(request.Motion)
	if err != nil {
		return result, err
	}
	if safeMotion == nil {
		return result, nil
	}

	safePath := buildSafeMotionPath(basePath)
	result.SafePath = safePath
	if safePath == "" {
		return result, nil
	}
	if err := request.Writer.Save(safePath, safeMotion, request.SaveOptions); err != nil {
		return result, err
	}
	return result, nil
}

// buildSafeMotionPath は安全モーションの保存先パスを生成する。
func buildSafeMotionPath(path string) string {
	if path == "" {
		return ""
	}
	dir, base := filepath.Split(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if ext == "" {
		ext = ".vmd"
	}
	if name == "" {
		return filepath.Join(dir, "_safe"+ext)
	}
	return filepath.Join(dir, name+"_safe"+ext)
}
