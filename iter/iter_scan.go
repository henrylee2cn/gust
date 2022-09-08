package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*scanIterator[any, any, any])(nil)
	_ iRealNext[any]     = (*scanIterator[any, any, any])(nil)
	_ iRealTryFold[any]  = (*scanIterator[any, any, any])(nil)
	_ iRealFold[any]     = (*scanIterator[any, any, any])(nil)
	_ iRealSizeHint      = (*scanIterator[any, any, any])(nil)
)

func newScanIterator[T any, St any, B any](iter innerIterator[T], initialState St, f func(*St, T) gust.Option[B]) innerIterator[B] {
	p := &scanIterator[T, St, B]{iter: iter, f: f, state: initialState}
	p.setFacade(p)
	return p
}

type scanIterator[T any, St any, B any] struct {
	iterBackground[B]
	iter  innerIterator[T]
	state St
	f     func(*St, T) gust.Option[B]
}

func (s *scanIterator[T, St, B]) realNext() gust.Option[B] {
	var a = s.iter.Next()
	if a.IsNone() {
		return gust.None[B]()
	}
	return s.f(&s.state, a.UnwrapUnchecked())
}

func (s *scanIterator[T, St, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = s.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (s *scanIterator[T, St, B]) realTryFold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return s.iter.TryFold(init, func(acc any, x T) gust.AnyCtrlFlow {
		r := s.f(&s.state, x)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.AnyBreak(acc)
	})
}

func (s *scanIterator[T, St, B]) realFold(init any, fold func(any, B) any) any {
	return s.TryFold(init, func(acc any, x B) gust.AnyCtrlFlow {
		return gust.AnyContinue(fold(acc, x))
	}).UnwrapContinue()
}
