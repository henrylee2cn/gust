package iter

import (
	"github.com/andeya/gust"
)

func newIntersperseIterator[T any](iter PeekableIterator[T], separator T) innerIterator[T] {
	p := &intersperseIterator[T]{iter: iter, separator: func() T { return separator }}
	p.setFacade(p)
	return p
}

func newIntersperseWithIterator[T any](iter PeekableIterator[T], separator func() T) innerIterator[T] {
	p := &intersperseIterator[T]{iter: iter, separator: separator}
	p.setFacade(p)
	return p
}

var (
	_ innerIterator[any] = (*intersperseIterator[any])(nil)
	_ iRealFold[any]     = (*intersperseIterator[any])(nil)
	_ iRealNext[any]     = (*intersperseIterator[any])(nil)
	_ iRealSizeHint      = (*intersperseIterator[any])(nil)
)

// intersperseIterator is an iterator adapter that places a separator between all elements.
type intersperseIterator[T any] struct {
	iterBackground[T]
	iter      PeekableIterator[T]
	separator func() T
	needsSep  bool
}

func (f *intersperseIterator[T]) realFold(init any, fold func(any, T) any) any {
	var accum = init
	if !f.needsSep {
		if x := f.iter.Next(); x.IsSome() {
			accum = fold(accum, x.Unwrap())
		} else {
			return accum
		}
	}
	return f.iter.Fold(accum, func(accum any, x T) any {
		accum = fold(accum, f.separator())
		accum = fold(accum, x)
		return accum
	})
}

func (f *intersperseIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	lo, hi := f.iter.SizeHint()
	var nextIsElem uint = 0
	if !f.needsSep {
		nextIsElem = 1
	}
	return saturatingAdd(saturatingSub(lo, nextIsElem), lo), hi.AndThen(func(hi uint) gust.Option[uint] {
		return checkedAdd(saturatingSub(hi, nextIsElem), hi)
	})
}

func (f *intersperseIterator[T]) realNext() gust.Option[T] {
	if f.needsSep && f.iter.Peek().IsSome() {
		f.needsSep = false
		return gust.Some(f.separator())
	}
	f.needsSep = true
	return f.iter.Next()
}
