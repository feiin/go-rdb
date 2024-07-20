package rdb

import "testing"

func TestSetAdd(t *testing.T) {
	set := newSet()
	set.Add("test", "test2")
	count := set.Card()
	if count != 2 {
		t.Error("invalid count")
	}
	t.Logf("count: %d", count)
}

func TestSetRemove(t *testing.T) {
	set := newSet()
	set.Add("test", "test2")
	count := set.Card()
	if count != 2 {
		t.Error("invalid count")
	}
	t.Logf("count: %d", count)
	set.Remove("test")
	count = set.Card()
	if count != 1 {
		t.Error("invalid count")
	}
	t.Logf("count: %d", count)
}
