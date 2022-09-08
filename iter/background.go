package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

var _ innerIterator[any] = iterBackground[any]{}

type iterBackground[T any] struct {
	facade iRealNext[T]
}

//goland:noinspection GoMixedReceiverTypes
func (iter *iterBackground[T]) setFacade(facade iRealNext[T]) {
	iter.facade = facade
}

func (iter iterBackground[T]) Collect() []T {
	lower, _ := iter.SizeHint()
	return Fold[T, []T](iter, make([]T, 0, lower), func(slice []T, x T) []T {
		return append(slice, x)
	})
}

func (iter iterBackground[T]) Next() gust.Option[T] {
	return iter.facade.realNext()
}

func (iter iterBackground[T]) NextChunk(n uint) gust.EnumResult[[]T, []T] {
	var chunk = make([]T, 0, n)
	for i := uint(0); i < n; i++ {
		item := iter.Next()
		if item.IsSome() {
			chunk = append(chunk, item.Unwrap())
		} else {
			return gust.EnumErr[[]T, []T](chunk)
		}
	}
	return gust.EnumOk[[]T, []T](chunk)
}

func (iter iterBackground[T]) SizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.facade.(iRealSizeHint); ok {
		return cover.realSizeHint()
	}
	return 0, gust.None[uint]()
}

func (iter iterBackground[T]) Count() uint {
	if cover, ok := iter.facade.(iRealCount); ok {
		return cover.realCount()
	}
	return Fold[T, uint](iter,
		uint(0),
		func(count uint, _ T) uint {
			return count + 1
		},
	)
}

func (iter iterBackground[T]) Fold(init any, f func(any, T) any) any {
	if cover, ok := iter.facade.(iRealFold[T]); ok {
		return cover.realFold(init, f)
	}
	return Fold[T, any](iter, init, f)
}

func (iter iterBackground[T]) TryFold(init any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if cover, ok := iter.facade.(iRealTryFold[T]); ok {
		return cover.realTryFold(init, f)
	}
	return TryFold[T, any](iter, init, f)
}

func (iter iterBackground[T]) Last() gust.Option[T] {
	if cover, ok := iter.facade.(iRealLast[T]); ok {
		return cover.realLast()
	}
	return Fold[T, gust.Option[T]](
		iter,
		gust.None[T](),
		func(_ gust.Option[T], x T) gust.Option[T] {
			return gust.Some(x)
		})
}

func (iter iterBackground[T]) AdvanceBy(n uint) gust.Errable[uint] {
	if cover, ok := iter.facade.(iRealAdvanceBy[T]); ok {
		return cover.realAdvanceBy(n)
	}
	for i := uint(0); i < n; i++ {
		if iter.Next().IsNone() {
			return gust.ToErrable[uint](i)
		}
	}
	return gust.NonErrable[uint]()
}

func (iter iterBackground[T]) Nth(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNth[T]); ok {
		return cover.realNth(n)
	}
	var res = iter.AdvanceBy(n)
	if res.IsErr() {
		return gust.None[T]()
	}
	return iter.Next()
}

func (iter iterBackground[T]) ForEach(f func(T)) {
	if cover, ok := iter.facade.(iRealForEach[T]); ok {
		cover.realForEach(f)
		return
	}
	var call = func(f func(T)) func(any, T) any {
		return func(_ any, item T) any {
			f(item)
			return nil
		}
	}
	_ = iter.Fold(nil, call(f))
}

func (iter iterBackground[T]) Reduce(f func(accum T, item T) T) gust.Option[T] {
	if cover, ok := iter.facade.(iRealReduce[T]); ok {
		return cover.realReduce(f)
	}
	var first = iter.Next()
	if first.IsNone() {
		return first
	}
	return gust.Some(Fold[T, T](iter, first.Unwrap(), func(accum T, item T) T {
		return f(accum, item)
	}))
}

func (iter iterBackground[T]) All(predicate func(T) bool) bool {
	if cover, ok := iter.facade.(iRealAll[T]); ok {
		return cover.realAll(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyContinue(nil)
			} else {
				return gust.AnyBreak(nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsContinue()
}

func (iter iterBackground[T]) Any(predicate func(T) bool) bool {
	if cover, ok := iter.facade.(iRealAny[T]); ok {
		return cover.realAny(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(nil)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsBreak()
}

func (iter iterBackground[T]) Find(predicate func(T) bool) gust.Option[T] {
	if cover, ok := iter.facade.(iRealFind[T]); ok {
		return cover.realFind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(x)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsBreak() {
		return gust.Some[T](r.UnwrapBreak().(T))
	}
	return gust.None[T]()
}

func (iter iterBackground[T]) FindMap(f func(T) gust.Option[T]) gust.Option[T] {
	if cover, ok := iter.facade.(iRealFindMap[T]); ok {
		return cover.realFindMap(f)
	}
	return FindMap[T, T](iter, f)
}

func (iter iterBackground[T]) XFindMap(f func(T) gust.Option[any]) gust.Option[any] {
	if cover, ok := iter.facade.(iRealFindMap[T]); ok {
		return cover.realXFindMap(f)
	}
	return FindMap[T, any](iter, f)
}

func (iter iterBackground[T]) TryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
	if cover, ok := iter.facade.(iRealTryFind[T]); ok {
		return cover.realTryFind(predicate)
	}
	var check = func(f func(T) gust.Result[bool]) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			r := f(x)
			if r.IsOk() {
				if r.Unwrap() {
					return gust.AnyBreak(gust.Ok[gust.Option[T]](gust.Some(x)))
				} else {
					return gust.AnyContinue(nil)
				}
			} else {
				return gust.AnyBreak(gust.Err[gust.Option[T]](r.Err()))
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsBreak() {
		return r.UnwrapBreak().(gust.Result[gust.Option[T]])
	}
	return gust.Ok[gust.Option[T]](gust.None[T]())
}

func (iter iterBackground[T]) Position(predicate func(T) bool) gust.Option[int] {
	if cover, ok := iter.facade.(iRealPosition[T]); ok {
		return cover.realPosition(predicate)
	}
	var check = func(f func(T) bool) func(int, T) gust.SigCtrlFlow[int] {
		return func(i int, x T) gust.SigCtrlFlow[int] {
			if f(x) {
				return gust.SigBreak[int](i)
			} else {
				return gust.SigContinue[int](i + 1)
			}
		}
	}
	r := TryFold[T, int](iter, 0, check(predicate))
	if r.IsBreak() {
		return gust.Some[int](r.UnwrapBreak())
	}
	return gust.None[int]()
}

var _ innerDeIterator[any] = deIterBackground[any]{}

type deIterBackground[T any] struct {
	iterBackground[T]
}

//goland:noinspection GoMixedReceiverTypes
func (iter *deIterBackground[T]) setFacade(facade iRealDeIterable[T]) {
	iter.iterBackground.facade = facade
}

func (iter deIterBackground[T]) Remaining() uint {
	if size, ok := iter.facade.(iRealRemaining); ok {
		return size.realRemaining()
	}
	return defaultRemaining[T](iter)
}

func defaultRemaining[T any](iter innerIterator[T]) uint {
	lo, hi := iter.SizeHint()
	if opt.MapOr[uint, bool](hi, false, func(x uint) bool {
		return x == lo
	}) {
		return lo
	}
	return lo
}

func (iter deIterBackground[T]) NextBack() gust.Option[T] {
	return iter.facade.(iRealDeIterable[T]).realNextBack()
}

func (iter deIterBackground[T]) AdvanceBackBy(n uint) gust.Errable[uint] {
	if cover, ok := iter.facade.(iRealAdvanceBackBy[T]); ok {
		return cover.realAdvanceBackBy(n)
	}
	for i := uint(0); i < n; i++ {
		if iter.NextBack().IsNone() {
			return gust.ToErrable[uint](i)
		}
	}
	return gust.NonErrable[uint]()
}

func (iter deIterBackground[T]) NthBack(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNthBack[T]); ok {
		return cover.realNthBack(n)
	}
	if iter.AdvanceBackBy(n).IsErr() {
		return gust.None[T]()
	}
	return iter.NextBack()
}

func (iter deIterBackground[T]) TryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if cover, ok := iter.facade.(iRealTryRfold[T]); ok {
		return cover.realTryRfold(init, fold)
	}
	return TryRfold[T, any](iter, init, fold)
}

func (iter deIterBackground[T]) Rfold(init any, fold func(any, T) any) any {
	if cover, ok := iter.facade.(iRealRfold[T]); ok {
		return cover.realRfold(init, fold)
	}
	return Rfold[T](iter, init, fold)
}

func (iter deIterBackground[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	if cover, ok := iter.facade.(iRealRfind[T]); ok {
		return cover.realRfind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(x)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	r := iter.TryRfold(nil, check(predicate))
	if r.IsBreak() {
		return gust.Some[T](r.UnwrapBreak().(T))
	}
	return gust.None[T]()
}

func (iter deIterBackground[T]) DeFuse() innerDeIterator[T] {
	return newDeFuseIterator[T](iter)
}

func (iter deIterBackground[T]) DePeekable() DePeekableIterator[T] {
	return newDePeekableIterator[T](iter)
}

func (iter deIterBackground[T]) DeSkip(n uint) innerDeIterator[T] {
	return newDeSkipIterator[T](iter, n)
}

func (iter deIterBackground[T]) DeTake(n uint) innerDeIterator[T] {
	return newDeTakeIterator[T](iter, n)
}

func (iter deIterBackground[T]) DeChain(other innerDeIterator[T]) innerDeIterator[T] {
	return newDeChainIterator[T](iter, other)
}

func (iter deIterBackground[T]) DeFilter(f func(T) bool) innerDeIterator[T] {
	return newDeFilterIterator[T](iter, f)
}

func (iter deIterBackground[T]) DeFilterMap(f func(T) gust.Option[T]) innerDeIterator[T] {
	return newDeFilterMapIterator[T, T](iter, f)
}

func (iter deIterBackground[T]) XDeFilterMap(f func(T) gust.Option[any]) innerDeIterator[any] {
	return newDeFilterMapIterator[T, any](iter, f)
}
func (iter deIterBackground[T]) DeInspect(f func(T)) innerDeIterator[T] {
	return newDeInspectIterator[T](iter, f)
}

func (iter deIterBackground[T]) DeMap(f func(T) T) innerDeIterator[T] {
	return newDeMapIterator[T, T](iter, f)
}

func (iter deIterBackground[T]) XDeMap(f func(T) any) innerDeIterator[any] {
	return newDeMapIterator[T, any](iter, f)
}
