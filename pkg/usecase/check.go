package usecase

import (
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

func CheckExists(model *pmx.PmxModel, motion *vmd.VmdMotion) (okBoneNames, okMorphNames, ngBoneNames, ngMorphNames []string) {
	okBoneNames = make([]string, 0)
	okMorphNames = make([]string, 0)
	ngBoneNames = make([]string, 0)
	ngMorphNames = make([]string, 0)

	if model == nil || motion == nil {
		return okBoneNames, okMorphNames, ngBoneNames, ngMorphNames
	}

	motion.BoneFrames.ForEach(func(boneName string, boneNameFrames *vmd.BoneNameFrames) {
		if boneNameFrames.ContainsActive() {
			if model.Bones.ContainsByName(boneName) {
				okBoneNames = append(okBoneNames, boneName)
			} else {
				ngBoneNames = append(ngBoneNames, boneName)
			}
		}
	})

	// OKボーンはモデルのボーンINDEX順にソート
	slices.SortFunc(okBoneNames, func(a, b string) int {
		aBone, err := model.Bones.GetByName(a)
		if err != nil {
			return -1
		}
		bBone, err := model.Bones.GetByName(b)
		if err != nil {
			return 1
		}
		if aBone.Index() < bBone.Index() {
			return -1
		} else if aBone.Index() > bBone.Index() {
			return 1
		}
		return 0
	})

	// NGボーンは名前順にソート
	sort.Strings(ngBoneNames)

	motion.MorphFrames.ForEach(func(morphName string, morphNameFrames *vmd.MorphNameFrames) {
		if morphNameFrames.ContainsActive() {
			if model.Morphs.ContainsByName(morphName) {
				okMorphNames = append(okMorphNames, morphName)
			} else {
				ngMorphNames = append(ngMorphNames, morphName)
			}
		}
	})

	// OKモーフはモデルのモーフINDEX順にソート
	slices.SortFunc(okMorphNames, func(a, b string) int {
		aMorph, err := model.Morphs.GetByName(a)
		if err != nil {
			return -1
		}
		bMorph, err := model.Morphs.GetByName(b)
		if err != nil {
			return 1
		}
		if aMorph.Index() < bMorph.Index() {
			return -1
		} else if aMorph.Index() > bMorph.Index() {
			return 1
		}
		return 0
	})

	// NGモーフは名前順にソート
	sort.Strings(ngMorphNames)

	return okBoneNames, okMorphNames, ngBoneNames, ngMorphNames
}
