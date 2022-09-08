package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*skipWhileIterator[any])(nil)
	_ iRealNext[any]     = (*skipWhileIterator[any])(nil)
	_ iRealFold[any]     = (*skipWhileIterator[any])(nil)
	_ iRealTryFold[any]  = (*skipWhileIterator[any])(nil)
	_ iRealSizeHint      = (*skipWhileIterator[any])(nil)
)

func newSkipWhileIterator[T any](iter innerIterator[T], predicate func(T) bool) innerIterator[T] {
	p := &skipWhileIterator[T]{iter: iter, predicate: predicate}
	p.setFacade(p)
	return p
}

type skipWhileIterator[T any] struct {
	iterBackground[T]
	iter      innerIterator[T]
	flag      bool
	predicate func(T) bool
}

func (s *skipWhileIterator[T]) realNext() gust.Option[T] {
	return s.iter.Find(func(v T) bool {
		if s.flag || !s.predicate(v) {
			s.flag = true
			return true
		}
		return false
	})
}

func (s *skipWhileIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = s.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (s *skipWhileIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if !s.flag {
		next := s.realNext()
		if next.IsNone() {
			return gust.AnyContinue(init)
		}
		r := fold(init, next.Unwrap())
		if r.IsBreak() {
			return r
		}
		init = r.UnwrapContinue()
	}
	return s.iter.TryFold(init, fold)
}

func (s *skipWhileIterator[T]) realFold(init any, fold func(any, T) any) any {
	if !s.flag {
		next := s.realNext()
		if next.IsNone() {
			return init
		}
		init = fold(init, next.Unwrap())
	}
	return s.iter.Fold(init, fold)
}
