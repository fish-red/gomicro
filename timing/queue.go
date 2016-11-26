package timing

// Item 定义队列的元素
type Item struct {
	ID          string
	Timestamp   uint32
	Event       string
	Description string
}

// Queue 包含多个元素
type Queue []*Item

// Queue支持Sort功能需要实现的接口

// Len 返回元素个数
func (q Queue) Len() int {
	return len(q)
}

// Less 比较两个元素
func (q Queue) Less(i, j int) bool {
	return q[i].Timestamp <= q[j].Timestamp
}

// Swap 交换两个元素的值
func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Push 队列插入元素
func (q *Queue) Push(x interface{}) {
	item := x.(*Item)
	*q = append(*q, item)
}

// Pop 队列输出元素，返回数值最小的元素（小根堆）
func (q *Queue) Pop() interface{} {
	old := *q
	// 队列的长度
	n := q.Len()
	// 返回数值最小的元素
	item := old[n-1]
	// 更新队列
	*q = old[0 : n-1]

	return item
}
