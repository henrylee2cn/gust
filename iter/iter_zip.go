package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[gust.Pair[any, any]] = (*zipIterator[any, any])(nil)
	_ iRealNext[gust.Pair[any, any]]     = (*zipIterator[any, any])(nil)
	_ iRealSizeHint                      = (*zipIterator[any, any])(nil)
	_ iRealNth[gust.Pair[any, any]]      = (*zipIterator[any, any])(nil)
)

func newZipIterator[A any, B any](a innerIterator[A], b innerIterator[B]) innerIterator[gust.Pair[A, B]] {
	p := &zipIterator[A, B]{a: a, b: b}
	p.setFacade(p)
	return p
}

type (
	zipIterator[A any, B any] struct {
		iterBackground[gust.Pair[A, B]]
		a innerIterator[A]
		b innerIterator[B]
	}
)

func (s *zipIterator[A, B]) superNth(n uint) gust.Option[gust.Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[gust.Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s *zipIterator[A, B]) realNext() gust.Option[gust.Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
}

func (s *zipIterator[A, B]) realSizeHint() (uint, gust.Option[uint]) {
	var aLower, aUpper = s.a.SizeHint()
	var bLower, bUpper = s.b.SizeHint()

	var lower = aLower
	if lower > bLower {
		lower = bLower
	}

	var upper gust.Option[uint]
	if aUpper.IsSome() && bUpper.IsSome() {
		if aUpper.Unwrap() <= bUpper.Unwrap() {
			upper = aUpper
		} else {
			upper = bUpper
		}
	} else if aUpper.IsSome() {
		upper = aUpper
	} else if bUpper.IsSome() {
		upper = bUpper
	}
	return lower, upper
}

func (s *zipIterator[A, B]) realNth(n uint) gust.Option[gust.Pair[A, B]] {
	return s.superNth(n)
}

var (
	_ innerDeIterator[gust.Pair[any, any]] = (*deZipIterator[any, any])(nil)
	_ iRealSizeHint                        = (*deZipIterator[any, any])(nil)
	_ iRealNth[gust.Pair[any, any]]        = (*deZipIterator[any, any])(nil)
	_ iRealDeIterable[gust.Pair[any, any]] = (*deZipIterator[any, any])(nil)
)

func newDeZipIterator[A any, B any](a innerDeIterator[A], b innerDeIterator[B]) innerDeIterator[gust.Pair[A, B]] {
	p := &deZipIterator[A, B]{a: a, b: b}
	p.setFacade(p)
	return p
}

// deZipIterator is a double-ended 2-in-1 iterator with explicit size
type deZipIterator[A any, B any] struct {
	deIterBackground[gust.Pair[A, B]]
	a innerDeIterator[A]
	b innerDeIterator[B]
}

func (s *deZipIterator[A, B]) realRemaining() uint {
	aLen := s.a.Remaining()
	bLen := s.b.Remaining()
	if aLen < bLen {
		return aLen
	}
	return bLen
}

func (s *deZipIterator[A, B]) realNextBack() gust.Option[gust.Pair[A, B]] {
	var aLen = s.a.Remaining()
	var bLen = s.b.Remaining()
	if aLen != bLen {
		// Adjust a, b to equal length
		if aLen > bLen {
			u := aLen - bLen
			for i := uint(0); i < u; i++ {
				s.a.NextBack()
			}
		} else {
			u := bLen - aLen
			for i := uint(0); i < u; i++ {
				s.b.NextBack()
			}
		}
	}
	var x = s.a.NextBack()
	var y = s.b.NextBack()
	if x.IsSome() && y.IsSome() {
		return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
	}
	return gust.None[gust.Pair[A, B]]()
}

func (s *deZipIterator[A, B]) SuperNth(n uint) gust.Option[gust.Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[gust.Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s *deZipIterator[A, B]) realNext() gust.Option[gust.Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
}

func (s *deZipIterator[A, B]) realSizeHint() (uint, gust.Option[uint]) {
	var aLower, aUpper = s.a.SizeHint()
	var bLower, bUpper = s.b.SizeHint()

	var lower = aLower
	if lower > bLower {
		lower = bLower
	}

	var upper gust.Option[uint]
	if aUpper.IsSome() && bUpper.IsSome() {
		if aUpper.Unwrap() <= bUpper.Unwrap() {
			upper = aUpper
		} else {
			upper = bUpper
		}
	} else if aUpper.IsSome() {
		upper = aUpper
	} else if bUpper.IsSome() {
		upper = bUpper
	}
	return lower, upper
}

func (s *deZipIterator[A, B]) realNth(n uint) gust.Option[gust.Pair[A, B]] {
	return s.SuperNth(n)
}
