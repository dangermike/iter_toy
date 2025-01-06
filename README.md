# Iter Toy - Some fun with custom iterators in Go

Since the release of Go 1.23, custom iterators have been available using the function/callback signatures in the [`iter`](https://pkg.go.dev/iter) package. The [`iter.Seq`](https://pkg.go.dev/iter#Seq) function has the following signature:

```go
type Seq[V any] func(yield func(V) bool)
```

## Making toys

This was meant to be a toy and I would not reach for the sort-of-functional style that I used in here in Go. It works, but it feels a little too cute and doesn't read like Go. That said, here's what I threw together:

```go
func reduce[T any, U any](src iter.Seq[T], initial U, combine func(agg U, val T) (U, error)) (U, error)
func numbers() iter.Seq[int]
func numbersFrom(i int) iter.Seq[int]
func fibs() iter.Seq[*big.Int]
func limit[T any](src iter.Seq[T], cnt int) iter.Seq[T]
func limit2[T any, U any](src iter.Seq2[T, U], cnt int) iter.Seq2[T, U]
func zip[T any, U any](a iter.Seq[T], b iter.Seq[U]) iter.Seq2[T, U]
```

The fact that we're using generics here reminds me of [this excellent article by Vicent Marti](https://planetscale.com/blog/generics-can-make-your-go-code-slower) about how generics create additional virtualization (allocations) when interfaces are involved[¹](#note1). I was going to compare what happened when using a plain type (`int`) and some interface. However, that didn't happen.

The first thing I did was just print some stuff to `stdout`, which worked fine. Then I threw together some tests, which highlighted a problem when passing around `*big.Int`, but was otherwise uneventful. Now that I knew it was working, I wrote some quick benchmarks.

## Wat?! (BenchmarkNumbers and BenchmarkReduceNumbers)

As per usual, I ran my benchmarks with `-benchmem`. I first just iterated a bunch of numbers, showing that the `for ... range` call over a `limit`ed sequence of numbers didn't do any allocations. This is what I expected and this is what happened. From there, I used the `reduce` function to sum up a `limit`ed sequence of values. There were suddenly 7 allocations! Again, this is all callbacks and I again expected no allocations at all.

## Maybe it's a type thing?

With Vincent Marti still whispering in my ear, I thought it might help to start taking the generics out. I made an `int`-specific `reduceInt` function. I had even tried taking the `limit` function out entirely by creating an `intsFromTo` generator. These are used in the

* `BenchmarkReduceIntNumbers` (`reduce(limit(ints))`)
* `BenchmarkReduceIntInt` ` (reduceInt(limit(ints))`)
* `BenchmarkReduceIntIntFromTo` `reduceInt(intsFromTo)`

Well, all of these had the same 7 allocations that got is here. So it's not a type thing. Incidentally, `BenchmarkReduceIntInt` was also significantly slower than the other versions.

## Se7en

I was aggregating 10,000 `int`s per cycle, so where was 7 coming from? On a whim, I tried changing the `reduce` function. Since there is no `error` returned when adding `int`s, I made a version of the `reduce` function that didn't return an error:

```go
func reduceIntNoErr(src func(yield func(int) bool), initial int, combine func(agg int, val int) int) int
```

As tested in `BenchmarkReduceIntNoErrIntFromTo`, this not only reduced the allocation count from 7 to 3, it also was nearly twice as fast.

## Conclusion: Should you write code like this?

As mentioned above, this doesn't read like Go. The `reduce` function used in  `BenchmarkReduceNumbers` is 13x slower than a simple loop over the `int` generator and the "tuned" `BenchmarkReduceIntNoErrIntFromTo` is still 7.5x slower. Taking the iteration out completely and doing it all with just a local loop doubles the speed again.

In other words, if you're thinking of replacing loops with a bunch of iterators calling iterators, just don't. You're fighting the language and it shows up in the benchmarks.

If you've got some other custom form of iteration (callbacks, iterator objects, or passing around channels), the standard `Seq` and `Seq2` functions are great. You're not going to get any worse performance and your callers get to use `for ... range`.

## Benchmark results

```plaintext
goos: darwin
goarch: amd64
pkg: github.com/dangermike/iter_toy
cpu: VirtualApple @ 2.50GHz
BenchmarkNumbers-10                    	   22617	     53053 ns/op	       0 B/op	       0 allocs/op
BenchmarkReduceNumbers-10              	   16316	     73545 ns/op	     144 B/op	       7 allocs/op
BenchmarkReduceIntNumbers-10           	   15338	     78362 ns/op	     144 B/op	       7 allocs/op
BenchmarkReduceIntInt-10               	   12376	     96826 ns/op	     128 B/op	       7 allocs/op
BenchmarkReduceIntIntFromTo-10         	   25448	     47433 ns/op	     128 B/op	       7 allocs/op
BenchmarkReduceIntNoErrIntFromTo-10    	   23570	     49921 ns/op	      48 B/op	       3 allocs/op
BenchmarkLoopIntsFromTo-10             	  179073	      6632 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoopNoIter-10                 	  378640	      3124 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/dangermike/iter_toy	14.265s
```

## Notes

<a id="note1">**1.**</a> Seriously, if you haven't read that article, you should.
