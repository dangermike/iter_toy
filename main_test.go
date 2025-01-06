package main

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumbers(t *testing.T) {
	agg, err := reduce(limit(numbersFrom(1), 100), 0, func(agg int, val int) (int, error) { return agg + val, nil })
	require.NoError(t, err)
	require.Equal(t, 5050, agg)
}

func BenchmarkNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range limit(numbers(), 10000) {
		}
	}
}

func BenchmarkReduceNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reduce(limit(numbersFrom(1), 10000), 0, func(agg int, val int) (int, error) { return agg + val, nil })
	}
}

func BenchmarkReduceIntNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reduce(limit(intsFrom(1), 10000), 0, func(agg int, val int) (int, error) { return agg + val, nil })
	}
}

func BenchmarkReduceIntInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reduceInt(limit(intsFrom(1), 10000), 0, func(agg int, val int) (int, error) { return agg + val, nil })
	}
}

func BenchmarkReduceIntIntFromTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reduceInt(intsFromTo(1, 10000), 0, func(agg int, val int) (int, error) { return agg + val, nil })
	}
}

func BenchmarkReduceIntNoErrIntFromTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reduceIntNoErr(intsFromTo(1, 10000), 0, func(agg int, val int) int { return agg + val })
	}
}

func BenchmarkLoopIntsFromTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sum int
		for v := range intsFromTo(1, 10000) {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkLoopNoIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sum int
		for i := 0; i < 10000; i++ {
			sum += i
		}
		_ = sum
	}
}

func TestFibs(t *testing.T) {
	arr := reduceToSlice(limit(fibs(), 10))
	require.Equal(t, []*big.Int{
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(8),
		big.NewInt(13),
		big.NewInt(21),
		big.NewInt(34),
		big.NewInt(55),
	}, arr)
}

func TestZip(t *testing.T) {
	fibSlice := reduceToSlice(limit(fibs(), 200))
	var exp int
	for i, f := range zip(limit(numbers(), 150), limit(fibs(), 200)) {
		require.Equal(t, exp, i)
		require.Equal(t, fibSlice[i], f)
		exp++
	}

	require.Equal(t, 150, exp)
}
