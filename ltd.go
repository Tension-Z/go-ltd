package lttb

import (
	"math"
)

// 参考论文：https://skemman.is/bitstream/1946/15343/3/SS_MSthesis.pdf

// calculateLinearRegressionCoefficients 计算 点 到 平均点 的线性回归系数
func calculateLinearRegressionCoefficients(points []Point) (float64, float64) {
	average := calculateAveragePoint(points)
	aNumerator := 0.0
	aDenominator := 0.0
	for i := 0; i < len(points); i++ {
		aNumerator += float64(points[i].Timestamp-average.Timestamp) * (points[i].Value - average.Value)
		aDenominator += (points[i].Value - average.Value) * float64(points[i].Timestamp-average.Timestamp)
	}

	a := aNumerator / aDenominator
	b := average.Value - a*float64(average.Timestamp)
	return a, b
}

// calculateSSEForBucket 计算桶的SSE
func calculateSSEForBucket(points []Point) float64 {
	a, b := calculateLinearRegressionCoefficients(points)
	sumStandardErrorsSquared := 0.0
	for _, p := range points {
		standardError := p.Value - (a*float64(p.Timestamp) + b)
		sumStandardErrorsSquared += standardError * standardError
	}
	return sumStandardErrorsSquared
}

// calculateSSEForBuckets 计算多个桶的SSE
func calculateSSEForBuckets(buckets [][]Point) []float64 {
	sse := make([]float64, len(buckets)-2)

	// 跳过第一个和最后一个桶，因为它们只包含一个数据点
	for i := 1; i < len(buckets)-1; i++ {
		prevBucket := buckets[i-1]
		currBucket := buckets[i]
		nextBucket := buckets[i+1]
		// var bucketWithAdjacentPoints []Point
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, prevBucket[len(prevBucket)-1])
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, currBucket...)
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, nextBucket[0])
		bucketWithAdjacentPoints := make([]Point, len(currBucket)+2)
		bucketWithAdjacentPoints[0] = prevBucket[len(prevBucket)-1]
		bucketWithAdjacentPoints[len(bucketWithAdjacentPoints)-1] = nextBucket[0]
		for i := 1; i < len(currBucket); i++ {
			bucketWithAdjacentPoints[i] = currBucket[i-1]
		}
		sse[i-1] = calculateSSEForBucket(bucketWithAdjacentPoints)
	}

	sse = append(sse, 0)
	return sse
}

// findLowestSSEAdjacentBucketsIndex 找到SSE最低的相邻桶的索引
func findLowestSSEAdjacentBucketIndex(sse []float64, ignoreIndex int) int {
	minSSE := float64(math.MaxInt64)
	minSSEIndex := -1
	for i := 1; i < len(sse)-2; i++ {
		if i == ignoreIndex || i+1 == ignoreIndex {
			continue
		}

		if sse[i]+sse[i+1] < minSSE {
			minSSE = sse[i] + sse[i+1]
			minSSEIndex = i
		}
	}
	return minSSEIndex
}

// findHighestSSEBucketIndex 找到SSE最高的相邻桶的索引
func findHighestSSEBucketIndex(buckets [][]Point, sse []float64) int {
	maxSSE := 0.0
	maxSSEIdx := -1
	for i := 1; i < len(sse)-1; i++ {
		if len(buckets[i]) > 1 && sse[i] > maxSSE {
			maxSSE = sse[i]
			maxSSEIdx = i
		}
	}
	return maxSSEIdx
}

// splitBucketAt 将桶分成大致相等的两个桶。如果桶包含奇数个点，则一个桶将包含比另一个桶多一个点
func splitBucketAt(buckets [][]Point, index int) [][]Point {
	if index < 0 || index >= len(buckets) {
		return buckets
	}
	bucket := buckets[index]
	bucketSize := len(bucket)
	if bucketSize < 2 {
		return buckets
	}

	bucketALength := int(math.Ceil(float64(bucketSize / 2)))
	bucketA := bucket[0 : bucketALength+1]
	bucketB := bucket[bucketALength:]

	var newBuckets [][]Point
	newBuckets = append(newBuckets, buckets[0:index]...)
	newBuckets = append(newBuckets, bucketA, bucketB)
	newBuckets = append(newBuckets, buckets[index+1:]...)

	return newBuckets
}

// mergeBucketAt 将两个桶合并为一个桶
func mergeBucketAt(buckets [][]Point, index int) [][]Point {
	if index < 0 || index >= len(buckets)-1 {
		return buckets
	}
	mergedBucket := buckets[index]
	mergedBucket = append(mergedBucket, buckets[index+1]...)

	var newBuckets [][]Point
	newBuckets = append(newBuckets, buckets[0:index]...)
	newBuckets = append(newBuckets, mergedBucket)
	newBuckets = append(newBuckets, buckets[index+2:]...)

	return newBuckets
}

// LTD - Largest triangle dynamic(LTD) 数据降采样算法的实现
//   - data . 原始数据
//   - threshold . 要返回的数据点数量
func LTD(data []Point, threshold int) []Point {

	if threshold >= len(data) || threshold == 0 {
		return data // Nothing to do
	}

	// 1.将数据拆分为与阈值相同数量的桶，第一个和最后一个桶由第一个和最后一个数据点组成
	buckets := splitDataBucket(data, threshold)
	numIterations := len(data) * 10 / threshold
	for iter := 0; iter < numIterations; iter++ {

		// 2: 相应的计算桶的SSE。
		sseForBuckets := calculateSSEForBuckets(buckets)

		// 4: 找到SSE最高的桶
		highestSSEBucketIndex := findHighestSSEBucketIndex(buckets, sseForBuckets)
		if highestSSEBucketIndex < 0 {
			break
		}

		// 5: 找到SSE和最小的一对相邻桶。
		lowestSSEAdjacentBucketIndex := findLowestSSEAdjacentBucketIndex(sseForBuckets, highestSSEBucketIndex)
		if lowestSSEAdjacentBucketIndex < 0 {
			break
		}

		// 6: 将桶分成大致相等的两个桶。如果桶F包含奇数个点，则一个桶将包含比另一个桶多一个点
		buckets = splitBucketAt(buckets, highestSSEBucketIndex)

		// 7: 将两个桶合并为一个桶
		if lowestSSEAdjacentBucketIndex > highestSSEBucketIndex {
			lowestSSEAdjacentBucketIndex++
		}
		buckets = mergeBucketAt(buckets, lowestSSEAdjacentBucketIndex)

	}
	// 对生成的桶使用最大三角三桶算法
	return LTTBForBuckets(buckets)
}
