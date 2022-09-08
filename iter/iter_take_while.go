package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*takeWhileIterator[any])(nil)
	_ iRealNext[any]     = (*takeWhileIterator[any])(nil)
	_ iRealFold[any]     = (*takeWhileIterator[any])(nil)
	_ iRealTryFold[any]  = (*takeWhileIterator[any])(nil)
	_ iRealSizeHint      = (*takeWhileIterator[any])(nil)
)

func newTakeWhileIterator[T any](iter innerIterator[T], predicate func(T) bool) innerIterator[T] {
	p := &takeWhileIterator[T]{iter: iter, predicate: predicate}
	p.setFacade(p)
	return p
}

type takeWhileIterator[T any] struct {
	iterBackground[T]
	iter      innerIterator[T]
	flag      bool
	predicate func(T) bool
}

func (s *takeWhileIterator[T]) realNext() gust.Option[T] {
	if s.flag {
		return gust.None[T]()
	}
	var x = s.iter.Next()
	if x.IsNone() {
		return x
	}
	next := x.Unwrap()
	if s.predicate(next) {
		return gust.Some(next)
	}
	s.flag = true
	return gust.None[T]()
}

func (s *takeWhileIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if s.flag {
		return 0, gust.Some[uint](0)
	}
	var _, upper = s.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (s *takeWhileIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if s.flag {
		return gust.AnyContinue(init)
	}
	return s.iter.TryFold(init, func(acc any, x T) gust.AnyCtrlFlow {
		if s.predicate(x) {
			return fold(acc, x)
		}
		s.flag = true
		return gust.AnyBreak(acc)
	})
}

func (s *takeWhileIterator[T]) realFold(init any, fold func(any, T) any) any {
	return s.iter.TryFold(init, func(acc any, x T) gust.AnyCtrlFlow {
		return gust.AnyContinue(fold(acc, x))
	}).UnwrapContinue()
}
