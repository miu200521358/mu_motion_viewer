package workflow

import "github.com/miu200521358/mlib_go/pkg/domain/motion"

// BuildSafeMotion はIKフレームを空にしたモーションを複製する。
func BuildSafeMotion(source *motion.VmdMotion) (*motion.VmdMotion, error) {
	if source == nil {
		return nil, nil
	}
	copied, err := source.Copy()
	if err != nil {
		return nil, err
	}
	copied.IkFrames = motion.NewIkFrames()
	return &copied, nil
}
