package cache

import (
	"testing"
)

func TestInsertOne(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("test"), "hello world")
	if tree.root == nil {
		t.Error("invalid insert")
	}
	if tree.size != 1 {
		t.Error("invalid size")
	}
	t.Logf("insert tree.innerNode:%+v", tree.root)

	if !tree.root.IsLeaf() {
		t.Error("invalid error")
	}

	value := tree.Search([]byte("test"))
	if value == nil {
		t.Error("invalid search")
	}
	if value.(string) != "hello world" {
		t.Error("invalid value")
	}

	t.Logf("insert value:%+v", value)
}
