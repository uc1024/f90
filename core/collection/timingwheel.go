package collection

import (
	"container/list"
	"errors"
	"fmt"

	"log"
	"time"

	"github.com/uc1024/f90/core/rescue"
	"github.com/uc1024/f90/core/threadingx"
	"github.com/uc1024/f90/core/timex"
)

const drainWorkers = 8

var (
	ErrClosed   = errors.New("TimingWheel is closed already")
	ErrArgument = errors.New("incorrect task argument")
)

type (
	// exec
	Execute func(key, value interface{})

	TimingWheel struct {
		interval      time.Duration // interval time (single time of scale)
		ticker        timex.Ticker  // promote time forward
		slots         []*list.List  // data performance of the timingWheel
		timers        *SafeMap      // store key value of the execute
		tickedPos     int           // numslots-1 ( pos the current slot  )
		numSlots      int           // one circle how much slot
		execute       Execute
		setChannel    chan timingEntry
		moveChannel   chan baseEntry
		removeChannel chan interface{}
		drainChannel  chan func(key, value interface{})
		stopChannel   chan struct{}
	}

	timingEntry struct {
		baseEntry
		value   interface{} // * 储存的value
		circle  int         // * 第几轮
		diff    int
		removed bool
	}

	baseEntry struct {
		delay time.Duration
		key   interface{}
	}

	positionEntry struct {
		pos  int
		item *timingEntry
	}

	timingTask struct {
		key   interface{}
		value interface{}
	}
)

func NewTimingWheel(interval time.Duration, numSlots int, execute Execute) (*TimingWheel, error) {
	if interval <= 0 || numSlots <= 0 || execute == nil {
		return nil, fmt.Errorf("interval: %v, slots: %d, execute: %p", interval, numSlots, execute)
	}

	return newTimingWheelWithClock(interval, numSlots, execute, timex.NewTicker(interval))
}

// * 执行所有任务
func (tw *TimingWheel) Drain(fn func(key, value interface{})) error {
	select {
	case tw.drainChannel <- fn:
		return nil
	case <-tw.stopChannel:
		return ErrClosed
	}
}

// * 修改任务执行时间
func (tw *TimingWheel) MoveTimer(key interface{}, delay time.Duration) error {
	if delay <= 0 || key == nil {
		return ErrArgument
	}
	select {
	case tw.moveChannel <- baseEntry{
		delay: delay,
		key:   key,
	}:
		return nil
	case <-tw.stopChannel:
		return ErrClosed
	}

}

// * 移除指定key的任务
func (tw *TimingWheel) RemoveTimer(key interface{}) error {
	if key == nil {
		return ErrArgument
	}
	select {
	case tw.removeChannel <- key:
		return nil
	case <-tw.stopChannel:
		return ErrClosed
	}
}

// * 设置一个定时任务
func (tw *TimingWheel) SetTimer(key, value interface{}, delay time.Duration) error {
	if delay <= 0 || key == nil {
		return ErrArgument
	}

	select {
	case tw.setChannel <- timingEntry{
		baseEntry: baseEntry{
			delay: delay,
			key:   key,
		},
		value: value,
	}:
		return nil
	case <-tw.stopChannel:
		return ErrClosed
	}

}

// 关闭
func (tw *TimingWheel) Stop() {
	close(tw.stopChannel)
}

// create timingWheel with clock
func newTimingWheelWithClock(interval time.Duration, numSlots int, execute Execute, ticker timex.Ticker) (
	*TimingWheel, error) {

	tw := &TimingWheel{
		interval:      interval,
		ticker:        ticker,
		slots:         make([]*list.List, numSlots),
		timers:        NewSafeMap(),
		tickedPos:     numSlots - 1, // at previous virtual circle
		execute:       execute,
		numSlots:      numSlots,
		setChannel:    make(chan timingEntry),
		moveChannel:   make(chan baseEntry),
		removeChannel: make(chan interface{}),
		drainChannel:  make(chan func(key, value interface{})),
		stopChannel:   make(chan struct{}),
	}

	tw.initSlots()
	go tw.run()

	return tw, nil
}

func (tw *TimingWheel) initSlots() {
	for i := 0; i < tw.numSlots; i++ {
		tw.slots[i] = list.New()
	}
}

func (tw *TimingWheel) run() {
	for {
		select {
		case <-tw.ticker.Chan(): // 定时器推动时间轮向前
			tw.onTick()
		case task := <-tw.setChannel: // 添加任务
			tw.setTask(&task)
		case key := <-tw.removeChannel: // 移除任务
			tw.removeTask(key)
		case task := <-tw.moveChannel: // 移动任务
			tw.moveTask(task)
		case fn := <-tw.drainChannel: //
			tw.drainAll(fn)
		case <-tw.stopChannel: // 停止
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimingWheel) moveTask(task baseEntry) {
	val, ok := tw.timers.Get(task.key) // 用任务key从map取出任务
	if !ok {
		log.Println("没有数据")
		return
	}

	timer := val.(*positionEntry) // val中存储的结构
	if task.delay < tw.interval { // 如果任务延迟时间小于间隔时间直接执行
		go func() {
			defer rescue.Catch()
			tw.execute(timer.item.key, timer.item.value)
		}()

		return
	}

	pos, circle := tw.getPositionAndCircle(task.delay) // 获取延迟任务在时间轮的上(轮数)(轮上的位置)
	if pos >= timer.pos {                              // 如果获取的位置大于等于原来位置
		timer.item.circle = circle        // 设置任务新的轮
		timer.item.diff = pos - timer.pos // 偏移量 (任务移动了多少)
	} else if circle > 0 { // 如果轮数大于0
		circle-- // 这里-- 是减去当前轮 转移到下一轮
		timer.item.circle = circle
		timer.item.diff = tw.numSlots + pos - timer.pos // 向后偏移了多少
	} else {
		// 任务执行的位置已经过去了,重新入队等待执行
		timer.item.removed = true // 原任务打上移除标识符
		newItem := &timingEntry{  // 新任务入队
			baseEntry: task,
			value:     timer.item.value,
		}
		tw.slots[pos].PushBack(newItem)
		tw.setTimerPosition(pos, newItem)
	}
}

func (tw *TimingWheel) onTick() {
	tw.tickedPos = (tw.tickedPos + 1) % tw.numSlots // 执行位置
	l := tw.slots[tw.tickedPos]                     // 取出当前的时间槽
	tw.scanAndRunTasks(l)                           // 扫描并执行任务
}

func (tw *TimingWheel) removeTask(key interface{}) {
	val, ok := tw.timers.Get(key)
	if !ok {
		return
	}

	timer := val.(*positionEntry)
	timer.item.removed = true
	tw.timers.Del(key)
}

func (tw *TimingWheel) runTasks(tasks []timingTask) {
	if len(tasks) == 0 {
		return
	}

	go func() {
		for i := range tasks {
			threadingx.RunSafe(func() {
				tw.execute(tasks[i].key, tasks[i].value)
			})
		}
	}()
}

func (tw *TimingWheel) scanAndRunTasks(l *list.List) {
	var tasks []timingTask

	for e := l.Front(); e != nil; {
		task := e.Value.(*timingEntry)
		if task.removed { // 移除的任务
			next := e.Next()
			l.Remove(e)
			e = next
			continue
		} else if task.circle > 0 { // 不是当前轮要执行的任务 -- 轮数
			task.circle--
			e = e.Next()
			continue
		} else if task.diff > 0 { // 如果任务偏移量大于0(说明任务需要重新入队)
			next := e.Next()
			l.Remove(e)
			// (tw.tickedPos+task.diff)%tw.numSlots
			// cannot be the same value of tw.tickedPos
			pos := (tw.tickedPos + task.diff) % tw.numSlots
			tw.slots[pos].PushBack(task)
			tw.setTimerPosition(pos, task)
			task.diff = 0
			e = next
			continue
		}
		// 需要执行的任务放入数组
		tasks = append(tasks, timingTask{
			key:   task.key,
			value: task.value,
		})
		next := e.Next() // 从slots.list 中移除
		l.Remove(e)
		tw.timers.Del(task.key) // 从任务map中移除
		e = next
	}

	// 执行
	tw.runTasks(tasks)
}

func (tw *TimingWheel) setTask(task *timingEntry) {
	if task.delay < tw.interval { // 如果任务的延迟时间小于 设定的运行间隔(执行间隔)
		task.delay = tw.interval
	}

	if val, ok := tw.timers.Get(task.key); ok { // 如果有相同的任务存在
		entry := val.(*positionEntry)
		entry.item.value = task.value // value覆盖
		tw.moveTask(task.baseEntry)   // 重新设置执行时间
	} else {
		pos, circle := tw.getPositionAndCircle(task.delay)
		// fmt.Println(pos, circle)
		task.circle = circle
		tw.slots[pos].PushBack(task)
		tw.setTimerPosition(pos, task)
	}
}

func (tw *TimingWheel) setTimerPosition(pos int, task *timingEntry) {
	if val, ok := tw.timers.Get(task.key); ok {
		timer := val.(*positionEntry)
		timer.item = task
		timer.pos = pos
	} else {
		tw.timers.Set(task.key, &positionEntry{
			pos:  pos,
			item: task,
		})
	}
}

// * 执行所有任务并清空
func (tw *TimingWheel) drainAll(fn func(key, value interface{})) {
	runner := threadingx.NewTaskRunner(drainWorkers)
	for _, slot := range tw.slots {
		for e := slot.Front(); e != nil; {
			task := e.Value.(*timingEntry)
			next := e.Next()
			slot.Remove(e)
			e = next
			if !task.removed {
				runner.Schedule(func() {
					fn(task.key, task.value)
				})
			}
		}
	}
}

// * circle 获取到在第几轮
// * pos 	获取在轮上的位子
func (tw *TimingWheel) getPositionAndCircle(d time.Duration) (pos, circle int) {
	steps := int(d / tw.interval)
	pos = (tw.tickedPos + steps) % tw.numSlots
	circle = (steps - 1) / tw.numSlots
	return
}
