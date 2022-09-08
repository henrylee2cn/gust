package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any]  = (*skipIterator[any])(nil)
	_ iRealNext[any]      = (*skipIterator[any])(nil)
	_ iRealNth[any]       = (*skipIterator[any])(nil)
	_ iRealCount          = (*skipIterator[any])(nil)
	_ iRealLast[any]      = (*skipIterator[any])(nil)
	_ iRealSizeHint       = (*skipIterator[any])(nil)
	_ iRealTryFold[any]   = (*skipIterator[any])(nil)
	_ iRealFold[any]      = (*skipIterator[any])(nil)
	_ iRealAdvanceBy[any] = (*skipIterator[any])(nil)
)

func newSkipIterator[T any](iter innerIterator[T], n uint) innerIterator[T] {
	p := &skipIterator[T]{iter: iter, n: n}
	p.setFacade(p)
	return p
}

type (
	skipIterator[T any] struct {
		deIterBackground[T]
		iter innerIterator[T]
		n    uint
	}
)

func (f *skipIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (f *skipIterator[T]) realNext() gust.Option[T] {
	if f.n > 0 {
		next := f.iter.Nth(f.n)
		f.n = 0
		return next
	} else {
		return f.iter.Next()
	}
}

func (f *skipIterator[T]) realNth(n uint) gust.Option[T] {
	if f.n <= 0 {
		return f.iter.Nth(n)
	}
	var skip uint = f.n
	f.n = 0
	// Checked add to handle overflow case.
	n2 := checkedAdd(skip, n).UnwrapOrElse(func() uint {
		// In case of overflow, load skip value, before loading `n`.
		// Because the amount of elements to iterate is beyond `usize::MAX`, this
		// is split into two `nth` calls where the `skip` `nth` call is discarded.
		f.iter.Nth(skip - 1)
		return n
	})
	// Load nth element including skip.
	return f.iter.Nth(n2)
}

func (f *skipIterator[T]) realCount() uint {
	if f.n > 0 {
		// Nth(n) skips n+1
		if f.iter.Nth(f.n - 1).IsNone() {
			return 0
		}
	}
	return f.iter.Count()
}

func (f *skipIterator[T]) realLast() gust.Option[T] {
	if f.n > 0 {
		// Nth(n) skips n+1
		x := f.iter.Nth(f.n - 1)
		if x.IsNone() {
			return x
		}
	}
	return f.iter.Last()
}

func (f *skipIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var lower, upper = f.iter.SizeHint()
	lower = saturatingSub(lower, f.n)
	if upper.IsSome() {
		upper = gust.Some(saturatingSub(upper.Unwrap(), f.n))
	}
	return lower, upper
}

func (f *skipIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var n = f.n
	f.n = 0
	if n > 0 {
		// Nth(n) skips n+1
		if f.iter.Nth(n - 1).IsNone() {
			return gust.AnyContinue(init)
		}
	}
	return f.iter.TryFold(init, fold)
}

func (f *skipIterator[T]) realFold(init any, fold func(any, T) any) any {
	if f.n > 0 {
		// Nth(n) skips n+1
		if f.iter.Nth(f.n - 1).IsNone() {
			return init
		}
	}
	return f.iter.Fold(init, fold)
}

func (f *skipIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var rem = n
	var stepOne = saturatingAdd(f.n, rem)
	var advanced = f.iter.AdvanceBy(stepOne)
	if advanced.IsErr() {
		var advancedWithoutSkip = saturatingSub(advanced.UnwrapErr(), f.n)
		f.n = saturatingSub(f.n, advanced.UnwrapErr())
		if n == 0 {
			return gust.NonErrable[uint]()
		} else {
			return gust.ToErrable(advancedWithoutSkip)
		}
	}
	rem -= stepOne - f.n
	f.n = 0
	// step_one calculation may have saturated
	if rem > 0 {
		var advanced = f.iter.AdvanceBy(stepOne)
		if advanced.IsErr() {
			rem -= advanced.UnwrapErr()
			return gust.ToErrable(n - rem)
		}
		return advanced
	}
	return gust.NonErrable[uint]()
}

var (
	_ innerDeIterator[any]    = (*deSkipIterator[any])(nil)
	_ iRealNext[any]          = (*deSkipIterator[any])(nil)
	_ iRealSizeHint           = (*deSkipIterator[any])(nil)
	_ iRealNth[any]           = (*deSkipIterator[any])(nil)
	_ iRealCount              = (*deSkipIterator[any])(nil)
	_ iRealLast[any]          = (*deSkipIterator[any])(nil)
	_ iRealTryFold[any]       = (*deSkipIterator[any])(nil)
	_ iRealFold[any]          = (*deSkipIterator[any])(nil)
	_ iRealAdvanceBy[any]     = (*deSkipIterator[any])(nil)
	_ iRealNextBack[any]      = (*deSkipIterator[any])(nil)
	_ iRealNthBack[any]       = (*deSkipIterator[any])(nil)
	_ iRealTryRfold[any]      = (*deSkipIterator[any])(nil)
	_ iRealRfold[any]         = (*deSkipIterator[any])(nil)
	_ iRealAdvanceBackBy[any] = (*deSkipIterator[any])(nil)
	_ iRealRemaining          = (*deSkipIterator[any])(nil)
)

func newDeSkipIterator[T any](iter innerDeIterator[T], n uint) innerDeIterator[T] {
	p := &deSkipIterator[T]{}
	p.iter = iter
	p.n = n
	p.setFacade(p)
	return p
}

type deSkipIterator[T any] struct {
	skipIterator[T]
}

func (d *deSkipIterator[T]) realNextBack() gust.Option[T] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	if sizeDeIter.Remaining() > 0 {
		return sizeDeIter.NextBack()
	}
	return gust.None[T]()
}

func (d *deSkipIterator[T]) realNthBack(n uint) gust.Option[T] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var remaining = sizeDeIter.Remaining()
	if n < remaining {
		return sizeDeIter.NthBack(n)
	}
	if remaining > 0 {
		// consume the original iterator
		sizeDeIter.NthBack(remaining - 1)
	}
	return gust.None[T]()
}

func (d *deSkipIterator[T]) realTryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var n = sizeDeIter.Remaining()
	if n == 0 {
		return gust.AnyContinue(init)
	}
	var r = sizeDeIter.TryRfold(init, func(acc any, x T) gust.AnyCtrlFlow {
		n -= 1
		var r = fold(acc, x)
		if n == 0 || r.IsBreak() {
			return gust.AnyBreak(r)
		}
		return r
	})
	if r.IsBreak() {
		return r.UnwrapBreak().(gust.AnyCtrlFlow)
	}
	return r
}

func (d *deSkipIterator[T]) realRfold(init any, fold func(any, T) any) any {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	return sizeDeIter.Rfold(init, func(acc any, x T) any {
		return fold(acc, x)
	})
}

func (d *deSkipIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var min = sizeDeIter.Remaining()
	if n < min {
		min = n
	}
	var ret = sizeDeIter.AdvanceBackBy(min)
	if ret.IsErr() {
		panic("iRemaining interface violation")
	}
	if n <= min {
		return ret
	}
	return gust.ToErrable(min)
}

func (d *deSkipIterator[T]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}
