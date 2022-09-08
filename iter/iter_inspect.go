package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*inspectIterator[any])(nil)
	_ iRealNext[any]     = (*inspectIterator[any])(nil)
	_ iRealSizeHint      = (*inspectIterator[any])(nil)
	_ iRealTryFold[any]  = (*inspectIterator[any])(nil)
	_ iRealFold[any]     = (*inspectIterator[any])(nil)
)

func newInspectIterator[T any](iter innerIterator[T], f func(T)) innerIterator[T] {
	p := &inspectIterator[T]{iter: iter, f: f}
	p.setFacade(p)
	return p
}

type inspectIterator[T any] struct {
	deIterBackground[T]
	iter innerIterator[T]
	f    func(T)
}

func (s *inspectIterator[T]) doInspect(elt gust.Option[T]) gust.Option[T] {
	if elt.IsSome() {
		s.f(elt.Unwrap())
	}
	return elt
}

func (*inspectIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (s *inspectIterator[T]) realNext() gust.Option[T] {
	return s.doInspect(s.iter.Next())
}

func (s *inspectIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return s.iter.SizeHint()
}

func (s *inspectIterator[T]) realTryFold(init any, g func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryFold[T, any](s.iter, init, func(acc any, elt T) gust.AnyCtrlFlow {
		s.f(elt)
		return g(acc, elt)
	})
}

func (s *inspectIterator[T]) realFold(init any, g func(any, T) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any {
		s.f(elt)
		return g(acc, elt)
	})
}

var (
	_ innerDeIterator[any] = (*deInspectIterator[any])(nil)
	_ iRealRemaining       = (*deInspectIterator[any])(nil)
	_ iRealNextBack[any]   = (*deInspectIterator[any])(nil)
	_ iRealTryRfold[any]   = (*deInspectIterator[any])(nil)
	_ iRealRfold[any]      = (*deInspectIterator[any])(nil)
)

func newDeInspectIterator[T any](iter innerDeIterator[T], f func(T)) innerDeIterator[T] {
	p := &deInspectIterator[T]{}
	p.iter = iter
	p.f = f
	p.setFacade(p)
	return p
}

type deInspectIterator[T any] struct {
	inspectIterator[T]
}

func (d *deInspectIterator[T]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}

func (d *deInspectIterator[T]) realNextBack() gust.Option[T] {
	return d.doInspect(d.iter.(innerDeIterator[T]).NextBack())
}

func (d *deInspectIterator[T]) realTryRfold(init any, g func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryRfold[T, any](d.iter.(innerDeIterator[T]), init, func(acc any, elt T) gust.AnyCtrlFlow {
		d.f(elt)
		return g(acc, elt)
	})
}

func (d *deInspectIterator[T]) realRfold(init any, g func(any, T) any) any {
	return Rfold[T, any](d.iter.(innerDeIterator[T]), init, func(acc any, elt T) any {
		d.f(elt)
		return g(acc, elt)
	})
}
