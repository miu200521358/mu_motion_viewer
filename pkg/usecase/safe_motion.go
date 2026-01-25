package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/motion"

	"github.com/miu200521358/mu_motion_viewer/pkg/workflow"
)

// BuildSafeMotion はIKフレームを空にしたモーションを複製する。
func BuildSafeMotion(source *motion.VmdMotion) (*motion.VmdMotion, error) {
	return workflow.BuildSafeMotion(source)
}
