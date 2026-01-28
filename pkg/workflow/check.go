// 指示: miu200521358
package workflow

import (
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/merrors"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"

	"github.com/miu200521358/mu_motion_viewer/pkg/ok_ng_rules"
)

// indexedName は表示順の並び替えに使う一時構造体。
type indexedName struct {
	Name  string
	Index int
}

// CheckExists はモーション内のボーン/モーフがモデルに存在するか判定する。
func CheckExists(modelData *model.PmxModel, motionData *motion.VmdMotion) (CheckResult, error) {
	result := CheckResult{}
	if motionData == nil {
		return result, nil
	}

	activeBoneNames := collectActiveBoneNames(motionData)
	okBoneEntries := make([]indexedName, 0, len(activeBoneNames))
	ngBoneNames := make([]string, 0, len(activeBoneNames))
	for _, name := range activeBoneNames {
		bone, ok, err := resolveBone(modelData, name)
		if err != nil {
			return CheckResult{}, err
		}
		if ok && bone != nil && ok_ng_rules.IsExactMatch(name, bone.Name()) {
			okBoneEntries = append(okBoneEntries, indexedName{Name: name, Index: bone.Index()})
			continue
		}
		ngBoneNames = append(ngBoneNames, name)
	}

	activeMorphNames := collectActiveMorphNames(motionData)
	okMorphEntries := make([]indexedName, 0, len(activeMorphNames))
	ngMorphNames := make([]string, 0, len(activeMorphNames))
	for _, name := range activeMorphNames {
		morph, ok, err := resolveMorph(modelData, name)
		if err != nil {
			return CheckResult{}, err
		}
		if ok && morph != nil && ok_ng_rules.IsExactMatch(name, morph.Name()) {
			okMorphEntries = append(okMorphEntries, indexedName{Name: name, Index: morph.Index()})
			continue
		}
		ngMorphNames = append(ngMorphNames, name)
	}

	result.OkBones = sortNamesByIndex(okBoneEntries)
	result.OkMorphs = sortNamesByIndex(okMorphEntries)
	result.NgBones = sortNamesByName(ngBoneNames)
	result.NgMorphs = sortNamesByName(ngMorphNames)
	return result, nil
}

// collectActiveBoneNames は有効なボーン名を列挙する。
func collectActiveBoneNames(motionData *motion.VmdMotion) []string {
	if motionData == nil || motionData.BoneFrames == nil {
		return nil
	}

	names := motionData.BoneFrames.Names()
	out := make([]string, 0, len(names))
	for _, name := range names {
		frames := motionData.BoneFrames.Get(name)
		if frames == nil || !frames.ContainsActive() {
			continue
		}
		out = append(out, name)
	}
	return out
}

// collectActiveMorphNames は有効なモーフ名を列挙する。
func collectActiveMorphNames(motionData *motion.VmdMotion) []string {
	if motionData == nil || motionData.MorphFrames == nil {
		return nil
	}

	names := motionData.MorphFrames.Names()
	out := make([]string, 0, len(names))
	for _, name := range names {
		frames := motionData.MorphFrames.Get(name)
		if frames == nil || !frames.ContainsActive() {
			continue
		}
		out = append(out, name)
	}
	return out
}

// resolveBone はモデル内のボーンを取得し、存在有無を返す。
func resolveBone(modelData *model.PmxModel, name string) (*model.Bone, bool, error) {
	if modelData == nil || modelData.Bones == nil {
		return nil, false, nil
	}
	bone, err := modelData.Bones.GetByName(name)
	if err != nil {
		if merrors.IsNameNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return bone, true, nil
}

// resolveMorph はモデル内のモーフを取得し、存在有無を返す。
func resolveMorph(modelData *model.PmxModel, name string) (*model.Morph, bool, error) {
	if modelData == nil || modelData.Morphs == nil {
		return nil, false, nil
	}
	morph, err := modelData.Morphs.GetByName(name)
	if err != nil {
		if merrors.IsNameNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return morph, true, nil
}

// sortNamesByIndex はインデックス順に名前を並べ替える。
func sortNamesByIndex(entries []indexedName) []string {
	if len(entries) == 0 {
		return nil
	}
	copyEntries := append([]indexedName(nil), entries...)
	sort.Slice(copyEntries, func(i, j int) bool {
		return copyEntries[i].Index < copyEntries[j].Index
	})
	out := make([]string, len(copyEntries))
	for i, entry := range copyEntries {
		out[i] = entry.Name
	}
	return out
}

// sortNamesByName は名前順に並べ替える。
func sortNamesByName(names []string) []string {
	if len(names) == 0 {
		return nil
	}
	out := append([]string(nil), names...)
	sort.Strings(out)
	return out
}
