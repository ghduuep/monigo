package monitor

import (
	"fmt"
	"math"
	"slices"
)

func isAnomaly(currentLatency int64, history []int64) (bool, string) {
	if len(history) < 10 {
		return false, ""
	}

	sortedHistory := slices.Clone(history)
	slices.Sort(sortedHistory)

	median := calculateMedian(sortedHistory)

	var deviations []float64
	for _, val := range history {
		diff := math.Abs(float64(val) - median)
		deviations = append(deviations, diff)
	}

	slices.Sort(deviations)
	mad := deviations[len(deviations)/2]

	if mad < 1.0 {
		mad = 1.0
	}

	modifiedZScore := 0.6745 * (float64(currentLatency) - median) / mad

	threshold := 3.5

	if modifiedZScore > threshold {
		return true, fmt.Sprintf("Anomaly detected! Latency %.0fms (Median: %.0fms, Score: %.2f)", float64(currentLatency), median, modifiedZScore)
	}

	return false, ""
}

func calculateMedian(sorted []int64) float64 {
	l := len(sorted)
	if l == 0 {
		return 0
	}

	if l%2 == 0 {
		return float64(sorted[l/2-1]+sorted[l/2]) / 2.0
	}

	return float64(sorted[l/2])
}
