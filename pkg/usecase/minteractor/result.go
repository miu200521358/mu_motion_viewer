// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult = commonusecase.ModelLoadResult

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult = commonusecase.MotionLoadResult

// CheckResult はOK/NG判定の結果一覧を表す。
type CheckResult struct {
	OkBones  []string
	OkMorphs []string
	NgBones  []string
	NgMorphs []string
}

// ExtractModelData は読み込み結果からモデルを取り出す。
func ExtractModelData(result *ModelLoadResult) *model.PmxModel {
	if result == nil {
		return nil
	}
	return result.Model
}

// ExtractMotionData は読み込み結果からモーションと最大フレームを取り出す。
func ExtractMotionData(result *MotionLoadResult) (*motion.VmdMotion, motion.Frame) {
	if result == nil {
		return nil, 0
	}
	return result.Motion, result.MaxFrame
}
