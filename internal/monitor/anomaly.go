package monitor

import (
	"fmt"
	"math"
)

func isAnomaly(currentLatency int64, history []int64) (bool, string) {
	if len(history) < 10 {
		return false, ""
	}

	var sum int64
	for _, v := range history {
		sum += v
	}

	mean := float64(sum) / float64(len(history))

	var varianceSum float64
	for _, v := range history {
		varianceSum += math.Pow(float64(v)-mean, 2)
	}
	stdDev := math.Sqrt(varianceSum / float64(len(history)))

	if stdDev < 10.0 {
		stdDev = 10.0
	}

	zScore := (float64(currentLatency) - mean) / stdDev

	thresholdZ := 3.0

	if zScore > thresholdZ {
		return true, fmt.Sprintf("Anomaly detected! Latency %.0fms is abnormal (Mean %.0fms, Z-Score> %.2f)", float64(currentLatency), mean, zScore)
	}

	return false, ""
}
