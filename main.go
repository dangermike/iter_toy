package main

import (
	"fmt"
	"iter"
	"math/big"
	"slices"
)

func reduce[T any, U any](src iter.Seq[T], initial U, combine func(agg U, val T) (U, error)) (U, error) {
	agg := initial
	var err error
	for val := range src {
		if agg, err = combine(agg, val); err != nil {
			return agg, err
		}
	}
	return agg, err
}

func reduceInt(src func(yield func(int) bool), initial int, combine func(agg int, val int) (int, error)) (int, error) {
	agg := initial
	var err error
	for val := range src {
		if agg, err = combine(agg, val); err != nil {
			return agg, err
		}
	}
	return agg, err
}
func reduceIntNoErr(src func(yield func(int) bool), initial int, combine func(agg int, val int) int) int {
	agg := initial
	for val := range src {
		agg = combine(agg, val)
	}
	return agg
}

func reduceToSlice[T any](src iter.Seq[T]) []T {
	slice, _ := reduce(src, nil, func(agg []T, val T) ([]T, error) {
		return append(agg, val), nil
	})
	return slice
}

func numbers() iter.Seq[int] {
	return numbersFrom(0)
}

func numbersFrom(i int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for yield(i) {
			i++
		}
	}
}

func intsFrom(i int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for yield(i) {
			i++
		}
	}
}
func intsFromTo(i int, j int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		ok := true
		for ; ok && i < j; i++ {
			ok = yield(i)
		}
	}
}

func fibs() iter.Seq[*big.Int] {
	var a, b *big.Int = big.NewInt(1), big.NewInt(1)
	return func(yield func(*big.Int) bool) {
		if !yield(a) {
			return
		}
		for {
			if !yield(b) {
				return
			}
			// Need new value for b so we don't mess up the returned value. That means
			// we will allocate on every iteration, but at least it's safe.
			a, b = b, new(big.Int).Add(a, b)
		}
	}
}

func limit[T any](src iter.Seq[T], cnt int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range src {
			if cnt <= 0 {
				break
			}
			cnt--
			if !yield(v) {
				break
			}
		}
	}
}

func limit2[T any, U any](src iter.Seq2[T, U], cnt int) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		for k, v := range src {
			if cnt <= 0 {
				break
			}
			cnt--
			if !yield(k, v) {
				break
			}
		}
	}
}

func zip[T any, U any](a iter.Seq[T], b iter.Seq[U]) iter.Seq2[T, U] {
	aNext, aStop := iter.Pull(a)
	bNext, bStop := iter.Pull(b)
	return func(yield func(T, U) bool) {
		defer aStop()
		defer bStop()
		aV, aOK := aNext()
		bV, bOK := bNext()
		for aOK && bOK {
			if !yield(aV, bV) {
				return
			}
			aV, aOK = aNext()
			bV, bOK = bNext()
		}
	}
}

func main() {
	fmt.Println("limit on iter")
	for i := range limit(numbers(), 3) {
		fmt.Println(i)
	}
	fmt.Println()
	fmt.Println("limit on slices.Values")
	for k := range limit(slices.Values([]string{"a", "b", "c", "d"}), 3) {
		fmt.Println(k)
	}
	fmt.Println()
	fmt.Println("limit on slices.All")
	for k, v := range limit2(slices.All([]string{"a", "b", "c", "d"}), 3) {
		fmt.Println(k, v)
	}
	fmt.Println()
	fmt.Println("zip")
	for a, b := range zip(numbers(), slices.Values([]string{"a", "b", "c", "d"})) {
		fmt.Println(a, b)
	}
	fmt.Println()
	fmt.Println("fibs")
	for i, f := range zip(numbers(), limit(fibs(), 200)) {
		fmt.Println(i, f)
	}
}
