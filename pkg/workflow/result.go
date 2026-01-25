package workflow

// CheckResult はOK/NG判定の結果一覧を表す。
type CheckResult struct {
	OkBones  []string
	OkMorphs []string
	NgBones  []string
	NgMorphs []string
}
