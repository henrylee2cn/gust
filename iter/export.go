package iter

import (
	"unicode/utf8"
	"unsafe"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

// FromIterable creates an iterator from an Iterable.
func FromIterable[T any](data gust.Iterable[T]) Iterator[T] {
	iter, _ := data.(Iterator[T])
	if iter != nil {
		return iter
	}
	return newIterator(fromIterable[T](data))
}

// EnumIterable creates an iterator with index from an Iterable.
func EnumIterable[T any](data gust.Iterable[T]) Iterator[KV[T]] {
	return Enumerate[T](FromIterable[T](data))
}

// FromDeIterable creates a double ended iterator from an Iterable.
func FromDeIterable[T any](data gust.DeIterable[T]) DeIterator[T] {
	iter, _ := data.(DeIterator[T])
	if iter != nil {
		return iter
	}
	return newDeIterator(fromDeIterable[T](data))
}

// EnumDeIterable creates a double ended iterator with index from an Iterable.
func EnumDeIterable[T any](data gust.DeIterable[T]) DeIterator[KV[T]] {
	return DeEnumerate[T](FromDeIterable[T](data))
}

// FromVec creates a double ended iterator from a slice.
func FromVec[T any](slice []T) DeIterator[T] {
	return NewIterableVec(slice).ToDeIterator()
}

// EnumVec creates a double ended iterator with index from a slice.
func EnumVec[T any](slice []T) DeIterator[KV[T]] {
	return DeEnumerate[T](FromVec[T](slice))
}

// FromElements creates a double ended iterator from a set of elements.
func FromElements[T any](elems ...T) DeIterator[T] {
	return NewIterableVec(elems).ToDeIterator()
}

// EnumElements creates a double ended iterator with index from a set of elements.
func EnumElements[T any](elems ...T) DeIterator[KV[T]] {
	return DeEnumerate[T](FromVec[T](elems))
}

// FromRange creates a double ended iterator from a range.
func FromRange[T digit.Integer](start T, end T, rightClosed ...bool) DeIterator[T] {
	return NewIterableRange[T](start, end, rightClosed...).ToDeIterator()
}

// EnumRange creates a double ended iterator with index from a range.
func EnumRange[T digit.Integer](start T, end T, rightClosed ...bool) DeIterator[KV[T]] {
	return DeEnumerate[T](FromRange[T](start, end, rightClosed...))
}

// FromChan creates an iterator from a channel.
func FromChan[T any](c chan T) Iterator[T] {
	return NewIterableChan[T](c).ToIterator()
}

// EnumChan creates an iterator with index from a channel.
func EnumChan[T any](c chan T) Iterator[KV[T]] {
	return Enumerate[T](FromChan(c))
}

// FromResult creates a double ended iterator from a result.
func FromResult[T any](ret gust.Result[T]) DeIterator[T] {
	return FromDeIterable[T](ret)
}

// EnumResult creates a double ended iterator with index from a result.
func EnumResult[T any](ret gust.Result[T]) DeIterator[KV[T]] {
	return EnumDeIterable[T](ret)
}

// FromOption creates a double ended iterator from an option.
func FromOption[T any](opt gust.Option[T]) DeIterator[T] {
	return FromDeIterable[T](opt)
}

// EnumOption creates a double ended iterator with index from an option.
func EnumOption[T any](opt gust.Option[T]) DeIterator[KV[T]] {
	return EnumDeIterable[T](opt)
}

// FromString creates a double ended iterator from a string.
func FromString[T ~byte | ~rune](s string) DeIterator[T] {
	if len(s) == 0 {
		return NewIterableVec[T]([]T{}).ToDeIterator()
	}
	const bn = rune(^byte(0))
	if bn == rune(^T(0)) {
		var rs = *(*[]T)(unsafe.Pointer(
			&struct {
				string
				Cap int
			}{s, len(s)},
		))
		return NewIterableVec[T](rs).ToDeIterator()
	}
	var rs = make([]T, 0, len(s))
	var b = *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rs = append(rs, T(r))
		b = b[size:]
	}
	return NewIterableVec[T](rs).ToDeIterator()
}

// EnumString creates a double ended iterator with index from a string.
func EnumString[T ~byte | ~rune](s string) DeIterator[KV[T]] {
	return DeEnumerate[T](FromString[T](s))
}

// TryFold a data method that applies a function as long as it returns
// successfully, producing a single, final value.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
//
// the checked sum of iAll the elements of the array
// var sum = FromVec(a).TryFold(0, func(acc int, x int) { return Ok(acc+x) });
//
// assert.Equal(t, sum, Ok(6));
func TryFold[T any, CB any](iter gust.Iterable[T], init CB, f func(CB, T) gust.SigCtrlFlow[CB]) gust.SigCtrlFlow[CB] {
	var accum = gust.SigContinue[CB](init)
	for {
		x := iter.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.UnwrapContinue(), x.Unwrap())
		if accum.IsBreak() {
			return accum
		}
	}
}

// Fold folds every element into an accumulator by applying an operation,
// returning the final
//
// `Fold()` takes two arguments: an initial value, and a closure with two
// arguments: an 'accumulator', and an element. The closure returns the value that
// the accumulator should have for the data iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the data, `Fold()`
// returns the accumulator.
//
// This operation is sometimes called 'iReduce' or 'inject'.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// Note: `Fold()`, and similar methods that traverse the entire data,
// might not terminate for infinite iterators, even on interfaces for which a
// result is determinable in finite time.
//
// Note: [`Reduce()`] can be used to use the first element as the initial
// value, if the accumulator type and item type is the same.
//
// Note: `Fold()` combines elements in a *left-associative* fashion. For associative
// operators like `+`, the order the elements are combined in is not important, but for non-associative
// operators like `-` the order will affect the final
//
// # Note to Implementors
//
// Several of the other (forward) methods have default implementations in
// terms of this one, so try to implement this explicitly if it can
// do something better than the default `for` loop implementation.
//
// In particular, try to have this call `Fold()` on the internal parts
// from which this data is composed.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
//
// the sum of iAll the elements of the array
// var sum = FromVec(a).Fold((0, func(acc int, x int) any { return acc + x });
//
// assert.Equal(t, sum, 6);
//
// Let's walk through each step of the iteration here:
//
// | element | acc | x | result |
// |---------|-----|---|--------|
// |         | 0   |   |        |
// | 1       | 0   | 1 | 1      |
// | 2       | 1   | 2 | 3      |
// | 3       | 3   | 3 | 6      |
//
// And so, our final result, `6`.
func Fold[T any, B any](iter gust.Iterable[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := iter.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum, x.Unwrap())
	}
}

// Map takes a closure and creates an iterator which calls that closure on each
// element.
//
// If you are good at thinking in types, you can think of `Map()` like this:
// If you have an iterator that gives you elements of some type `A`, and
// you want an iterator of some other type `B`, you can use `Map()`,
// passing a closure that takes an `A` and returns a `B`.
//
// `Map()` is conceptually similar to a [`for`] loop. However, as `Map()` is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use [`for`] than `Map()`.
//
// # Examples
//
// Basic usage:
//
// ```
// var a = []int{1, 2, 3};
//
// var iter = FromVec(a).Map(func(x)int{ return 2 * x});
//
// assert.Equal(iter.Next(), gust.Some(2));
// assert.Equal(iter.Next(), gust.Some(4));
// assert.Equal(iter.Next(), gust.Some(6));
// assert.Equal(iter.Next(), gust.None[int]());
// ```
// NOTE: `iter` will be cleared
func Map[T any, B any](iter Iterator[T], f func(T) B) Iterator[B] {
	r := newIterator(newMapIterator[T, B](iter.inner(), f))
	iter.clear()
	return r
}

// DeMap takes a closure and creates a double ended iterator which calls that closure on each
// element.
//
// If you are good at thinking in types, you can think of `DeMap()` like this:
// If you have an iterator that gives you elements of some type `A`, and
// you want an iterator of some other type `B`, you can use `DeMap()`,
// passing a closure that takes an `A` and returns a `B`.
//
// `DeMap()` is conceptually similar to a [`for`] loop. However, as `DeMap()` is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use [`for`] than `DeMap()`.
func DeMap[T any, B any](iter DeIterator[T], f func(T) B) DeIterator[B] {
	r := newDeIterator(newDeMapIterator(iter.deInner(), f))
	iter.clear()
	return r
}

// FilterMap creates an iterator that both filters and maps.
//
// The returned iterator yields only the `value`s for which the supplied
// closure returns `gust.Some(value)`.
// NOTE: `iter` will be cleared
func FilterMap[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) Iterator[B] {
	r := newIterator(newFilterMapIterator[T, B](iter.inner(), f))
	iter.clear()
	return r
}

// DeFilterMap creates a double ended iterator that both filters and maps.
//
// The returned iterator yields only the `value`s for which the supplied
// closure returns `gust.Some(value)`.
func DeFilterMap[T any, B any](iter DeIterator[T], f func(T) gust.Option[B]) DeIterator[B] {
	r := newDeIterator(newDeFilterMapIterator[T, B](iter.deInner(), f))
	iter.clear()
	return r
}

// FindMap applies function to the elements of data and returns
// the first non-none
//
// `FindMap(iter, f)` is equivalent to `FilterMap(iter, f).Next()`.
//
// # Examples
//
// var a = []string{"lol", "NaN", "2", "5"};
//
// var first_number = FromVec(a).FindMap(func(s A) Option[any]{ return Wrap[any](strconv.Atoi(s))});
//
// assert.Equal(t, first_number, gust.Some(2));
func FindMap[T any, B any](iter gust.Iterable[T], f func(T) gust.Option[B]) gust.Option[B] {
	for {
		x := iter.Next()
		if x.IsNone() {
			break
		}
		y := f(x.Unwrap())
		if y.IsSome() {
			return y
		}
	}
	return gust.None[B]()
}

// Zip 'Zips up' two iterators into a single iterator of pairs.
//
// `Zip()` returns a new iterator that will iterate over two other
// iterators, returning a tuple where the first element comes from the
// first iterator, and the second element comes from the second iterator.
//
// In other words, it zips two iterators together, into a single one.
//
// If either iterator returns [`gust.None[A]()`], [`Next`] from the zipped iterator
// will return [gust.None[A]()].
// If the zipped iterator has no more elements to return then each further attempt to advance
// it will first try to advance the first iterator at most one time and if it still yielded an item
// try to advance the second iterator at most one time.
// NOTE: `a` and `b` will be cleared
func Zip[A any, B any](a Iterator[A], b Iterator[B]) Iterator[gust.Pair[A, B]] {
	r := newIterator(newZipIterator[A, B](a.inner(), b.inner()))
	a.clear()
	b.clear()
	return r
}

// DeZip is similar to `Zip`, but it supports take elements starting from the back of the iterator.
func DeZip[A any, B any](a DeIterator[A], b DeIterator[B]) DeIterator[gust.Pair[A, B]] {
	r := newDeIterator(newDeZipIterator[A, B](a.deInner(), b.deInner()))
	a.clear()
	b.clear()
	return r
}

// TryRfold is the reverse version of [`Iterator[T].TryFold()`]: it takes
// elements starting from the back of the iterator.
func TryRfold[T any, CB any](iter gust.DeIterable[T], init CB, f func(CB, T) gust.SigCtrlFlow[CB]) gust.SigCtrlFlow[CB] {
	var accum = gust.SigContinue[CB](init)
	for {
		x := iter.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.UnwrapContinue(), x.Unwrap())
		if accum.IsBreak() {
			return accum
		}
	}
}

// Rfold is an iterator method that reduces the iterator's elements to a single,
// final value, starting from the back.
func Rfold[T any, B any](iter gust.DeIterable[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := iter.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum, x.Unwrap())
	}
}

// Enumerate creates an iterator that yields pairs of the index and the value.
// NOTE: `iter` will be cleared
func Enumerate[T any](iter Iterator[T]) Iterator[KV[T]] {
	r := newIterator(newEnumerateIterator(iter.inner()))
	iter.clear()
	return r
}

// DeEnumerate creates a double ended iterator that yields pairs of the index and the value.
func DeEnumerate[T any](iter DeIterator[T]) DeIterator[KV[T]] {
	r := newDeIterator(newDeEnumerateIterator(iter.deInner()))
	iter.clear()
	return r
}

// MapWhile creates an iterator that both yields elements based on a predicate and maps.
//
// `MapWhile()` takes a closure as an argument. It will call this
// closure on each element of the iterator, and yield elements
// while it returns [`Some`].
// NOTE: `iter` will be cleared
func MapWhile[T any, B any](iter Iterator[T], predicate func(T) gust.Option[B]) Iterator[B] {
	r := newIterator(newMapWhileIterator[T, B](iter, predicate))
	iter.clear()
	return r
}

// Scan is an iterator adapter similar to [`Fold`] that holds internal state and
// produces a new iterator.
//
// [`Fold`]: Iterator.Fold
//
// `Scan()` takes two arguments: an initial value which seeds the internal
// state, and a closure with two arguments, the first being a mutable
// reference to the internal state and the second an iterator element.
// The closure can assign to the internal state to share state between
// iterations.
//
// On iteration, the closure will be applied to each element of the
// iterator and the return value from the closure, an [`Option`], is
// yielded by the iterator.
// NOTE: `iter` will be cleared
func Scan[T any, St any, B any](iter Iterator[T], initialState St, f func(state *St, item T) gust.Option[B]) Iterator[B] {
	r := newIterator(newScanIterator[T, St, B](iter, initialState, f))
	iter.clear()
	return r
}

// Flatten creates an iterator that flattens nested structure.
// NOTE: `iter` will be cleared
func Flatten[I gust.Iterable[T], T any](iter Iterator[I]) Iterator[T] {
	r := newIterator[T](newFlattenIterator[I, T](iter.inner()))
	iter.clear()
	return r
}

// DeFlatten creates a double ended iterator that flattens nested structure.
func DeFlatten[I gust.DeIterable[T], T any](iter DeIterator[I]) DeIterator[T] {
	r := newDeIterator[T](newDeFlattenIterator[I, T](iter.deInner()))
	iter.clear()
	return r
}

// FlatMap creates an iterator that works like map, but flattens nested structure.
//
// The [`Map`] adapter is very useful, but only when the closure
// argument produces values. If it produces an iterator instead, there's
// an extra layer of indirection. `FlatMap()` will remove this extra layer
// on its own.
//
// You can think of `FlatMap(f)` as the semantic equivalent
// of [`Map`]ping, and then [`Flatten`]ing as in `Map(f).Flatten()`.
//
// Another way of thinking about `FlatMap()`: [`Map`]'s closure returns
// one item for each element, and `FlatMap()`'s closure returns an
// iterator for each element.
// NOTE: `iter` will be cleared
func FlatMap[T any, B any](iter Iterator[T], f func(T) Iterator[B]) Iterator[B] {
	r := newIterator[B](newFlatMapIterator[T, B](iter.inner(), f))
	iter.clear()
	return r
}

// DeFlatMap creates a double ended iterator that works like map, but flattens nested structure.
//
// The [`DeMap`] adapter is very useful, but only when the closure
// argument produces values. If it produces an iterator instead, there's
// an extra layer of indirection. `DeFlatMap()` will remove this extra layer
// on its own.
//
// You can think of `DeFlatMap(f)` as the semantic equivalent
// of [`DeMap`]ping, and then [`DeFlatten`]ing as in `DeFlatten(DeMap(f))`.
//
// Another way of thinking about `DeFlatMap()`: [`DeMap`]'s closure returns
// one item for each element, and `DeFlatMap()`'s closure returns an
// iterator for each element.
func DeFlatMap[T any, B any](iter DeIterator[T], f func(T) DeIterator[B]) DeIterator[B] {
	r := newDeIterator[B](newDeFlatMapIterator[T, B](iter.deInner(), f))
	iter.clear()
	return r
}
