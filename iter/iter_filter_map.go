package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

var (
	_ innerIterator[any] = (*filterMapIterator[any, any])(nil)
	_ iRealNext[any]     = (*filterMapIterator[any, any])(nil)
	_ iRealSizeHint      = (*filterMapIterator[any, any])(nil)
	_ iRealTryFold[any]  = (*filterMapIterator[any, any])(nil)
	_ iRealFold[any]     = (*filterMapIterator[any, any])(nil)
)

func newFilterMapIterator[T any, B any](iter innerIterator[T], filterMap func(T) gust.Option[B]) innerIterator[B] {
	p := &filterMapIterator[T, B]{iter: iter, f: filterMap}
	p.setFacade(p)
	return p
}

type filterMapIterator[T any, B any] struct {
	deIterBackground[B]
	iter innerIterator[T]
	f    func(T) gust.Option[B]
}

func (*filterMapIterator[T, B]) realNextBack() gust.Option[B] {
	panic("unreachable")
}

func (f *filterMapIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the f
}

func (f *filterMapIterator[T, B]) realNext() gust.Option[B] {
	return FindMap[T, B](f.iter, f.f)
}

func (f *filterMapIterator[T, B]) realTryFold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(init, func(acc any, item T) gust.AnyCtrlFlow {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.AnyContinue(acc)
	})
}

func (f *filterMapIterator[T, B]) realFold(init any, fold func(any, B) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return acc
	})
}

var (
	_ innerDeIterator[any] = (*deFilterMapIterator[any, any])(nil)
	_ iRealRemaining       = (*deFilterMapIterator[any, any])(nil)
	_ iRealNextBack[any]   = (*deFilterMapIterator[any, any])(nil)
	_ iRealTryRfold[any]   = (*deFilterMapIterator[any, any])(nil)
	_ iRealRfold[any]      = (*deFilterMapIterator[any, any])(nil)
)

func newDeFilterMapIterator[T any, B any](iter innerDeIterator[T], filterMap func(T) gust.Option[B]) innerDeIterator[B] {
	p := &deFilterMapIterator[T, B]{}
	p.iter = iter
	p.f = filterMap
	p.setFacade(p)
	return p
}

type deFilterMapIterator[T any, B any] struct {
	filterMapIterator[T, B]
}

func (d *deFilterMapIterator[T, B]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}

func (d *deFilterMapIterator[T, B]) realNextBack() gust.Option[B] {
	return opt.XAssert[B](d.iter.(innerDeIterator[T]).TryRfold(nil, func(_ any, x T) gust.AnyCtrlFlow {
		v := d.f(x)
		if v.IsSome() {
			return gust.AnyBreak(v.Unwrap())
		}
		return gust.AnyContinue(nil)
	}).BreakValue())
}

func (d *deFilterMapIterator[T, B]) realTryRfold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return d.iter.(innerDeIterator[T]).TryRfold(init, func(acc any, item T) gust.AnyCtrlFlow {
		r := d.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.AnyContinue(acc)
	})
}

func (d *deFilterMapIterator[T, B]) realRfold(init any, fold func(any, B) any) any {
	return d.iter.(innerDeIterator[T]).Rfold(init, func(acc any, item T) any {
		r := d.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return acc
	})
}
