//go:build ignore

package main

import (
	"fmt"
	"time"
)

type Sorter struct {
	comparisons int
	swaps       int
}

func (s *Sorter) swap(arr []float64, i, j int) {
	if i != j {
		arr[i], arr[j] = arr[j], arr[i]
		s.swaps++
	}
}

func (s *Sorter) partition(arr []float64, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		s.comparisons++
		if arr[j] <= pivot {
			i++
			s.swap(arr, i, j)
		}
	}
	s.swap(arr, i+1, high)
	return i + 1
}

func (s *Sorter) sort(arr []float64) {
	stack := []int{0, len(arr) - 1}
	for len(stack) > 0 {
		high := stack[len(stack)-1]
		low := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		if low < high {
			p := s.partition(arr, low, high)
			if p-1 > low {
				stack = append(stack, low, p-1)
			}
			if p+1 < high {
				stack = append(stack, p+1, high)
			}
		}
	}
}

func main() {
	data150 := []float64{655, 115, 26, 760, 282, 251, 229, 143, 755, 105, 693, 759, 914, 559, 90, 605, 433, 33, 31, 96, 224, 239, 518, 617, 28, 575, 204, 734, 666, 719, 559, 430, 226, 460, 604, 285, 829, 891, 7, 778, 826, 164, 715, 433, 349, 285, 160, 221, 981, 782, 345, 105, 95, 390, 100, 368, 868, 353, 619, 271, 827, 45, 748, 471, 550, 128, 997, 945, 388, 81, 566, 301, 850, 644, 634, 907, 883, 371, 592, 197, 722, 72, 47, 678, 234, 792, 297, 82, 876, 239, 888, 104, 390, 285, 465, 651, 855, 374, 167, 380, 364, 215, 687, 274, 719, 960, 700, 664, 74, 624, 651, 176, 547, 747, 251, 168, 474, 389, 277, 948, 656, 705, 571, 225, 702, 333, 864, 787, 795, 58, 235, 842, 33, 825, 324, 411, 275, 68, 217, 936, 966, 581, 898, 736, 323, 218, 672, 512, 406, 906}
	s1 := &Sorter{}
	t0 := time.Now()
	s1.sort(data150)
	d1 := time.Since(t0)
	fmt.Printf("Go 150: min=%.0f max=%.0f comps=%d time=%.6fs\n", data150[0], data150[len(data150)-1], s1.comparisons, d1.Seconds())

	data10k := make([]float64, 0, 10000)
	for i := 10000; i >= 1; i-- {
		data10k = append(data10k, float64(i))
	}
	s2 := &Sorter{}
	t1 := time.Now()
	s2.sort(data10k)
	d2 := time.Since(t1)
	fmt.Printf("Go 10000 worst: min=%.0f max=%.0f comps=%d time=%.6fs\n", data10k[0], data10k[len(data10k)-1], s2.comparisons, d2.Seconds())
}
