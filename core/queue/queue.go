package queue

import (
	"context"
	"fmt"

	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uc1024/f90/core/rescue"
	"github.com/uc1024/f90/core/slogx"
	"github.com/uc1024/f90/core/threadingx"
)

const queueName = "default queue"

type (
	Queue struct {
		name                 string
		producerFactory      ProducerFactory
		producerRoutineGroup *threadingx.RoutineGroup
		consumerFactory      ConsumerFactory
		consumerRoutineGroup *threadingx.RoutineGroup
		producerQuantity     int
		consumerQuentity     int
		active               int32
		channel              chan string
		quit                 chan struct{}
		listeners            []Listener // * 队列事件监听器
		eventLock            sync.Mutex // * 事件管道锁,在操作事件管道时需要加锁
		eventChannels        []chan any
	}
)

func NewQueue(
	producerFactory ProducerFactory,
	consumerFactory ConsumerFactory) *Queue {
	queue := Queue{
		name:                 queueName,
		producerFactory:      producerFactory,
		producerRoutineGroup: threadingx.NewRoutineGroup(),
		consumerFactory:      consumerFactory,
		consumerRoutineGroup: threadingx.NewRoutineGroup(),
		producerQuantity:     runtime.NumCPU(),
		consumerQuentity:     runtime.NumCPU() << 1,
		channel:              make(chan string),
		quit:                 make(chan struct{}),
	}
	return &queue
}

/*
启动队列
开始生产者和消费者的线程，并等待所有线程结束后关闭通道
*/
func (q *Queue) Start() {
	q.startProducers(q.producerQuantity)
	q.startConsumers(q.consumerQuentity)

	q.producerRoutineGroup.Wait()
	close(q.channel)
	q.consumerRoutineGroup.Wait()
}

/*
停止队列，关闭quit通道，使所有线程都退出
*/
func (q *Queue) Stop() {
	close(q.quit)
}

/*
添加监听器
将监听器添加到队列的listeners切片中
*/
func (q *Queue) AddListener(listener Listener) {
	q.listeners = append(q.listeners, listener)
}

/*
广播事件，将事件发送到所有监听器的事件通道中
*/
func (q *Queue) Broadcast(message any) {
	go func() {
		q.eventLock.Lock()
		defer q.eventLock.Unlock()

		for _, channel := range q.eventChannels {
			channel <- message
		}
	}()
}

/*
启动生产者线程，根据指定数量创建并启动生产者线程
*/
func (q *Queue) startProducers(count int) {
	for i := 0; i < count; i++ {
		q.producerRoutineGroup.Run(func() {
			q.produce()
		})
	}
}

/*
生产消息，循环从生产者工厂获取生产者并执行生产操作。
*/
func (q *Queue) produce() {
	var producer Producer

	for {
		var err error
		if producer, err = q.producerFactory(); err != nil {
			slogx.Default.Error(context.Background(),
				fmt.Sprintf("producer factory error: %s", err.Error()))
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	// * 活动实例++
	atomic.AddInt32(&q.active, 1)
	// * 添加监听器
	producer.AddListener(routineListener{
		queue: q,
	})

	for {
		select {
		case <-q.quit: // * 监听退出信号
			slogx.Default.Info(context.Background(), "quitting producer")
			return
		default: // * 生产消息
			if v, ok := q.produceOne(producer); ok {
				q.channel <- v //
			}
		}
	}
}

/*
生产一条消息，执行生产者的Produce方法，并返回生成的消息和是否成功
*/
func (q *Queue) produceOne(producer Producer) (string, bool) {
	// * 捕获panic
	defer rescue.Catch()

	return producer.Produce()
}

/*
启动消费者线程，根据指定数量创建并启动消费者线程。
*/
func (q *Queue) startConsumers(count int) {
	for i := 0; i < count; i++ {
		eventChan := make(chan any)
		q.eventLock.Lock()
		q.eventChannels = append(q.eventChannels, eventChan)
		q.eventLock.Unlock()
		q.consumerRoutineGroup.Run(func() {
			q.consume(eventChan)
		})
	}
}

/*
消费消息，循环从通道接收消息并执行消费操作，同时监听事件通道并执行相应的操作。
*/
func (q *Queue) consume(eventChan chan any) {
	var consumer Consumer

	// * 构建消费者直至成功
	for {
		var err error
		if consumer, err = q.consumerFactory(); err != nil {
			slogx.Default.Error(context.Background(), fmt.Sprintf("Error on creating consumer: %v", err))
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	// * 阻塞并竞争消费数据
	for {
		select {
		case message, ok := <-q.channel:
			if ok {
				q.consumeOne(consumer, message) // * 管道取出消息并消费
			} else {
				slogx.Default.Info(context.Background(), "Task channel was closed, quitting consumer...")
				return
			}
		case event := <-eventChan:
			consumer.OnEvent(event)
		}
	}
}

/*
消费一条消息，执行消费者的Consume方法，并捕获panic。
*/
func (q *Queue) consumeOne(consumer Consumer, message string) {
	threadingx.RunSafe(func() {
		if err := consumer.Consume(message); err != nil {
			slogx.Default.Error(context.Background(), fmt.Sprintf("Error on consuming message: %s error: %s", message, err.Error()))
		}
	})
}

func (q *Queue) resume() {
	for _, listener := range q.listeners {
		listener.OnResume()
	}
}

func (q *Queue) pause() {
	for _, listener := range q.listeners {
		listener.OnPause()
	}
}

func (q *Queue) SetName(s string) {
	q.name = s
}

func (q *Queue) SetConsumerQuantity(i int) {
	q.consumerQuentity = i
}

func (q *Queue) SetProducerQuantity(i int) {
	q.producerQuantity = i
}

func (q *Queue) GetActiveCount(f ProducerFactory) int32 {
	return atomic.LoadInt32(&q.active)
}

/*
常规监听器
*/
type routineListener struct {
	queue *Queue
}

func (rl routineListener) OnProducerPause() {
	if atomic.AddInt32(&rl.queue.active, -1) <= 0 {
		rl.queue.pause()
	}
}

func (rl routineListener) OnProducerResume() {
	if atomic.AddInt32(&rl.queue.active, 1) == 1 {
		rl.queue.resume()
	}
}
