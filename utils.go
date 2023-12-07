package lttb

import (
	"math"
)

// 计算三角形面积
func calculateTriangleArea(pa, pb, pc Point) float64 {
	area := (float64(pa.Timestamp-pc.Timestamp)*(pb.Value-pa.Value) - float64(pa.Timestamp-pb.Timestamp)*(pc.Value-pa.Value)) * 0.5
	return math.Abs(area)
}

// 计算平均数点
func calculateAverageDataPoint(points []Point) (avg Point) {
	for _, point := range points {
		avg.Timestamp += point.Timestamp
		avg.Value += point.Value
	}
	l := len(points)
	avg.Timestamp /= uint64(l)
	avg.Value /= float64(l)
	return avg
}

// 分割桶
func splitDataBucket(data []Point, threshold int) [][]Point {

	buckets := make([][]Point, threshold)
	for i := range buckets {
		buckets[i] = make([]Point, 0)
	}
	// 第一个和最后一个桶由第一个和最后一个数据点组成
	buckets[0] = append(buckets[0], data[0])
	buckets[threshold-1] = append(buckets[threshold-1], data[len(data)-1])

	bucketSize := float64(len(data)-2) / float64(threshold-2)

	// 删除第一个和最后一个点
	d := data[1 : len(data)-1]

	for i := 0; i < threshold-2; i++ {
		bucketStartIdx := int(math.Floor(float64(i) * bucketSize))
		bucketEndIdx := int(math.Floor(float64(i+1)*bucketSize)) + 1
		if i == threshold-3 {
			bucketEndIdx = len(d)
		}
		buckets[i+1] = append(buckets[i+1], d[bucketStartIdx:bucketEndIdx]...)
	}

	return buckets
}

// 计算平均点
func calculateAveragePoint(points []Point) Point {
	l := len(points)
	var p Point
	for i := 0; i < l; i++ {
		p.Timestamp += points[i].Timestamp
		p.Value += points[i].Value
	}

	p.Timestamp /= uint64(l)
	p.Value /= float64(l)
	return p

}
