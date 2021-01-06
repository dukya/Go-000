package rolling

//包含各种聚合计算函数（参数为滑动窗口迭代器）

// 返回滑动窗口中数据之和
func Sum(iterator Iterator) float64 {
	result := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		result += bucket.Point
	}
	return result
}

// 返回滑动窗口中数据平均值
func Avg(iterator Iterator) float64 {
	result := 0.0
	count := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		result += bucket.Point
		if bucket.Count != 0 {
			count++
		}
	}
	return result / count
}

// 返回滑动窗口中数据的最小值
func Min(iterator Iterator) float64 {
	result := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		if bucket.Point < result {
			result = bucket.Point
		}
	}
	return result
}

// 返回滑动窗口中数据的最大值
func Max(iterator Iterator) float64 {
	result := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		if bucket.Point > result {
			result = bucket.Point
		}
	}
	return result
}

// Count sums the count value within the window.
func Count(iterator Iterator) float64 {
	var result int64
	for iterator.Next() {
		bucket := iterator.Bucket()
		result += bucket.Count
	}
	return float64(result)
}
