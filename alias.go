// Package alias implements random sampling from a cumulative distribution using the Walker-Vose alias method,
package alias

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Alias struct {
	rand  *rand.Rand
	table *table
}

func NewFreq(freqs []int) (*Alias, error) {
	probs, err := freqToProbs(freqs)
	if err != nil {
		return nil, err
	}
	return New(probs)
}

func New(probs []float64) (*Alias, error) {
	t, err := newTable(probs)
	if err != nil {
		return nil, err
	}
	return &Alias{
		table: t,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano()))}, nil
}

func (alias *Alias) Draw() int {
	u := alias.rand.Float64()
	n := alias.rand.Intn(alias.table.len)
	if u <= alias.table.prob[n] {
		return int(n)
	}
	return alias.table.alias[n]
}

type table struct {
	prob  []float64
	alias []int
	len   int
}

// freqToProbs converts a freq distribution into an probability distribution
func freqToProbs(freqs []int) ([]float64, error) {
	sum := 0
	for _, value := range freqs {
		sum += value
	}
	if sum == 0 {
		return nil, errors.New("sum of freqs must be > 0")
	}
	n := len(freqs)
	probs := make([]float64, n)
	for i, w := range freqs {
		probs[i] = float64(w*n) / float64(sum)
	}
	return probs, nil
}

func newTable(probs []float64) (*table, error) {
	n := len(probs)
	h := 0
	l := n - 1
	hl := make([]int, n)
	for i, p := range probs {
		if p < 1 {
			hl[l] = i
			l--
		}
		if p > 1 {
			hl[h] = i
			h++
		}
	}

	a := make([]int, n)
	for h != 0 && l != n-1 {
		j := hl[l+1]
		k := hl[h-1]

		if 1 < probs[j] {
			return nil, fmt.Errorf("MUST: %f <= 1", probs[j])
		}
		if probs[k] < 1 {
			return nil, fmt.Errorf("MUST: 1 <= %f", probs[k])
		}
		a[j] = k
		probs[k] -= (1 - probs[j]) // - residual weight
		l++
		if probs[k] < 1 {
			hl[l] = k
			l--
			h--
		}
	}
	return &table{prob: probs, alias: a, len: n}, nil
}
