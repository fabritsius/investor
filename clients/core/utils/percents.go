package utils

import "fmt"

func GetPercDiff(old float64, new float64) string {
	percDiff := (new - old) / old * 100
	if percDiff >= 0 {
		return fmt.Sprintf("+%3.2f", percDiff)
	}

	return fmt.Sprintf("%3.2f", percDiff)
}
