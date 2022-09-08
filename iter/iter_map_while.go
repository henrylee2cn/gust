package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*mapWhileIterator[any, any])(nil)
	_ iRealNext[any]     = (*mapWhileIterator[any, any])(nil)
	_ iRealTryFold[any]  = (*mapWhileIterator[any, any])(nil)
	_ iRealFold[any]     = (*mapWhileIterator[any, any])(nil)
	_ iRealSizeHint      = (*mapWhileIterator[any, any])(nil)
)

func newMapWhileIterator[T any, B any](iter innerIterator[T], f func(T) gust.Option[B]) innerIterator[B] {
	p := &mapWhileIterator[T, B]{iter: iter, f: f}
	p.setFacade(p)
	return p
}

type mapWhileIterator[T any, B any] struct {
	iterBackground[B]
	iter innerIterator[T]
	f    func(T) gust.Option[B]
}

func (s *mapWhileIterator[T, B]) realNext() gust.Option[B] {
	r := s.iter.Next()
	if r.IsSome() {
		return s.f(r.Unwrap())
	}
	return gust.None[B]()
}

func (s *mapWhileIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = s.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (s *mapWhileIterator[T, B]) realTryFold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	ret := TryFold[T, any](s.iter, init, func(acc any, x T) gust.AnyCtrlFlow {
		r := s.f(x)
		if r.IsSome() {
			y := fold(acc, r.Unwrap())
			if y.IsBreak() {
				return gust.AnyBreak(y)
			}
			return y
		}
		return gust.AnyBreak(gust.AnyContinue(acc))
	})
	if ret.IsBreak() {
		return ret.UnwrapBreak().(gust.AnyCtrlFlow)
	}
	return ret
}

func (s *mapWhileIterator[T, B]) realFold(init any, fold func(any, B) any) any {
	return s.TryFold(init, func(acc any, x B) gust.AnyCtrlFlow {
		return gust.AnyContinue(fold(acc, x))
	}).UnwrapContinue()
}
