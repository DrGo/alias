package aliasmethod

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type AliasMethod struct {
	rand *rand.Rand
}

func NewAliasMethod() *AliasMethod {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &AliasMethod{rand: r}
}

func (self *AliasMethod) Random(table *AliasTable) int {
	u := self.rand.Float64()
	n := self.rand.Intn(table.Len)

	if u <= table.Prob[n] {
		return int(n)
	} else {
		return table.Alias[n]
	}
}

type AliasTable struct {
	Len   int
	Prob  []float64
	Alias []int
}

func NewAliasTable(weights []int) (*AliasTable, error) {
	n := len(weights)

	sum := 0
	for _, value := range weights {
		sum += value
	}
	if sum == 0 {
		return nil, errors.New("sum of weights is 0.")
	}

	prob := make([]float64, n)
	for i, w := range weights {
		prob[i] = float64(w) * float64(n) / float64(sum)
	}

	h := 0
	l := n - 1
	hl := make([]int, n)
	for i, p := range prob {
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

		if 1 < prob[j] {
			panic(fmt.Sprintf("MUST: %f <= 1", prob[j]))
		}
		if prob[k] < 1 {
			panic(fmt.Sprintf("MUST: 1 <= %f", prob[k]))
		}

		a[j] = k
		prob[k] -= (1 - prob[j]) // - residual weight
		l++
		if prob[k] < 1 {
			hl[l] = k
			l--
			h--
		}
	}

	return &AliasTable{Len: n, Prob: prob, Alias: a}, nil
}
