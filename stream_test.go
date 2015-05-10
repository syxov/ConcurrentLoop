package stream

import (
	"math"
	"runtime"
	"testing"
)

func TestEach(t *testing.T) {
	slice := []int{1, 2, 3}
	testSlice := make([]int, 3, 3)
	Each(slice, func(s int, index int) {
		testSlice[index] = int(math.Pow(float64(s), 2))
	})
	if len(testSlice) != 3 || (testSlice[0] != 1 || testSlice[1] != 4 || testSlice[2] != 9) {
		t.Error("Fail")
	}
}

func TestMap(t *testing.T) {
	testResult := Map([]int{1, 2, 3}, func(s int, index int) float64 {
		return float64(s)
	}).([]float64)
	if len(testResult) != 3 || (testResult[0] != 1.0 || testResult[1] != 2.0 || testResult[2] != 3.0) {
		t.Error("Fail")
	}
}

func BenchmarkEachConcurrent(b *testing.B) {
	slice := []int{1000, 2000, 3000, 4000, 5000, 6000}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Each(slice, func(s int, index int) int {
			var i = 0
			square := s * s
			for i < square {
				i++
			}
			return i
		})
	}
}

func BenchmarkEach(b *testing.B) {
	slice := []int{1000, 2000, 3000, 4000, 5000, 6000}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for index, i := range slice {
			func(s, index int) int {
				var i = 0
				square := s * s
				for i < square {
					i++
				}
				return i
			}(i, index)
		}
	}
}
