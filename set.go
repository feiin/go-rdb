package rdb

type set map[interface{}]struct{}

func newSet() set {
	return make(map[interface{}]struct{})
}

func (s set) Add(items ...interface{}) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

func (s set) Card() int {
	return len(s)
}

func (s set) Remove(item ...interface{}) {
	for _, item := range item {
		delete(s, item)
	}
}
