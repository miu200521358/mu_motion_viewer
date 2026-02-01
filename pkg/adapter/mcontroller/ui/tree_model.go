//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/walk/pkg/walk"
)

const (
	modelExtPmx = ".pmx"
	modelExtPmd = ".pmd"
	modelExtX   = ".x"
)

// TreeNode はツリー表示用のノードを表す。
type TreeNode struct {
	name     string
	fullPath string
	parent   *TreeNode
	children []*TreeNode
	isDir    bool
}

// NewTreeNode はTreeNodeを生成する。
func NewTreeNode(name string, fullPath string, parent *TreeNode, isDir bool) *TreeNode {
	return &TreeNode{
		name:     name,
		fullPath: fullPath,
		parent:   parent,
		isDir:    isDir,
	}
}

// Text はツリー表示用のラベルを返す。
func (n *TreeNode) Text() string {
	if n == nil {
		return ""
	}
	return n.name
}

// Parent は親ノードを返す。
func (n *TreeNode) Parent() walk.TreeItem {
	if n == nil {
		return nil
	}
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// ChildCount は子ノード数を返す。
func (n *TreeNode) ChildCount() int {
	if n == nil {
		return 0
	}
	return len(n.children)
}

// ChildAt は指定インデックスの子ノードを返す。
func (n *TreeNode) ChildAt(index int) walk.TreeItem {
	if n == nil {
		return nil
	}
	if index < 0 || index >= len(n.children) {
		return nil
	}
	return n.children[index]
}

// HasChild は子ノードが存在するか判定する。
func (n *TreeNode) HasChild() bool {
	return n != nil && len(n.children) > 0
}

// Path はノードのフルパスを返す。
func (n *TreeNode) Path() string {
	if n == nil {
		return ""
	}
	return n.fullPath
}

// IsDir はディレクトリノードか判定する。
func (n *TreeNode) IsDir() bool {
	if n == nil {
		return false
	}
	return n.isDir
}

// addChild は子ノードを追加する。
func (n *TreeNode) addChild(child *TreeNode) {
	if n == nil || child == nil {
		return
	}
	n.children = append(n.children, child)
}

// sortChildren は子ノードを名前順で並べ替える。
func (n *TreeNode) sortChildren() {
	if n == nil {
		return
	}
	if len(n.children) > 1 {
		sort.SliceStable(n.children, func(i, j int) bool {
			return strings.ToLower(n.children[i].name) < strings.ToLower(n.children[j].name)
		})
	}
	for _, child := range n.children {
		child.sortChildren()
	}
}

// TreeModel はツリービュー表示用のモデルを表す。
type TreeModel struct {
	walk.TreeModelBase
	roots []*TreeNode
}

// NewTreeModel はTreeModelを生成する。
func NewTreeModel() *TreeModel {
	return &TreeModel{}
}

// RootCount はルートノード数を返す。
func (m *TreeModel) RootCount() int {
	if m == nil {
		return 0
	}
	return len(m.roots)
}

// RootAt は指定インデックスのルートノードを返す。
func (m *TreeModel) RootAt(index int) walk.TreeItem {
	if m == nil {
		return nil
	}
	if index < 0 || index >= len(m.roots) {
		return nil
	}
	return m.roots[index]
}

// SetRoots はルートパス一覧からツリー構造を再構築する。
func (m *TreeModel) SetRoots(paths []string) error {
	if m == nil {
		return errors.New("tree model is nil")
	}
	roots, err := buildRoots(paths)
	m.roots = roots
	m.PublishItemsReset(nil)
	return err
}

// buildRoots はルート配下のモデルファイルからツリーを構成する。
func buildRoots(paths []string) ([]*TreeNode, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	roots := make([]*TreeNode, 0, len(paths))
	var errs []error
	for _, rootPath := range paths {
		if rootPath == "" {
			continue
		}
		rootNode, err := buildRootNode(rootPath)
		if err != nil {
			errs = append(errs, err)
		}
		if rootNode != nil {
			roots = append(roots, rootNode)
		}
	}
	if len(roots) > 1 {
		sort.SliceStable(roots, func(i, j int) bool {
			return strings.ToLower(roots[i].fullPath) < strings.ToLower(roots[j].fullPath)
		})
	}
	return roots, errors.Join(errs...)
}

// buildRootNode は指定ルートのツリーノードを生成する。
func buildRootNode(rootPath string) (*TreeNode, error) {
	modelPaths, err := collectModelPaths(rootPath)
	if len(modelPaths) == 0 {
		return nil, err
	}
	rootLabel := fmt.Sprintf("【%s】", rootPath)
	rootNode := NewTreeNode(rootLabel, rootPath, nil, true)

	dirNodes := map[string]*TreeNode{}
	dirNodes[strings.ToLower(rootPath)] = rootNode

	var errs []error
	if err != nil {
		errs = append(errs, err)
	}

	for _, modelPath := range modelPaths {
		rel, relErr := filepath.Rel(rootPath, modelPath)
		if relErr != nil {
			errs = append(errs, relErr)
			continue
		}
		parts := splitPath(rel)
		if len(parts) == 0 {
			continue
		}
		current := rootNode
		currentPath := rootPath
		for i, part := range parts {
			if part == "" {
				continue
			}
			if i == len(parts)-1 {
				fileNode := NewTreeNode(part, modelPath, current, false)
				current.addChild(fileNode)
				continue
			}
			currentPath = filepath.Join(currentPath, part)
			key := strings.ToLower(currentPath)
			child := dirNodes[key]
			if child == nil {
				child = NewTreeNode(part, currentPath, current, true)
				dirNodes[key] = child
				current.addChild(child)
			}
			current = child
		}
	}

	rootNode.sortChildren()
	return rootNode, errors.Join(errs...)
}

// collectModelPaths はモデルファイルのパスを収集する。
func collectModelPaths(rootPath string) ([]string, error) {
	if rootPath == "" {
		return nil, nil
	}
	var paths []string
	var errs []error
	walkErr := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, err)
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry == nil || entry.IsDir() {
			return nil
		}
		if isModelFile(path) {
			paths = append(paths, path)
		}
		return nil
	})
	if walkErr != nil {
		errs = append(errs, walkErr)
	}
	return paths, errors.Join(errs...)
}

// isModelFile はモデル拡張子か判定する。
func isModelFile(path string) bool {
	if path == "" {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case modelExtPmx, modelExtPmd, modelExtX:
		return true
	default:
		return false
	}
}

// splitPath はOS依存区切りで分割する。
func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, string(os.PathSeparator))
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}
