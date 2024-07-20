package rdb

import "testing"

func TestSetGetString(t *testing.T) {

	SetString("test", "hello world")
	SetString("test2", "hello world2")

	v1, _ := GetString("test")
	if v1 != "hello world" {
		t.Error("invalid value1")
	}

	v2, _ := GetString("test2")
	if v2 != "hello world2" {
		t.Error("invalid value2")
	}
}

func TestSetGet(t *testing.T) {

	Set("test", "hello world")
	Set("test2", "hello world2")

	v1, _ := Get("test")
	if v1.(string) != "hello world" {
		t.Error("invalid value1")
	}

	v2, _ := Get("test2")
	if v2.(string) != "hello world2" {
		t.Error("invalid value2")
	}
}
