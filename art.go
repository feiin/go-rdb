package cache

import (
	"bytes"
)

const (
	Node4 = iota
	Node16
	Node48
	Node256
	Leaf

	Node4Min = 2
	Node4Max = 4

	Node16Min = Node4Max + 1
	Node16Max = 16

	Node48Min = Node16Max + 1
	Node48Max = 48

	Node256Min = Node48Max + 1
	Node256Max = 256

	MaxPrefixLen = 8

	nullIdx = -1
)

type LeafNode struct {
	key   []byte
	value interface{}
}

type InnerNode struct {
	nodeType int
	keys     []byte
	children []*ArtNode

	// 存储最大MaxPrefixLen数据
	prefix []byte
	// 实际前缀长度
	prefixLen int
	// 子节点数量
	size int
}

type ArtNode struct {
	leaf      *LeafNode
	innerNode *InnerNode
}

type ArtTree struct {
	root *ArtNode
	size int
}

func NewTree() *ArtTree {
	return &ArtTree{root: nil, size: 0}
}
func newNode4() *ArtNode {
	return &ArtNode{
		innerNode: &InnerNode{
			nodeType: Node4,
			keys:     make([]byte, Node4Max),
			children: make([]*ArtNode, Node4Max),
			prefix:   make([]byte, MaxPrefixLen),
		},
	}
}

func newNode16() *ArtNode {
	return &ArtNode{
		innerNode: &InnerNode{
			nodeType: Node16,
			keys:     make([]byte, Node16Max),
			children: make([]*ArtNode, Node16Max),
			prefix:   make([]byte, MaxPrefixLen),
		},
	}
}

func newNode48() *ArtNode {
	return &ArtNode{
		innerNode: &InnerNode{
			nodeType: Node48,
			keys:     make([]byte, Node256Max),
			children: make([]*ArtNode, Node48Max),
			prefix:   make([]byte, MaxPrefixLen),
		},
	}
}

func newNode256() *ArtNode {
	return &ArtNode{
		innerNode: &InnerNode{
			nodeType: Node256,
			children: make([]*ArtNode, Node256Max),
			prefix:   make([]byte, MaxPrefixLen),
		},
	}
}

func terminateKey(key []byte) []byte {
	index := bytes.Index(key, []byte{0})
	if index < 0 {
		return append(key, byte(0))
	}
	return key
}

func newLeafNode(key []byte, value interface{}) *ArtNode {
	newKey := make([]byte, len(key))
	copy(newKey, key)
	return &ArtNode{
		leaf: &LeafNode{newKey, value},
		innerNode: &InnerNode{
			nodeType: Leaf,
		},
	}

}

func (n *ArtNode) IsLeaf() bool {
	return n.innerNode.nodeType == Leaf
}

func (n *ArtNode) prefixMatch(key []byte, depth int) int {
	index := 0
	in := n.innerNode

	if in.prefixLen <= MaxPrefixLen {

		for ; index < in.prefixLen; index++ {
			if key[depth+index] != in.prefix[index] {
				return index
			}
		}
		return index
	}

	// in prefixLen > MaxPrefixLen condition
	for ; index < MaxPrefixLen; index++ {
		if key[depth+index] != in.prefix[index] {
			return index
		}
	}

	// 超出MaxPrefixLen
	min := n.Minimum()
	for ; index < n.innerNode.prefixLen; index++ {
		if min.leaf.key[depth+index] != key[depth+index] {
			return index
		}
	}

	return index

}

func (n *ArtNode) Minimum() *ArtNode {
	if n == nil {
		return nil
	}

	in := n.innerNode
	switch in.nodeType {
	case Leaf:
		return n
	case Node4, Node16:
		return in.children[0].Minimum()
	case Node48:
		i := 0
		for in.keys[i] == 0 {
			i++
		}
		child := in.children[in.keys[i]-1]
		return child.Minimum()
	case Node256:
		i := 0
		for in.children[i] == nil {
			i++
		}

		return in.children[i].Minimum()
	}
	return nil
}

func (in *InnerNode) findChild(key byte) **ArtNode {
	if in == nil {
		return nil
	}

	index := in.index(key)
	switch in.nodeType {
	case Node4, Node16, Node48:
		if index >= 0 {
			return &in.children[index]
		}

	case Node256:
		child := in.children[key]
		if child != nil {
			return &in.children[key]
		}
	}
	return nil

}

func (in *InnerNode) index(key byte) int {
	switch in.nodeType {
	case Node4, Node16:
		for i := 0; i < in.size; i++ {
			if in.keys[i] == key {
				return i
			}
		}
	case Node48:
		index := int(in.keys[key])
		if index > 0 {
			return index - 1
		}
	case Node256:
		return int(key)
	}
	return -1

}

func (l *LeafNode) IsMatch(key []byte) bool {
	return bytes.Equal(l.key, key)
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func copyBytes(dest, src []byte, nums int) {
	for i := 0; i < len(dest) && i < len(src) && i < nums; i++ {
		dest[i] = src[i]
	}
}

// prefixMatch 计算相同前缀长度
func (l *LeafNode) prefixMatch(key []byte, depth int) int {

	maxLen := min(len(l.key), len(key)) - depth

	i := 0
	for ; i < maxLen; i++ {
		if l.key[depth+i] != key[depth+i] {
			return i
		}
	}
	return i
}

func (n *InnerNode) IsFull() bool {
	return n.size == n.maxSize()
}

func (n *InnerNode) maxSize() int {
	switch n.nodeType {
	case Node4:
		return Node4Max
	case Node16:
		return Node16Max
	case Node48:
		return Node48Max
	case Node256:
		return Node256Max
	}

	return 0
}

func (in *InnerNode) minSize() int {
	switch in.nodeType {
	case Node4:
		return Node4Min
	case Node16:
		return Node16Min
	case Node48:
		return Node48Min
	case Node256:
		return Node256Min
	}
	return 0
}

func (n *InnerNode) grow() {

	switch n.nodeType {
	case Node4:
		n16 := newNode16().innerNode
		n16.prefix = n.prefix
		n16.size = n.size
		n16.prefixLen = n.prefixLen

		for i := 0; i < n.size; i++ {
			n16.keys[i] = n.keys[i]
			n16.children[i] = n.children[i]
		}
		*n = *n16
	case Node16:
		// n48 keys长度为256
		n48 := newNode48().innerNode
		n48.prefix = n.prefix
		n48.size = n.size
		n48.prefixLen = n.prefixLen

		index := 0
		for i := 0; i < n.size; i++ {
			child := n.children[i]
			if child != nil {
				n48.keys[n.keys[i]] = byte(index + 1)
				n48.children[index] = n.children[i]
				index++
			}
		}

		*n = *n48
	case Node48:
		n256 := newNode256().innerNode
		n256.prefix = n.prefix
		n256.size = n.size
		n256.prefixLen = n.prefixLen

		// key 为index
		for i := 0; i < len(n.keys); i++ {
			index := n.keys[i]
			if index > 0 {
				child := n.children[index-1]
				if child != nil {
					n256.children[byte(i)] = child
				}
			}
		}
		*n = *n256

	case Node256:
		// 无需扩展
	}

}

func (n *InnerNode) addChild(key byte, node *ArtNode) {

	if n.IsFull() {
		n.grow()
		n.addChild(key, node)
		return
	}

	switch n.nodeType {
	case Node4, Node16:
		{
			idx := 0
			for ; idx < n.size; idx++ {
				if key < n.keys[idx] {
					break
				}
			}

			for i := n.size; i > idx; i-- {

				if n.keys[i-1] > key {
					n.keys[i] = n.keys[i-1]
					n.children[i] = n.children[i-1]
				}
			}
			n.keys[idx] = key
			n.children[idx] = node
			n.size++
		}
	case Node48:
		// key 256  node 48
		idx := 0
		for i := 0; i < len(n.children); i++ {
			if n.children[idx] != nil {
				idx++
			}
		}

		n.children[idx] = node
		n.keys[key] = byte(idx + 1)
		n.size++
	case Node256:
		// key做为index node作为value
		n.children[key] = node
		n.size++

	}
}

func (n *ArtNode) deleteChild(key byte) {
	in := n.innerNode
	switch in.nodeType {
	case Node4, Node16:
		idx := in.index(key)
		if idx < 0 {
			return
		}
		in.keys[idx] = 0
		in.children[idx] = nil

		for i := idx; i < in.size-1; i++ {
			in.keys[i] = in.keys[i+1]
			in.children[i] = in.children[i+1]
		}

		in.keys[in.size-1] = 0
		in.children[in.size-1] = nil
		in.size -= 1

	case Node48:
		idx := in.index(key)
		if idx < 0 {
			return
		}
		child := in.children[idx]
		if child != nil {
			in.children[idx] = nil
			in.keys[key] = 0
			in.size -= 1
		}

	case Node256:
		idx := in.index(key)
		child := in.children[idx]
		if child != nil {
			in.children[idx] = nil
			in.size -= 1
		}

	}

	if in.size < in.minSize() {
		n.shrink()
	}
}

func (n *ArtNode) shrink() {
	in := n.innerNode
	switch in.nodeType {
	case Node4:
		c := in.children[0]
		if !c.IsLeaf() {
			child := c.innerNode
			currentPrefixLen := in.prefixLen

			// 前缀小于最大MaxPrefixLen时情况
			if currentPrefixLen < MaxPrefixLen {
				in.prefix[currentPrefixLen] = in.keys[0]
				currentPrefixLen++
			}

			if currentPrefixLen < MaxPrefixLen {
				childPrefixLen := min(child.prefixLen, MaxPrefixLen-currentPrefixLen)

				copyBytes(in.prefix[currentPrefixLen:], child.prefix, childPrefixLen)
				currentPrefixLen += childPrefixLen
			}

			copyBytes(child.prefix, in.prefix, min(currentPrefixLen, MaxPrefixLen))
			child.prefixLen += in.prefixLen + 1

		}
		*n = *c
	case Node16:
		n4 := newNode4()
		n4in := n4.innerNode
		n4in.prefix = n.innerNode.prefix
		n4in.prefixLen = n.innerNode.prefixLen

		for i := 0; i < len(in.keys); i++ {
			n4in.keys[i] = in.keys[i]
			n4in.children[i] = in.children[i]
			n4in.size++
		}
		*n = *n4

	case Node48:
		n16 := newNode16()
		n16in := n16.innerNode
		n16in.prefix = in.prefix
		n16in.prefixLen = in.prefixLen

		for i := 0; i < len(in.keys); i++ {
			idx := in.keys[byte(i)]
			if idx > 0 {
				child := in.children[idx-1]
				if child != nil {
					n16in.children[n16in.size] = child
					n16in.keys[n16in.size] = byte(i)
					n16in.size++
				}
			}

		}

		*n = *n16
	case Node256:
		n48 := newNode48()
		n48in := n48.innerNode
		n48in.prefix = in.prefix
		n48in.prefixLen = in.prefixLen

		for i := 0; i < len(in.children); i++ {
			child := in.children[byte(i)]
			if child != nil {
				n48in.keys[byte(i)] = byte(n48in.size + 1)
				n48in.children[n48in.size] = child
				n48in.size++
			}
		}

		*n = *n48

	}
}

func (t *ArtTree) Insert(key []byte, value interface{}) bool {
	key = terminateKey(key)

	updated := t.insert(&t.root, key, value, 0)
	if !updated {
		t.size++
	}
	return updated
}

func (t *ArtTree) insert(currentRef **ArtNode, key []byte, value interface{}, depth int) bool {
	current := *currentRef
	if current == nil {
		// 如果节点是空节点，则用叶子节点替换该节点
		*currentRef = newLeafNode(key, value)
		return false
	}

	if current.IsLeaf() {

		// 如果当前节点是叶子节点，如果key存在则退出
		// 叶子节点key匹配，直接修改value
		if current.leaf.IsMatch(key) {
			current.leaf.value = value
			return true
		}

		// 否则使用新的内部节点替换当前叶子节点，同时获取两个key前缀
		currentLeaf := current.leaf
		newLeaf := newLeafNode(key, value)
		limit := currentLeaf.prefixMatch(key, depth)

		n4 := newNode4()
		n4.innerNode.prefixLen = limit

		copyBytes(n4.innerNode.prefix, key[depth:], min(n4.innerNode.prefixLen, MaxPrefixLen))

		depth += n4.innerNode.prefixLen

		n4.innerNode.addChild(currentLeaf.key[depth], current)
		n4.innerNode.addChild(key[depth], newLeaf)

		*currentRef = n4
		return false
	}

	// 如果是内部节点
	in := current.innerNode

	if in.prefixLen != 0 {

		// 先比较前缀

		mMatch := current.prefixMatch(key, depth)

		// 前缀不符合，则生成新节点
		if mMatch != in.prefixLen {

			// 公共前缀节点
			n4 := newNode4()
			n4in := n4.innerNode
			*currentRef = n4
			n4in.prefixLen = mMatch // 新节点公共前缀

			copyBytes(n4in.prefix, in.prefix, mMatch)

			// 老节点插入新节点中
			if in.prefixLen < MaxPrefixLen {

				n4in.addChild(in.prefix[mMatch], current)
				in.prefixLen -= (mMatch + 1) // 减去新增节点前缀
				copyBytes(in.prefix, in.prefix[mMatch+1:], min(in.prefixLen, MaxPrefixLen))

			} else {

				in.prefixLen -= (mMatch + 1)
				minKey := current.Minimum().leaf.key
				n4in.addChild(minKey[depth+mMatch], current)
				copyBytes(in.prefix, minKey[depth+mMatch+1:], min(in.prefixLen, MaxPrefixLen))
			}

			// 叶子节点插入新节点中
			newLeafNode := newLeafNode(key, value)
			n4in.addChild(key[depth+mMatch], newLeafNode)
			return false

		}

		depth += in.prefixLen
	}

	// 如果前缀等，继续进行往子节点搜索
	next := in.findChild(key[depth])
	if next != nil {
		return t.insert(next, key, value, depth+1)
	}

	// 不存在下一层节点， 直接插入叶子节点
	newLeafNode := newLeafNode(key, value)
	in.addChild(key[depth], newLeafNode)

	return false

}

func (t *ArtTree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}

	key = terminateKey(key)
	result := t.delete(&t.root, key, 0)
	if result {
		t.size--
		return true
	}
	return false
}

func (t *ArtTree) delete(currentRef **ArtNode, key []byte, depth int) bool {
	current := *currentRef

	if current == nil {
		return false
	}

	if current.IsLeaf() {
		if current.leaf.IsMatch(key) {
			*currentRef = nil
			return true
		}
		return false
	}

	// 当前是非叶子节点
	in := current.innerNode
	if in.prefixLen > 0 {
		mMatch := current.prefixMatch(key, depth)
		if mMatch != in.prefixLen {
			// 前缀不一致，直接返回
			return false
		}
		depth += in.prefixLen
	}
	// 接着找下个节点
	next := in.findChild(key[depth])
	if *next != nil && (*next).IsLeaf() && (*next).leaf.IsMatch(key) {
		current.deleteChild(key[depth])
		return true
	}
	return t.delete(next, key, depth+1)

}

func (t *ArtTree) Search(key []byte) interface{} {
	key = terminateKey(key)
	return t.search(t.root, key, 0)
}

func (t *ArtTree) search(current *ArtNode, key []byte, depth int) interface{} {
	for current != nil {

		if current.IsLeaf() {
			if current.leaf.IsMatch(key) {
				return current.leaf.value
			}
			return nil
		}

		in := current.innerNode
		if current.prefixMatch(key, depth) != in.prefixLen {
			return nil
		}

		depth += in.prefixLen
		next := in.findChild(key[depth])

		if next == nil {
			return nil
		}

		current = *next

		depth++

	}

	return nil

}
