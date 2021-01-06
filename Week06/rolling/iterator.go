package rolling

import (
	"fmt"
)

// 滑动窗口的迭代器
type Iterator struct {
	count         int     //需要遍历的桶总数
	iteratedCount int     //已遍历的桶个数
	cur           *Bucket //指向当前正在遍历的桶
}

// 当遍历完所有桶后返回false
func (i *Iterator) Next() bool {
	return i.iteratedCount != i.count
}

// 获取当前正在遍历的桶
func (i *Iterator) Bucket() Bucket {
	if !(i.Next()) {
		panic(fmt.Errorf("iteration out of range iteratedCount: %d count: %d", i.iteratedCount, i.count))
	}
	bucket := *i.cur
	i.iteratedCount++    //遍历数自增
	i.cur = i.cur.Next() //指向下一个桶
	return bucket
}
