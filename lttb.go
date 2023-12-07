package lttb

import (
	"math"
)

func LTTB(data []Point, threshold int) []Point {

	if threshold >= len(data) || threshold == 0 {
		return data // Nothing to do
	}

	sampledData := make([]Point, 0, threshold)

	// 桶大小. 为开始和结束数据点留出空间
	bucketSize := float64(len(data)-2) / float64(threshold-2)

	sampledData = append(sampledData, data[0]) // 添加第一个点

	bucketLow := 1
	bucketMiddle := int(math.Floor(bucketSize)) + 1

	var prevMaxAreaPoint int

	for i := 0; i < threshold-2; i++ {

		bucketHigh := int(math.Floor(float64(i+2)*bucketSize)) + 1
		if bucketHigh >= len(data)-1 {
			bucketHigh = len(data) - 2
		}

		// 计算下一个桶的平均点
		avgPoint := calculateAverageDataPoint(data[bucketMiddle : bucketHigh+1])

		// 获取当前桶的范围
		currBucketStart := bucketLow
		currBucketEnd := bucketMiddle

		pointA := data[prevMaxAreaPoint]

		maxArea := -1.0

		var maxAreaPoint int
		for ; currBucketStart < currBucketEnd; currBucketStart++ {

			area := calculateTriangleArea(pointA, avgPoint, data[currBucketStart])
			if area > maxArea {
				maxArea = area
				maxAreaPoint = currBucketStart
			}
		}
		sampledData = append(sampledData, data[maxAreaPoint])
		prevMaxAreaPoint = maxAreaPoint

		// 移至下一个窗口
		bucketLow = bucketMiddle
		bucketMiddle = bucketHigh
	}

	sampledData = append(sampledData, data[len(data)-1]) // 添加到最后

	return sampledData
}

// LTTBForBuckets - 在每个桶上应用 LTTB 算法
func LTTBForBuckets(buckets [][]Point) []Point {
	bucketCount := len(buckets)
	sampledData := make([]Point, 0)

	sampledData = append(sampledData, buckets[0][0])

	lastSelectedDataPoint := buckets[0][0]
	for i := 1; i < bucketCount-1; i++ {
		bucket := buckets[i]
		averagePoint := calculateAveragePoint(buckets[i+1])

		maxArea := -1.0
		maxAreaIndex := -1
		for j := 0; j < len(bucket); j++ {
			point := bucket[j]
			area := calculateTriangleArea(lastSelectedDataPoint, point, averagePoint)

			if area > maxArea {
				maxArea = area
				maxAreaIndex = j
			}
		}
		lastSelectedDataPoint := bucket[maxAreaIndex]
		sampledData = append(sampledData, lastSelectedDataPoint)
	}
	sampledData = append(sampledData, buckets[len(buckets)-1][0])
	return sampledData
}
