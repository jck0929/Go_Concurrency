package main

// CAS操作，当时还没有atomic包
func cas(val *int32, old, new int32) bool {}
func semacquire(*int32)                   {}
func semrelease(*int32)                   {}

// 互斥锁的结构
type Mutex struct {
	key  int32 // 锁是否被持有的标致
	sema int32 // 信号量专用，用以唤醒/阻塞goroutine
}

// 保证成功在val上增加delta的值
func xadd(val *int32, delta int32) (new int32) {
	for {
		v := *val
		if cas(val, v, v+delta) {
			return v + delta
		}

		panic("unreached")
	}
}

// 请求锁
func (m *Mutex) Lock() {
	// 标识+1，如果等于1，则成功获取到锁
	if xadd(&m.key, 1) == 1 {
		return
	}
	// 否则阻塞等待
	semacquire(&m.sema)
}

// 释放锁
func (m *Mutex) UnLock() {
	// 标识-1，如果等于0，则没有其他等待者
	if xadd(&m.key, -1) == 0 {
		return
	}
	// 否则唤醒其他等待者
	semrelease(&m.sema)
}
