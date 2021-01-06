package work

import (
	"sync"
	"time"
)

// 滑动窗口计数有很多使用场景，比如说限流防止系统雪崩。相比计数实现，滑动窗口实现会更加平滑，能自动消除毛刺。
// 滑动窗口原理是在每次有访问进来时，先判断前N个单位时间内的总访问量是否超过了设置的阈值，并对当前时间片上的请求数+1。
// 3s 总访问量>10 限流
const (
	typeSuccess int = 1
	typeFail    int = 2
)

//Metrics 指标
type Metrics struct {
	success int64
}

//SlidingWindow 滑动窗口
type SlidingWindow struct {
	bucket int                //桶数
	curKey int64              //当前key
	m      map[int64]*Metrics //统计
	sync.RWMutex
}

//NewSlidingWindow 创建滑动窗口
func NewSlidingWindow(bucket int) *SlidingWindow {
	sw := &SlidingWindow{}
	sw.bucket = bucket
	return sw
}

//自增操作
func (sw *SlidingWindow) increment() bool {
	sw.Lock()
	defer sw.Unlock()

	nowTime := time.Now().Unix()
	if _, ok := sw.m[nowTime]; !ok {
		sw.m = make(map[int64]*Metrics)
		sw.m[nowTime] = &Metrics{}
		sw.curKey = nowTime
	}
	var sum int64 = 0
	for i := nowTime - 1; i < nowTime-3; i-- {
		if _, ok := sw.m[i]; !ok {
			continue
		}
		sum += sw.m[i].success
	}
	if 10 < sum {
		return false
	}
	sw.m[nowTime].success++
	return true
}

//Len 获取数据长度
func (sw *SlidingWindow) Len() int {
	//return sw.data.Len()
	return 0
}
