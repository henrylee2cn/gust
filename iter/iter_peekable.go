package iter

import (
	"github.com/andeya/gust"
)

var (
	_ iRealNext[any]     = (*implPeekable[any])(nil)
	_ iRealCount         = (*implPeekable[any])(nil)
	_ iRealNth[any]      = (*implPeekable[any])(nil)
	_ iRealLast[any]     = (*implPeekable[any])(nil)
	_ iRealSizeHint      = (*implPeekable[any])(nil)
	_ iRealTryFold[any]  = (*implPeekable[any])(nil)
	_ iRealFold[any]     = (*implPeekable[any])(nil)
	_ iRealNextBack[any] = (*implPeekable[any])(nil)
	_ iRealTryRfold[any] = (*implPeekable[any])(nil)
	_ iRealRfold[any]    = (*implPeekable[any])(nil)
	_ iRealRemaining     = (*implPeekable[any])(nil)
)

func newPeekableIterator[T any](iter innerIterator[T]) PeekableIterator[T] {
	p := &peekableIterator[T]{
		implPeekable: implPeekable[T]{iter: iter},
	}
	p.setFacade(p)
	return p
}

type peekableIterator[T any] struct {
	iterBackground[T]
	implPeekable[T]
}

func (s *peekableIterator[T]) Intersperse(separator T) innerIterator[T] {
	return newIntersperseIterator[T](s, separator)
}

func (s *peekableIterator[T]) IntersperseWith(separator func() T) innerIterator[T] {
	return newIntersperseWithIterator[T](s, separator)
}

func newDePeekableIterator[T any](iter innerDeIterator[T]) DePeekableIterator[T] {
	p := &dePeekableIterator[T]{
		implPeekable: implPeekable[T]{iter: iter},
	}
	p.setFacade(p)
	return p
}

type dePeekableIterator[T any] struct {
	deIterBackground[T]
	implPeekable[T]
}

func (s *dePeekableIterator[T]) Intersperse(separator T) innerIterator[T] {
	return newIntersperseIterator[T](s, separator)
}

func (s *dePeekableIterator[T]) IntersperseWith(separator func() T) innerIterator[T] {
	return newIntersperseWithIterator[T](s, separator)
}

type implPeekable[T any] struct {
	iter innerIterator[T]
	// Remember a peeked value, even if it was None.
	peeked gust.Option[gust.Option[T]]
}

func (s *implPeekable[T]) realRemaining() uint {
	if x, ok := s.iter.(iRemaining[T]); ok {
		return x.Remaining()
	}
	return defaultRemaining(s.iter)
}

func (s *implPeekable[T]) NextIf(f func(T) bool) gust.Option[T] {
	next := s.realNext()
	if next.IsSome() {
		matched := next.Unwrap()
		if f(matched) {
			return gust.Some(matched)
		}
	}
	s.peeked = gust.Some(next)
	return gust.None[T]()
}

func (s *implPeekable[T]) Peek() gust.Option[T] {
	return *s.peeked.GetOrInsertWith(s.iter.Next)
}

func (s *implPeekable[T]) PeekPtr() gust.Option[*T] {
	x := s.peeked.GetOrInsertWith(s.iter.Next)
	if x.IsNone() {
		return gust.None[*T]()
	}
	return gust.Some[*T](x.GetOrInsertWith(nil))
}

func (s *implPeekable[T]) realNext() gust.Option[T] {
	taken := s.peeked.Take()
	if taken.IsSome() {
		return taken.Unwrap()
	}
	return s.iter.Next()
}

func (s *implPeekable[T]) realNextBack() gust.Option[T] {
	if s.peeked.IsSome() {
		peeked := s.peeked.Unwrap()
		if peeked.IsSome() {
			return s.iter.(iNextBack[T]).NextBack().OrElse(func() gust.Option[T] { return peeked.Take() })
		}
		return gust.None[T]()
	}
	return s.iter.(iNextBack[T]).NextBack()
}

func (s *implPeekable[T]) realCount() uint {
	taken := s.peeked.Take()
	if taken.IsSome() {
		peeked := taken.Unwrap()
		if peeked.IsNone() {
			return 0
		}
		return 1 + s.iter.Count()
	}
	return s.iter.Count()
}

func (s *implPeekable[T]) realNth(n uint) gust.Option[T] {
	taken := s.peeked.Take()
	if taken.IsNone() {
		return gust.None[T]()
	}
	peeked := taken.Unwrap()
	if peeked.IsNone() {
		return gust.None[T]()
	}
	if n == 0 {
		return peeked
	}
	return s.iter.Nth(n - 1)
}

func (s *implPeekable[T]) realLast() gust.Option[T] {
	var peekOpt gust.Option[T]
	taken := s.peeked.Take()
	if taken.IsSome() {
		peeked := taken.Unwrap()
		if peeked.IsNone() {
			return gust.None[T]()
		}
		peekOpt = peeked
	}
	return s.iter.Last().Or(peekOpt)
}

func (s *implPeekable[T]) realSizeHint() (uint, gust.Option[uint]) {
	var peekLen uint
	if s.peeked.IsSome() {
		peeked := s.peeked.Unwrap()
		if peeked.IsNone() {
			return 0, gust.Some[uint](0)
		}
		peekLen = 1
	}
	lo, hi := s.iter.SizeHint()
	lo = saturatingAdd(lo, peekLen)
	if hi.IsSome() {
		hi = checkedAdd(hi.Unwrap(), peekLen)
	}
	return lo, hi
}

func (s *implPeekable[T]) realTryFold(init any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var acc = init
	taken := s.peeked.Take()
	if taken.IsSome() {
		peeked := taken.Unwrap()
		if peeked.IsNone() {
			return gust.AnyContinue(init)
		}
		x := f(init, peeked.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x.UnwrapContinue()
	}
	return TryFold[T, any](s.iter, acc, f)
}

func (s *implPeekable[T]) realTryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var taken = s.peeked.Take()
	if taken.IsNone() {
		return TryRfold[T, any](s.iter.(innerDeIterator[T]), init, fold)
	}
	peeked := taken.Unwrap()
	if peeked.IsNone() {
		return gust.AnyContinue(init)
	}
	v := peeked.Unwrap()
	r := TryRfold[T, any](s.iter.(innerDeIterator[T]), init, fold)
	if r.IsContinue() {
		return fold(r.UnwrapContinue(), v)
	}
	s.peeked = taken
	return r
}

func (s *implPeekable[T]) realFold(init any, f func(any, T) any) any {
	var acc = init
	taken := s.peeked.Take()
	if taken.IsSome() {
		peeked := taken.Unwrap()
		if peeked.IsNone() {
			return init
		}
		acc = f(init, peeked.Unwrap())
	}
	return Fold[T, any](s.iter, acc, f)
}

func (s *implPeekable[T]) realRfold(init any, fold func(any, T) any) any {
	var taken = s.peeked.Take()
	if taken.IsNone() {
		return Rfold[T, any](s.iter.(innerDeIterator[T]), init, fold)
	}
	peeked := taken.Unwrap()
	if peeked.IsNone() {
		return init
	}
	v := peeked.Unwrap()
	acc := Rfold[T, any](s.iter.(innerDeIterator[T]), init, fold)
	return fold(acc, v)
}
