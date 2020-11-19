package main

import (
	"time"
)

// Go needs generics for this case
func MinMaxDuration(array []time.Duration) (time.Duration, time.Duration) {
	var max time.Duration = array[0]
	var min time.Duration = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
func MinMaxInt(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
