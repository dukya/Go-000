package rolling

//参考https://github.com/go-kratos/kratos/tree/v2
import (
	"fmt"
	"sync"
	"time"
)

// 滑动窗口计数器配置
type RollingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

// 滑动窗口计数器
type RollingCounter struct {
	mu     sync.RWMutex
	size   int
	window *Window
	offset int

	bucketDuration time.Duration
	lastAddTime    time.Time
}

// 新建一个滑动窗口计数器
func NewRollingCounter(opt RollingCounterOpts) *RollingCounter {
	w := NewWindow(WindowOpt{Size: opt.Size})
	return &RollingCounter{
		window:         w,
		size:           w.Size(),
		offset:         0,
		bucketDuration: opt.BucketDuration,
		lastAddTime:    time.Now(),
	}
}

// 获取当前时间与最近一次累加操作所在桶的起始时间之间的间隔桶数
func (r *RollingCounter) TimeSpan() int {
	num := int(time.Since(r.lastAddTime) / r.bucketDuration)
	if num > -1 {
		return num //正常情况下应该是0,1,2,3,...
	}
	return r.size
}

// 累加数据操作(写锁)
func (r *RollingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("cannot decrease in value. val: %d", val))
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	timespan := r.TimeSpan()
	// 若timespan==0则直接写入到当前桶中
	if timespan > 0 {
		r.lastAddTime = r.lastAddTime.Add(time.Duration(timespan * int(r.bucketDuration)))
		offset := r.offset

		if timespan > r.size { //超过滑动窗口的总时长，直接重置所有桶
			//timespan = r.size
			r.window.ResetWindow()
			t := offset + (timespan % r.size)
			r.offset = t % r.size
		} else {
			// reset the expired buckets
			s := offset + 1 //必须从offset+1的桶开始重置

			e, e1 := s+timespan, 0
			if e > r.size {
				e1 = e - r.size
				e = r.size
			}
			//[offset+1,e)的桶进行过期重置	offset+1 <= e <= r.size
			for i := s; i < e; i++ {
				r.window.ResetBucket(i)
				offset = i
			}
			//[0,e1)的桶进行过期重置		0 <= e1 <= r.size
			for i := 0; i < e1; i++ {
				r.window.ResetBucket(i)
				offset = i
			}
			r.offset = offset
		}
	}
	r.window.Add(r.offset, float64(val))
}

// 进行计数器数据聚合操作的通用方法
func (r *RollingCounter) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	timespan := r.TimeSpan()
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	r.mu.RUnlock()
	return val
}

// 计算平均值(读锁)
func (r *RollingCounter) Avg() float64 {
	return r.Reduce(Avg)
}

// 计算总和(读锁)
func (r *RollingCounter) Sum() float64 {
	return r.Reduce(Sum)
}

// 计算最大值(读锁)
func (r *RollingCounter) Max() float64 {
	return r.Reduce(Max)
}

// 计算最小值(读锁)
func (r *RollingCounter) Min() float64 {
	return r.Reduce(Min)
}

func (r *RollingCounter) Value() int64 {
	return int64(r.Sum())
}
