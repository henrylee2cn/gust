package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any]  = (*takeIterator[any])(nil)
	_ iRealNext[any]      = (*takeIterator[any])(nil)
	_ iRealNth[any]       = (*takeIterator[any])(nil)
	_ iRealSizeHint       = (*takeIterator[any])(nil)
	_ iRealTryFold[any]   = (*takeIterator[any])(nil)
	_ iRealFold[any]      = (*takeIterator[any])(nil)
	_ iRealAdvanceBy[any] = (*takeIterator[any])(nil)
)

func newTakeIterator[T any](iter innerIterator[T], n uint) innerIterator[T] {
	p := &takeIterator[T]{iter: iter, n: n}
	p.setFacade(p)
	return p
}

type (
	takeIterator[T any] struct {
		deIterBackground[T]
		iter innerIterator[T]
		n    uint
	}
)

func (f *takeIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (f *takeIterator[T]) realNext() gust.Option[T] {
	if f.n != 0 {
		f.n -= 1
		return f.iter.Next()
	}
	return gust.None[T]()
}

func (f *takeIterator[T]) realNth(n uint) gust.Option[T] {
	if f.n > n {
		f.n -= n + 1
		return f.iter.Nth(n)
	}
	if f.n > 0 {
		f.iter.Nth(f.n - 1)
		f.n = 0
	}
	return gust.None[T]()
}

func (f *takeIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if f.n == 0 {
		return 0, gust.Some[uint](0)
	}
	var lower, upper = f.iter.SizeHint()
	if lower > f.n {
		lower = f.n
	}
	if upper.IsNone() || upper.Unwrap() >= f.n {
		upper = gust.Some(f.n)
	}
	return lower, upper
}

func (f *takeIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if f.n == 0 {
		return gust.AnyContinue(init)
	}
	var r = f.iter.TryFold(init, func(acc any, x T) gust.AnyCtrlFlow {
		f.n -= 1
		var r = fold(acc, x)
		if f.n == 0 || r.IsBreak() {
			return gust.AnyBreak(r)
		}
		return r
	})
	if r.IsBreak() {
		return r.UnwrapBreak().(gust.AnyCtrlFlow)
	}
	return r
}

func (f *takeIterator[T]) realFold(init any, fold func(any, T) any) any {
	return f.TryFold(init, func(acc any, x T) gust.AnyCtrlFlow {
		return gust.AnyContinue(fold(acc, x))
	}).UnwrapContinue()
}

func (f *takeIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var min = f.n
	if min > n {
		min = n
	}
	var r = f.iter.AdvanceBy(min)
	if r.IsErr() {
		f.n -= r.UnwrapErr()
		return r
	}
	f.n -= min
	if min < n {
		return gust.ToErrable(min)
	}
	return gust.NonErrable[uint]()
}

var (
	_ innerDeIterator[any]    = (*deTakeIterator[any])(nil)
	_ iRealNext[any]          = (*deTakeIterator[any])(nil)
	_ iRealSizeHint           = (*deTakeIterator[any])(nil)
	_ iRealNth[any]           = (*deTakeIterator[any])(nil)
	_ iRealTryFold[any]       = (*deTakeIterator[any])(nil)
	_ iRealFold[any]          = (*deTakeIterator[any])(nil)
	_ iRealAdvanceBy[any]     = (*deTakeIterator[any])(nil)
	_ iRealNextBack[any]      = (*deTakeIterator[any])(nil)
	_ iRealNthBack[any]       = (*deTakeIterator[any])(nil)
	_ iRealTryRfold[any]      = (*deTakeIterator[any])(nil)
	_ iRealRfold[any]         = (*deTakeIterator[any])(nil)
	_ iRealAdvanceBackBy[any] = (*deTakeIterator[any])(nil)
	_ iRealRemaining          = (*deTakeIterator[any])(nil)
)

func newDeTakeIterator[T any](iter innerDeIterator[T], n uint) innerDeIterator[T] {
	p := &deTakeIterator[T]{}
	p.iter = iter
	p.n = n
	p.setFacade(p)
	return p
}

type deTakeIterator[T any] struct {
	takeIterator[T]
}

func (d *deTakeIterator[T]) realNextBack() gust.Option[T] {
	if d.n == 0 {
		return gust.None[T]()
	}
	var n = d.n
	d.n -= 1
	var sizeDeIter = d.iter.(innerDeIterator[T])
	return sizeDeIter.NthBack(saturatingSub(sizeDeIter.Remaining(), n))
}

func (d *deTakeIterator[T]) realNthBack(n uint) gust.Option[T] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var remaining = sizeDeIter.Remaining()
	if d.n > n {
		var m = saturatingSub(remaining, d.n) + n
		d.n -= n + 1
		return sizeDeIter.NthBack(m)
	}
	if remaining > 0 {
		sizeDeIter.NthBack(remaining - 1)
	}
	return gust.None[T]()
}

func (d *deTakeIterator[T]) realTryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if d.n == 0 {
		return gust.AnyContinue(init)
	}
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var remaining = sizeDeIter.Remaining()
	if remaining > d.n && sizeDeIter.NthBack(remaining-d.n-1).IsNone() {
		return gust.AnyContinue(init)
	}
	return sizeDeIter.TryRfold(init, fold)
}

func (d *deTakeIterator[T]) realRfold(init any, fold func(any, T) any) any {
	if d.n == 0 {
		return init
	}
	var sizeDeIter = d.iter.(innerDeIterator[T])
	var remaining = sizeDeIter.Remaining()
	if remaining > d.n && sizeDeIter.NthBack(remaining-d.n-1).IsNone() {
		return init
	}
	return sizeDeIter.Rfold(init, fold)
}

func (d *deTakeIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	var sizeDeIter = d.iter.(innerDeIterator[T])
	// The amount by which the inner iterator needs to be shortened for it to be
	// at most as long as the take() amount.
	var trimInner = saturatingSub(sizeDeIter.Remaining(), d.n)
	// The amount we need to advance inner to fulfill the caller's request.
	// take(), advance_by() and len() all can be at most usize, so we don't have to worry
	// about having to advance more than usize::MAX here.
	var advanceBy = saturatingAdd(trimInner, n)
	var r = sizeDeIter.AdvanceBackBy(advanceBy)
	var advanced uint
	if r.IsErr() {
		advanced = r.UnwrapErr() - trimInner
	} else {
		advanced = advanceBy - trimInner
	}
	d.n -= advanced
	if advanced < n {
		return gust.ToErrable(advanced)
	}
	return gust.NonErrable[uint]()
}

func (d *deTakeIterator[T]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}
