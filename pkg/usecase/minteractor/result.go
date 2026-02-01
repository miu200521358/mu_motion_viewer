// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult struct {
	Model *model.PmxModel
}

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult struct {
	Motion   *motion.VmdMotion
	MaxFrame motion.Frame
}
