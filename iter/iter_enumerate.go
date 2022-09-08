package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[KV[any]]  = (*enumerateIterator[any])(nil)
	_ iRealNext[KV[any]]      = (*enumerateIterator[any])(nil)
	_ iRealSizeHint           = (*enumerateIterator[any])(nil)
	_ iRealNth[KV[any]]       = (*enumerateIterator[any])(nil)
	_ iRealCount              = (*enumerateIterator[any])(nil)
	_ iRealTryFold[KV[any]]   = (*enumerateIterator[any])(nil)
	_ iRealFold[KV[any]]      = (*enumerateIterator[any])(nil)
	_ iRealAdvanceBy[KV[any]] = (*enumerateIterator[any])(nil)
)

func newEnumerateIterator[T any](iter innerIterator[T]) innerIterator[KV[T]] {
	p := &enumerateIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type (
	// enumerateIterator is an iterator that yields the current count and the element during iteration.
	enumerateIterator[T any] struct {
		deIterBackground[KV[T]]
		iter  innerIterator[T]
		count uint
	}
	// KV is an index-value pair.
	KV[T any] struct {
		Index uint
		Value T
	}
)

func (f *enumerateIterator[T]) realNextBack() gust.Option[KV[T]] {
	panic("unreachable")
}

func (f *enumerateIterator[T]) realNext() gust.Option[KV[T]] {
	var a = f.iter.Next()
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	var i = f.count
	f.count += 1
	return gust.Some(KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *enumerateIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return f.iter.SizeHint()
}

func (f *enumerateIterator[T]) realNth(n uint) gust.Option[KV[T]] {
	var a = f.iter.Nth(n)
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	var i = f.count + n
	f.count = i + 1
	return gust.Some(KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *enumerateIterator[T]) realCount() uint {
	return f.iter.Count()
}

func (f *enumerateIterator[T]) realTryFold(acc any, fold func(any, KV[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(acc, func(acc any, item T) gust.AnyCtrlFlow {
		var r = fold(acc, KV[T]{Index: f.count, Value: item})
		f.count += 1
		return r
	})
}

func (f *enumerateIterator[T]) realFold(acc any, fold func(any, KV[T]) any) any {
	return f.iter.Fold(acc, func(acc any, item T) any {
		var r = fold(acc, KV[T]{Index: f.count, Value: item})
		f.count += 1
		return r
	})
}

func (f *enumerateIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var ret = f.iter.AdvanceBy(n)
	if !ret.IsErr() {
		f.count += n
		return ret
	}
	f.count += ret.UnwrapErr()
	return ret
}

var (
	_ innerDeIterator[KV[any]]    = (*deEnumerateIterator[any])(nil)
	_ iRealNext[KV[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealSizeHint               = (*deEnumerateIterator[any])(nil)
	_ iRealNth[KV[any]]           = (*deEnumerateIterator[any])(nil)
	_ iRealCount                  = (*deEnumerateIterator[any])(nil)
	_ iRealTryFold[KV[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealFold[KV[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBy[KV[any]]     = (*deEnumerateIterator[any])(nil)
	_ iRealNextBack[KV[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealNthBack[KV[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealTryRfold[KV[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealRfold[KV[any]]         = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBackBy[KV[any]] = (*deEnumerateIterator[any])(nil)
	_ iRealRemaining              = (*deEnumerateIterator[any])(nil)
)

func newDeEnumerateIterator[T any](iter innerDeIterator[T]) innerDeIterator[KV[T]] {
	p := &deEnumerateIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// deEnumerateIterator is an iterator that yields the current count and the element during iteration.
type deEnumerateIterator[T any] struct {
	enumerateIterator[T]
}

func (d *deEnumerateIterator[T]) realNextBack() gust.Option[KV[T]] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var a = sizeDeIter.NextBack()
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	return gust.Some(KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realNthBack(n uint) gust.Option[KV[T]] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var a = sizeDeIter.NthBack(n)
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	return gust.Some(KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realTryRfold(init any, fold func(any, KV[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.TryRfold(init, func(acc any, item T) gust.AnyCtrlFlow {
		count -= 1
		return fold(acc, KV[T]{Index: count, Value: item})
	})
}

func (d *deEnumerateIterator[T]) realRfold(acc any, fold func(any, KV[T]) any) any {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.Rfold(acc, func(acc any, item T) any {
		count -= 1
		return fold(acc, KV[T]{Index: count, Value: item})
	})
}

func (d *deEnumerateIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	return d.iter.(innerDeIterator[T]).AdvanceBackBy(n)
}

func (d *deEnumerateIterator[T]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}
