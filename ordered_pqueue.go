package utils

import (
	"sync"
)

type OrderedParallelQueue struct {
	concurrency     int
	concurrencyChan chan struct{}

	inputChanCap, outputChanCap int
	inputChan, outputChan       chan interface{}

	DoFunc func(x interface{}) (y interface{})

	quit chan struct{}
	wg   sync.WaitGroup
}

func NewOrderedParallelQueue(concurrency int, do_func func(x interface{}) (y interface{})) (q *OrderedParallelQueue) {
	q = &OrderedParallelQueue{
		concurrency:     concurrency,
		inputChanCap:    10,
		outputChanCap:   10,
		DoFunc:          do_func,
		concurrencyChan: make(chan struct{}, concurrency),
		quit:            make(chan struct{}),
	}

	q.inputChan = make(chan interface{}, q.inputChanCap)
	q.outputChan = make(chan interface{}, q.outputChanCap)

	return
}

func NewOrderedParallelQueueWithCap(concurrency, inputChan_cap, outputChan_cap int, do_func func(x interface{}) (y interface{})) (q *OrderedParallelQueue) {
	q = &OrderedParallelQueue{
		concurrency:     concurrency,
		inputChanCap:    inputChan_cap,
		outputChanCap:   outputChan_cap,
		DoFunc:          do_func,
		concurrencyChan: make(chan struct{}, concurrency),
		quit:            make(chan struct{}),
	}

	q.inputChan = make(chan interface{}, q.inputChanCap)
	q.outputChan = make(chan interface{}, q.outputChanCap)

	return
}

// 开始并发处理输入数据
func (q *OrderedParallelQueue) Start() {
	go func() {
		for i := 0; i < q.concurrency; i++ {
			q.concurrencyChan <- struct{}{}
		}

		prev := make(chan (struct{}))
		close(prev)

		for x := range q.inputChan {
			<-q.concurrencyChan

			this := make(chan (struct{}))

			go func(x interface{}, prev, this chan struct{}) {
				y := q.DoFunc(x)

				<-prev
				q.outputChan <- y
				q.concurrencyChan <- struct{}{}
				close(this)
				q.wg.Done()
			}(x, prev, this)

			prev = this
		}

	}()
}

// 关闭输入channel，等待所有处理完毕后，关闭输出队列。
func (q *OrderedParallelQueue) Stop() {
	close(q.quit)
	close(q.inputChan)

	// TODO 等待output消耗完？
	q.wg.Wait()
	close(q.outputChan)

}

// 异步关闭
func (q *OrderedParallelQueue) AsyncStop() {
	go q.Stop()
}

// 将处理数据压入队列异步并发执行。如果队列已被关闭则返回false; 如果成功则返回true。如果队列已满，则阻塞直至成功插入队列。。
func (q *OrderedParallelQueue) Enqueue(x interface{}) (ok bool) {

	select {
	case q.inputChan <- x:
		q.wg.Add(1)
		return true
	case <-q.quit:
		return false
	}
}

// 获取一个处理完毕的数据。如chan已关闭，则ok为false，否则阻塞直至下一个处理完毕的数据返回。
func (q *OrderedParallelQueue) Dequeue() (y interface{}, ok bool) {
	y, ok = <-q.outputChan

	return
}

// 获取input chan以进行高级操作。 注意手动传值的时候需要wg.Add()
func (q *OrderedParallelQueue) InputChan() chan interface{} {
	return q.inputChan
}

// 获取output chan以进行高级操作。 不涉及wg
func (q *OrderedParallelQueue) OutputChan() chan interface{} {
	return q.outputChan
}

// 获取输入暂存队列的长度
func (q *OrderedParallelQueue) InputChanCap() int {
	return q.inputChanCap
}

// 获取输出暂存队列的长度
func (q *OrderedParallelQueue) OutputChanCap() int {
	return q.outputChanCap
}

// 获取能够同时并行执行逻辑的数量
func (q *OrderedParallelQueue) Concurrency() int {
	return q.concurrency
}

func (q *OrderedParallelQueue) Wg() sync.WaitGroup {
	return q.wg
}
