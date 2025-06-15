package repository

import "testing"

func TestReadGraph(t *testing.T) {
	// テスト用のパッケージを読み込む
	nodes, err := ReadGraph("fmt")
	if err != nil {
		t.Fatalf("Failed to read graph: %v", err)
	}

	// ノード数を確認
	if len(nodes) == 0 {
		t.Error("Expected graph to have nodes, but it is empty")
	}

	// 少なくとも1つのノードに子要素があることを確認
	hasChildren := false
	for _, node := range nodes {
		if len(node.GetChildren()) > 0 {
			hasChildren = true
			break
		}
	}
	
	if !hasChildren {
		t.Error("Expected at least one node to have children (edges)")
	}
}
