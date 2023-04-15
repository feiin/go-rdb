package cache

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("test"), "hello world")
	if tree.root == nil {
		t.Error("invalid insert")
	}
	if tree.size != 1 {
		t.Error("invalid size")
	}
	t.Logf("insert tree.innerNode:%+v", tree.root.innerNode)

	if tree.root != nil && tree.root.innerNode.nodeType != Leaf {
		t.Error("invalid nodeType")
	}
	t.Logf("insert tree:%+v", tree)
}
