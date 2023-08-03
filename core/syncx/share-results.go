package syncx

import "sync"

type (
	ShareResults interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
		DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
	}

	call struct {
		wg        sync.WaitGroup
		resources interface{}
		err       error
	}

	shareCall struct {
		lock  sync.Mutex
		calls map[string]*call
	}
)

func NewShareCall() ShareResults {
	return &shareCall{
		calls: map[string]*call{},
	}
}

func (s *shareCall) createCall(key string) (*call, bool) {
	s.lock.Lock()
	if c, ok := s.calls[key]; ok {
		s.lock.Unlock()
		// * 说明有相同key的在获取资源
		// * 所以在这等待比在这跟早的回去的回调,共享它获取的数据
		c.wg.Wait()
		// * 对应key的回调完成了
		return c, true 
	}

	// * 没有相同的key调用(自己来)
	call := new(call)
	call.wg.Add(1)
	s.calls[key] = call
	s.lock.Unlock()

	return call, false
}

func (s *shareCall) makeCall(c *call, key string, fn func() (interface{}, error)) {

	// & 释放资源并通知其他等待进程解除阻塞获取资源
	defer func() {
		s.lock.Lock()
		delete(s.calls, key)
		s.lock.Unlock()
		c.wg.Done()
	}()

	// * 调用回调获取资源
	c.resources, c.err = fn()
}

func (s *shareCall) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	call, done := s.createCall(key)

	if done {
		return call.resources, call.err
	}

	s.makeCall(call, key, fn)

	return call.resources, call.err

}

func (s *shareCall) DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error) {
	call, done := s.createCall(key)

	if done {
		return call.resources, !done, call.err
	}

	s.makeCall(call, key, fn)

	return call.resources, !done, call.err

}
