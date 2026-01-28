// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"

	"github.com/miu200521358/mu_motion_viewer/pkg/workflow"
)

// CheckResult はOK/NG判定の結果一覧を表す。
type CheckResult = workflow.CheckResult

// CheckExists はモーション内のボーン/モーフがモデルに存在するか判定する。
func CheckExists(modelData *model.PmxModel, motionData *motion.VmdMotion) (CheckResult, error) {
	return workflow.CheckExists(modelData, motionData)
}
