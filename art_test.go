package cache

import (
	"fmt"
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
func TestInsertExpand(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("1key"), "v1")
	tree.Insert([]byte("2key"), "v2")
	tree.Insert([]byte("3key"), "v3")
	tree.Insert([]byte("4key"), "v4")

	if tree.root.innerNode.nodeType != Node4 {
		t.Errorf("invalid nodeType:%d", tree.root.innerNode.nodeType)
	}
	tree.Insert([]byte("5key"), "v5")
	if tree.root.innerNode.nodeType != Node16 {
		t.Errorf("invalid nodeType:%d", tree.root.innerNode.nodeType)
	}
	t.Logf("insert tree.innerNode:%+v", tree.root.innerNode)

	v1 := tree.Search([]byte("1key"))
	t.Logf("insert value1:%+v", v1)

	if v1.(string) != "v1" {
		t.Error("invalid value1")
	}

	v2 := tree.Search([]byte("2key"))
	t.Logf("insert value2:%+v", v2)

	if v2.(string) != "v2" {
		t.Error("invalid value2")
	}

	v3 := tree.Search([]byte("3key"))
	t.Logf("insert value3:%+v", v3)

	if v3.(string) != "v3" {
		t.Error("invalid value3")
	}

	v4 := tree.Search([]byte("4key"))
	t.Logf("insert value4:%+v", v4)

	if v4.(string) != "v4" {
		t.Error("invalid value4")
	}

	v5 := tree.Search([]byte("5key"))
	t.Logf("insert value5:%+v", v5)

	if v5.(string) != "v5" {
		t.Error("invalid value5")
	}

}

func TestShrink2Node4(t *testing.T) {
	tree := NewTree()
	tree.Insert([]byte("1key"), "v1")
	tree.Insert([]byte("2key"), "v2")
	tree.Insert([]byte("3key"), "v3")
	tree.Insert([]byte("4key"), "v4")
	tree.Insert([]byte("5key"), "v5")

	if tree.root.innerNode.nodeType != Node16 {
		t.Errorf("invalid nodeType:%d", tree.root.innerNode.nodeType)
	}

	tree.Delete([]byte("5key"))

	if tree.root.innerNode.nodeType != Node4 {
		t.Errorf("invalid nodeType:%d", tree.root.innerNode.nodeType)
	}

	v1 := tree.Search([]byte("1key"))
	t.Logf("insert value1:%+v", v1)

	if v1.(string) != "v1" {
		t.Error("invalid value1")
	}

	v5 := tree.Search([]byte("5key"))
	t.Logf("insert value5:%+v", v5)
}

func TestShrink2Node16(t *testing.T) {
	n := newNode48()

	for i := 0; i < n.innerNode.minSize(); i++ {
		n.innerNode.addChild(byte(i), newNode16())
	}

	if n.innerNode.nodeType != Node48 {
		t.Errorf("invalid nodeType:%d", n.innerNode.nodeType)
	}

	n.deleteChild(byte(0))

	if n.innerNode.nodeType != Node16 {
		t.Errorf("shrink invalid nodeType:%d", n.innerNode.nodeType)
	}
}

func TestShrink2Node48(t *testing.T) {
	n := newNode256()

	for i := 0; i < n.innerNode.minSize(); i++ {
		n.innerNode.addChild(byte(i), newNode16())
	}

	if n.innerNode.nodeType != Node256 {
		t.Errorf("invalid nodeType:%d", n.innerNode.nodeType)
	}

	n.deleteChild(byte(0))

	if n.innerNode.nodeType != Node48 {
		t.Errorf("shrink invalid nodeType:%d", n.innerNode.nodeType)
	}

	fmt.Printf("result innerNode:%+v", n.innerNode)

}
