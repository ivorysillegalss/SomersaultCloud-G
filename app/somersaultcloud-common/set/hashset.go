package set

type HashSet struct {
	set map[interface{}]struct{}
}

// NewHashSet 创建一个新的 HashSet 实例
func NewHashSet() *HashSet {
	return &HashSet{
		set: make(map[interface{}]struct{}),
	}
}

// Add 向 HashSet 中添加一个元素
func (h *HashSet) Add(value interface{}) {
	h.set[value] = struct{}{}
}

// Remove 从 HashSet 中删除一个元素
func (h *HashSet) Remove(value interface{}) {
	delete(h.set, value)
}

// Contains 检查 HashSet 是否包含一个元素
func (h *HashSet) Contains(value interface{}) bool {
	_, exists := h.set[value]
	return exists
}

// Size 返回 HashSet 中元素的数量
func (h *HashSet) Size() int {
	return len(h.set)
}

// Clear 清空 HashSet 中的所有元素
func (h *HashSet) Clear() {
	h.set = make(map[interface{}]struct{})
}

// Iter 返回一个迭代器，允许遍历 HashSet 中的所有元素
func (h *HashSet) Iter() []interface{} {
	var keys []interface{}
	for key := range h.set {
		keys = append(keys, key)
	}
	return keys
}
