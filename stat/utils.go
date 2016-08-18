package statutils

import "sort"

//getPercentile generates percentiles for input arrray of floats
func getPercentile(numbers []float64) map[float64]float64 {
	sort.Float64s(numbers)

	var percentiles = make(map[float64]float64)

	var length = len(numbers)

	if length == 0 {
		return percentiles
	}

	var steps = 1000
	var step = 1.0 / float64(steps)

	for i := 0; i < steps; i++ {
		var percentile = float64(i) * step
		var index = int(float64(length) * percentile)
		//fmt.Fprintf(os.Stdout, "index: %d/%d\n", index, length)
		percentiles[percentile] = numbers[index]
	}

	return percentiles
}
