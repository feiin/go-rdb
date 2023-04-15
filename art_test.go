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

func TestInsert2(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("test"), "hello world")
	tree.Insert([]byte("test2"), "hello world2")

	t.Logf("insert tree.innerNode:%+v", tree.root.innerNode)
	t.Logf("prefix :%s", string(tree.root.innerNode.prefix))

	if tree.root.innerNode.prefixLen != 4 {
		t.Errorf("invalid prefixLen:%d", tree.root.innerNode.prefixLen)
	}

	v1 := tree.Search([]byte("test"))
	t.Logf("insert value1:%+v", v1)

	if v1.(string) != "hello world" {
		t.Error("invalid value1")
	}

	v2 := tree.Search([]byte("test2"))
	t.Logf("insert value2:%+v", v2)

	if v2.(string) != "hello world2" {
		t.Error("invalid value2")
	}

}
