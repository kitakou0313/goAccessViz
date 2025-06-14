package repository

import "testing"

func TestReadGraph(t *testing.T) {
	// テスト用のグラフを読み込む
	graph, err := ReadGraph("testdata/test_graph.json")
	if err != nil {
		t.Fatalf("Failed to read graph: %v", err)
	}

	// グラフのノード数を確認
	if len(graph.Nodes()) == 0 {
		t.Error("Expected graph to have nodes, but it is empty")
	}

	// グラフのエッジ数を確認
	if len(graph.Edges()) == 0 {
		t.Error("Expected graph to have edges, but it is empty")
	}

}
