package timing

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/pborman/uuid"
)

// HandlerFunc 处理函数类型
type HandlerFunc func(items ...*Item)

var (
	// PersistFunc 定时项存储，便于系统重启后恢复
	PersistFunc HandlerFunc = func(items ...*Item) {}
	// DeleteFunc 时间到了，存储中移除定时项
	DeleteFunc HandlerFunc = func(items ...*Item) {}
	// RemindFunc 时间到了，提醒处理定时项
	RemindFunc HandlerFunc = func(items ...*Item) {
		fmt.Printf("default remind: %#v\n", items)
	}
)

var (
	// 用于通知新增一个定时项
	notify = make(chan *Item)
	// 保证定时器初始化先于其它操作
	inited = make(chan struct{})
	// 仅执行一次
	once sync.Once
)

// Init 初始化定时器
func Init(items ...*Item) {
	once.Do(func() {
		// 定时项插入队列
		q := make(Queue, len(items), 1024)
		for i, item := range items {
			q[i] = item
		}
		go start(q)
	})
}

// Add 添加定时项
func Add(items ...*Item) {
	// 等待定时项初始化工作结束
	<-inited

	// 当前时间
	now := uint32(time.Now().Unix())
	for _, item := range items {
		// 提醒时间到了，需要处理定时器
		if item.Timestamp < now {
			RemindFunc(item)
		}

		item.ID = uuid.New()
		// 通知新增一个定时项
		notify <- item
	}
}

func start(q Queue) {
	// 初始化堆，和队列关联，队列元素按照大->小的顺序排列
	heap.Init(&q)

	var min *Item
	// 创建一个Timer
	var timer = time.NewTimer(24 * time.Hour)

	if len(q) > 0 {
		// 获取数值最小的元素
		min = heap.Pop(&q).(*Item)
		// 重新设置timer的超时时间
		timer.Reset(time.Unix(int64(min.Timestamp), 0).Sub(time.Now()))
	}

	// 允许定时项的插入等操作
	close(inited)

	for {
		select {
		case item := <-notify: // 新增一个定时项
			// 持久化新增的定时项
			PersistFunc(item)

			if min == nil {
				// 如果当前处理的定时项为空，直接更新当前处理定时项
			} else if item.Timestamp < min.Timestamp {
				// 如果小于当前定时项数值，替换当前正在处理的定时项
				heap.Push(&q, min)
			} else {
				// 如果大于当前定时项数值，直接插入队列
				heap.Push(&q, item)
				break
			}

			// 更新当前处理的定时项
			min = item
			// 重新设置定时器
			timer.Reset(time.Unix(int64(min.Timestamp), 0).Sub(time.Now()))

		case <-timer.C: // 定时器超时
			if min != nil {
				// 缓存超时的定时项
				timeout := []*Item{min}
				for {
					// 判断队列是否为空
					if q.Len() == 0 {
						min = nil
						break
					}
					// 获取下一个数值最小的定时项
					next := heap.Pop(&q).(*Item)
					// 如果数值非当前定时项的值，退出循环
					if next.Timestamp != min.Timestamp {
						min = next
						break
					}
					// 缓存所有和当前定时项数值相同的项
					timeout = append(timeout, next)
				}

				// 提醒或者删除定时项
				go func(timeout []*Item) {
					RemindFunc(timeout...)
					DeleteFunc(timeout...)
				}(timeout)
			}

			if min == nil {
				// 如果定时项为空，默认定时器为24小时
				timer.Reset(24 * time.Hour)
				break
			}

			// 设置定时器的下一个触发定时项
			timer.Reset(time.Unix(int64(min.Timestamp), 0).Sub(time.Now()))
		}
	}
}
