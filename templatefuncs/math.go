package templatefuncs

import "fmt"

func mod(a, b int) int {
	return a % b
}

func percent(v1, v2 int) string {
	if v2 == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", float64(v1)/float64(v2)*100)
}

func intDiv(v1, v2 int) float64 {
	if v2 == 0 {
		return 0.0
	}
	return float64(v1) / float64(v2)
}
