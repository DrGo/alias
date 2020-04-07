package alias

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func ExampleDraw() {
	a, err := NewFreq([]int{10, 10, 20, 40, 20})
	if err != nil {
		panic(err)
	}
	counts := make(map[int]int)
	for i := 0; i < 100000; i++ {
		counts[a.Draw()]++
	}
	fmt.Println(counts)
	// output:

}

func TestProbability(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	var params = []struct {
		freqs []int
		rates []int
	}{
		{[]int{10, 15}, []int{40, 60}},
		{[]int{20, 30}, []int{40, 60}},
		{[]int{20, 5}, []int{80, 20}},
		{[]int{25}, []int{100}},
		{[]int{1, 99}, []int{1, 99}},
		{[]int{1, 1, 8}, []int{10, 10, 80}},
	}

	for id, param := range params {
		sample := 1000000
		results := map[int]int{}
		a, err := NewFreq(param.freqs)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < sample; i++ {
			r := a.Draw()
			results[r]++
		}

		for key, rate := range param.rates {
			count := results[key]

			p := float64(rate) / 100
			q := 1.0 - p

			expected := float64(sample) * p
			// 3.89 = inverse of normal distribution function with alpha=0.9999
			delta := 3.89 * math.Sqrt(expected*q)

			if !(expected-delta <= float64(count) && float64(count) <= expected+delta) {
				w := param.freqs[key]
				t.Errorf("[%d] The probability is out of by interval estimation. key=%d weight=%d actual=%d, expected=%f, delta=%f", id, key, w, count, expected, delta)
			}
		}
	}
}
