package rolling

// 桶定义
type Bucket struct {
	Point float64 //累加数据值
	Count int64   //累加次数
	next  *Bucket //指向下一个桶
}

// 累加数据
func (b *Bucket) Add(val float64) {
	b.Point += val
	b.Count++
}

// 重置桶
func (b *Bucket) Reset() {
	b.Point = 0
	b.Count = 0
}

// 返回下一个桶
func (b *Bucket) Next() *Bucket {
	return b.next
}

// 包含多个桶的窗口
type Window struct {
	window []Bucket
	size   int
}

// 窗口配置
type WindowOpt struct {
	Size int
}

// 新建一个窗口
func NewWindow(opt WindowOpt) *Window {
	buckets := make([]Bucket, opt.Size)
	//创建环形数组
	for i := range buckets {
		buckets[i] = Bucket{}
		nextOffset := i + 1
		if nextOffset == opt.Size {
			nextOffset = 0
		}
		buckets[i].next = &buckets[nextOffset]
	}
	return &Window{window: buckets, size: opt.Size}
}

// 重置整个窗口
func (w *Window) ResetWindow() {
	for i := range w.window {
		w.ResetBucket(i)
	}
}

// 重置指定index的单个桶
func (w *Window) ResetBucket(offset int) {
	w.window[offset].Reset()
}

// 重置一组指定index的桶
func (w *Window) ResetBuckets(offsets []int) {
	for _, offset := range offsets {
		w.ResetBucket(offset)
	}
}

// 向窗口中指定index的桶添加数据
func (w *Window) Add(offset int, val float64) {
	w.window[offset].Add(val)
}

// 返回指定index的桶
func (w *Window) Bucket(offset int) Bucket {
	return w.window[offset]
}

// 获取窗口的大小
func (w *Window) Size() int {
	return w.size
}

// 生成窗口的桶迭代器
func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count: count,
		cur:   &w.window[offset],
	}
}
