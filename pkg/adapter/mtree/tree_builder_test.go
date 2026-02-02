// 指示: miu200521358
package mtree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildModelTree_PrunesAndSorts(t *testing.T) {
	root := t.TempDir()
	createDir(t, root, "empty")
	createDir(t, root, "alpha")
	createDir(t, root, "zeta")
	createDir(t, root, "nested")
	createDir(t, filepath.Join(root, "nested"), "sub")

	createFile(t, root, "ignore.txt")
	createFile(t, root, "A.pmx")
	createFile(t, root, "b.PMD")
	createFile(t, filepath.Join(root, "alpha"), "a.pmx")
	createFile(t, filepath.Join(root, "zeta"), "z.pmx")
	createFile(t, filepath.Join(root, "nested", "sub"), "c.x")

	roots, err := BuildModelTree([]string{root})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}
	rootNode := roots[0]
	absRoot, _ := filepath.Abs(root)
	expectedName := fmt.Sprintf(rootLabelFormat, absRoot)
	if rootNode.Name != expectedName {
		t.Fatalf("expected root name %q, got %q", expectedName, rootNode.Name)
	}
	if rootNode.Path != absRoot {
		t.Fatalf("expected root path %q, got %q", absRoot, rootNode.Path)
	}
	if !isSortedByName(rootNode.Children) {
		t.Fatalf("root children are not sorted: %v", collectNames(rootNode.Children))
	}
	if containsName(rootNode.Children, "empty") {
		t.Fatalf("unexpected empty folder in children: %v", collectNames(rootNode.Children))
	}
	if containsName(rootNode.Children, "ignore.txt") {
		t.Fatalf("unexpected ignored file in children: %v", collectNames(rootNode.Children))
	}
	if !containsName(rootNode.Children, "alpha") || !containsName(rootNode.Children, "zeta") {
		t.Fatalf("expected model folders in children: %v", collectNames(rootNode.Children))
	}
}

func TestBuildModelTree_RootOrder(t *testing.T) {
	base := t.TempDir()
	rootB := filepath.Join(base, "broot")
	rootA := filepath.Join(base, "aroot")
	createDir(t, base, "broot")
	createDir(t, base, "aroot")
	createFile(t, rootB, "b.pmx")
	createFile(t, rootA, "a.pmx")

	roots, err := BuildModelTree([]string{rootB, rootA})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}
	if roots[0].Path != rootA {
		t.Fatalf("expected root order %q first, got %q", rootA, roots[0].Path)
	}
	if roots[1].Path != rootB {
		t.Fatalf("expected root order %q second, got %q", rootB, roots[1].Path)
	}
}

// createDir はテスト用のディレクトリを作成する。
func createDir(t *testing.T, base, name string) {
	t.Helper()
	path := filepath.Join(base, name)
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
}

// createFile はテスト用の空ファイルを作成する。
func createFile(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
}

// collectNames はノード名一覧を返す。
func collectNames(nodes []*Node) []string {
	out := make([]string, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.Name)
	}
	return out
}

// containsName は指定名が含まれているか判定する。
func containsName(nodes []*Node, name string) bool {
	for _, node := range nodes {
		if node.Name == name {
			return true
		}
	}
	return false
}

// isSortedByName は名前が昇順で並んでいるか判定する。
func isSortedByName(nodes []*Node) bool {
	if len(nodes) <= 1 {
		return true
	}
	prev := strings.ToLower(nodes[0].Name)
	for _, node := range nodes[1:] {
		current := strings.ToLower(node.Name)
		if prev > current {
			return false
		}
		prev = current
	}
	return true
}
