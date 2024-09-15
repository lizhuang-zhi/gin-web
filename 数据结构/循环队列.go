package main

import "fmt"

/*
循环队列使用情况:
1. 如果你预先知道队列的最大大小，并且希望避免频繁地分配和释放内存
2. 数据缓冲区
3. 生产者消费者问题(生产者和消费者的速度不一致)
*/

type RingQueue struct {
	data  []int // 数据
	head  int   // 头指针
	tail  int   // 尾指针
	count int   // 元素个数
	size  int   // 队列大小
}

func NewRingQueue(size int) *RingQueue {
	return &RingQueue{
		data: make([]int, size),
		size: size,
	}
}

// 入队(从尾部插入)
func (r *RingQueue) Enqueue(item int) bool {
	if r.count == r.size { // 队列已满
		return false
	}

	r.data[r.tail] = item
	r.tail = (r.tail + 1) % r.size // 循环队列尾指针+1
	r.count++

	return true
}

// 出队(从头部删除)
func (r *RingQueue) Dequeue() (int, bool) {
	if r.count == 0 { // 队列为空
		return 0, false
	}

	item := r.data[r.head]
	r.head = (r.head + 1) % r.size // 循环队列头指针+1
	r.count--

	return item, true
}

func main() {
	queue := NewRingQueue(5)
	b := queue.Enqueue(1)
	b = queue.Enqueue(2)
	b = queue.Enqueue(3)
	b = queue.Enqueue(4)
	item, b := queue.Dequeue()
	if b {
		fmt.Println(item)
	}
	item, b = queue.Dequeue()
	if b {
		fmt.Println(item)
	}
	item, b = queue.Dequeue()
	if b {
		fmt.Println(item)
	}
	item, b = queue.Dequeue()
	if b {
		fmt.Println(item)
	}
	b = queue.Enqueue(55)
	b = queue.Enqueue(66)
	b = queue.Enqueue(77)
	item, b = queue.Dequeue()
	if b {
		fmt.Println(item)
	}
	item, b = queue.Dequeue()
	if b {
		fmt.Println(item)
	}
}
