package aliasmethod

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestNewAliasTable(t *testing.T) {

	var params = []struct {
		weights []int
	}{
		// {[]int{99, 1}},
		// {[]int{5, 5}},
		{[]int{1, 1, 8}},
	}

	for _, p := range params {
		sample_weights := p.weights
		actual, err := NewAliasTable(sample_weights)

		if err != nil {
			t.Fatal(err)
		}

		if actual.Len != len(sample_weights) {
			t.Error("AliasTable.Len does not equals length of weights.")
		}
		if len(actual.Prob) != len(sample_weights) {
			t.Error("Length of AliasTable.Prob does not equals length of weights.")
		}
	}

}

func TestProbability(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	var params = []struct {
		weights []int
		rates   []int
	}{
		{[]int{10, 15}, []int{40, 60}},
		{[]int{20, 30}, []int{40, 60}},
		{[]int{20, 5}, []int{80, 20}},
		{[]int{25}, []int{100}},
		{[]int{1, 99}, []int{1, 99}},
		{[]int{1, 1, 8}, []int{10, 10, 80}},
	}

	for param_no, param := range params {
		table, _ := NewAliasTable(param.weights)

		sample := 1000000.0
		results := map[int]int{}

		a := NewAliasMethod()
		for i := 0.0; i < sample; i++ {
			r := a.Random(table)
			results[r] += 1
		}

		for key, rate := range param.rates {
			count := results[key]

			p := float64(rate) / 100
			q := 1.0 - p

			expected := sample * p
			// 3.89 = inverse of normal distribution function with alpha=0.9999
			delta := 3.89 * math.Sqrt(sample*p*q)

			ok := expected-delta <= float64(count) && float64(count) <= expected+delta
			if !ok {
				w := param.weights[key]
				t.Errorf("[%d] The probability is out of by interval estimation. key=%d weight=%d actual=%d, expected=%f, delta=%f", param_no, key, w, count, expected, delta)
			}
		}
	}
}
