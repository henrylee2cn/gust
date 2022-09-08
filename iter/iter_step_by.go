package iter

import (
	"math"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

var (
	_ innerIterator[any] = (*stepByIterator[any])(nil)
	_ iRealNext[any]     = (*stepByIterator[any])(nil)
	_ iRealSizeHint      = (*stepByIterator[any])(nil)
	_ iRealNth[any]      = (*stepByIterator[any])(nil)
	_ iRealTryFold[any]  = (*stepByIterator[any])(nil)
	_ iRealFold[any]     = (*stepByIterator[any])(nil)
)

func newStepByIterator[T any](iter innerIterator[T], step uint) innerIterator[T] {
	if step == 0 {
		panic("step must be non-zero")
	}
	p := &stepByIterator[T]{iter: iter, step: step - 1, firstTake: true}
	p.setFacade(p)
	return p
}

type stepByIterator[T any] struct {
	iterBackground[T]
	iter      innerIterator[T]
	step      uint
	firstTake bool
}

func (s *stepByIterator[T]) realFold(acc any, f func(any, T) any) any {
	if s.firstTake {
		s.firstTake = false
		r := s.iter.Next()
		if r.IsNone() {
			return acc
		}
		acc = f(acc, r.Unwrap())
	}
	r := s.iter.Nth(s.step)
	if r.IsSome() {
		return f(acc, r.Unwrap())
	}
	return acc
}

func (s *stepByIterator[T]) realTryFold(acc any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if s.firstTake {
		s.firstTake = false
		r := s.iter.Next()
		if r.IsNone() {
			return gust.AnyContinue(acc)
		}
		v := f(acc, r.Unwrap())
		if v.IsBreak() {
			return v
		}
		acc = v.UnwrapContinue()
	}
	r := s.iter.Nth(s.step)
	if r.IsSome() {
		return f(acc, r.Unwrap())
	}
	return gust.AnyContinue(acc)
}

func (s *stepByIterator[T]) realNth(n uint) gust.Option[T] {
	if s.firstTake {
		s.firstTake = false
		var first = s.iter.Next()
		if n == 0 {
			return first
		}
		n -= 1
	}
	// n and s.step are indices, we need to add 1 to get the amount of elements
	// When calling `.Nth`, we need to subtract 1 again to convert back to an index
	// step + 1 can't overflow because `.step_by` sets `s.step` to `step - 1`
	var step = s.step + 1
	// n + 1 could overflow
	// thus, if n is math.MaxUint, instead of adding one, we call .Nth(step)
	if n == math.MaxUint {
		s.iter.Nth(step - 1)
	} else {
		n += 1
	}

	// overflow handling
	for {
		var mul = uintCheckedMul(n, step)
		if mul.IsSome() {
			return s.iter.Nth(mul.Unwrap() - 1)
		}
		var divN = math.MaxUint / n
		var divStep = math.MaxUint / step
		var nthN = divN * n
		var nthStep = divStep * step
		var nth uint
		if nthN > nthStep {
			step -= divN
			nth = nthN
		} else {
			n -= divStep
			nth = nthStep
		}
		s.iter.Nth(nth - 1)
	}
}

func (s *stepByIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var firstSize = func(step uint) func(uint) uint {
		return func(n uint) uint {
			if n == 0 {
				return 0
			}
			return 1 + (n-1)/(step+1)
		}
	}

	var otherSize = func(step uint) func(uint) uint {
		return func(n uint) uint { return n / (step + 1) }
	}

	var low, high = s.iter.SizeHint()

	if s.firstTake {
		var f = firstSize(uint(s.step))
		return f(low), opt.Map[uint](high, f)
	}
	var f = otherSize(uint(s.step))
	return f(low), opt.Map[uint](high, f)
}

func (s *stepByIterator[T]) realNext() gust.Option[T] {
	if s.firstTake {
		s.firstTake = false
		return s.iter.Next()
	}
	return s.iter.Nth(s.step)
}
