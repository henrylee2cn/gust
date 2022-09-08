package iter

import "github.com/andeya/gust"

var _ Iterator[any] = (*iterImpl[any])(nil)

type iterImpl[T any] struct {
	innerIterator[T]
}

func newIterator[T any](inner innerIterator[T]) *iterImpl[T] {
	return &iterImpl[T]{innerIterator: inner}
}

func (p *iterImpl[T]) inner() innerIterator[T] {
	return p.innerIterator
}

func (p *iterImpl[T]) reset(inner innerIterator[T]) *iterImpl[T] {
	p.innerIterator = inner
	return p
}

func (p *iterImpl[T]) renewAny(inner innerIterator[any]) *iterImpl[any] {
	p.clear()
	return newIterator(inner)
}

func (p *iterImpl[T]) clear() {
	*p = iterImpl[T]{}
}

func (p *iterImpl[T]) StepBy(step uint) Iterator[T] {
	return p.reset(newStepByIterator(p.innerIterator, step))
}

func (p *iterImpl[T]) Filter(f func(T) bool) Iterator[T] {
	return p.reset(newFilterIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) FilterMap(f func(T) gust.Option[T]) Iterator[T] {
	return p.reset(newFilterMapIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) IntoFilterMap(f func(T) gust.Option[any]) Iterator[any] {
	return p.renewAny(newFilterMapIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) Chain(other Iterator[T]) Iterator[T] {
	return p.reset(newChainIterator(p.inner(), other.inner()))
}

func (p *iterImpl[T]) Map(f func(T) T) Iterator[T] {
	return p.reset(newMapIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) IntoMap(f func(T) any) Iterator[any] {
	return p.renewAny(newMapIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) Inspect(f func(T)) Iterator[T] {
	return p.reset(newInspectIterator(p.innerIterator, f))
}

func (p *iterImpl[T]) Fuse() Iterator[T] {
	return p.reset(newFuseIterator(p.innerIterator))
}

func (p *iterImpl[T]) IntoPeekable() PeekableIterator[T] {
	q := &peekableIterImpl[T]{
		iterImpl: *p,
	}
	p.clear()
	return q
}

func (p *iterImpl[T]) SkipWhile(predicate func(T) bool) Iterator[T] {
	return p.reset(newSkipWhileIterator(p.innerIterator, predicate))
}

func (p *iterImpl[T]) TakeWhile(predicate func(T) bool) Iterator[T] {
	return p.reset(newTakeWhileIterator(p.innerIterator, predicate))
}

func (p *iterImpl[T]) MapWhile(predicate func(T) gust.Option[T]) Iterator[T] {
	return p.reset(newMapWhileIterator(p.innerIterator, predicate))
}

func (p *iterImpl[T]) IntoMapWhile(predicate func(T) gust.Option[any]) Iterator[any] {
	return p.renewAny(newMapWhileIterator(p.innerIterator, predicate))
}

func (p *iterImpl[T]) Skip(n uint) Iterator[T] {
	return p.reset(newSkipIterator(p.innerIterator, n))
}

func (p *iterImpl[T]) Take(n uint) Iterator[T] {
	return p.reset(newTakeIterator(p.innerIterator, n))
}

func (p *iterImpl[T]) IntoScan(initialState any, f func(state *any, item T) gust.Option[any]) Iterator[any] {
	return p.renewAny(newScanIterator(p.innerIterator, initialState, f))
}

// -----------------------------------------------------------------

var _ DeIterator[any] = (*deIterImpl[any])(nil)

func newDeIterator[T any](inner innerDeIterator[T]) *deIterImpl[T] {
	p := &deIterImpl[T]{}
	p.innerIterator = inner
	return p
}

type deIterImpl[T any] struct {
	iterImpl[T]
}

func (p *deIterImpl[T]) deInner() innerDeIterator[T] {
	return p.innerIterator.(innerDeIterator[T])
}

func (p *deIterImpl[T]) deReset(inner innerDeIterator[T]) *deIterImpl[T] {
	p.innerIterator = inner
	return p
}

func (p *deIterImpl[T]) deRenewAny(inner innerDeIterator[any]) *deIterImpl[any] {
	p.clear()
	return newDeIterator(inner)
}

func (p *deIterImpl[T]) Remaining() uint {
	return p.deInner().Remaining()
}

func (p *deIterImpl[T]) NextBack() gust.Option[T] {
	return p.deInner().NextBack()
}

func (p *deIterImpl[T]) AdvanceBackBy(n uint) gust.Errable[uint] {
	return p.deInner().AdvanceBackBy(n)
}

func (p *deIterImpl[T]) NthBack(n uint) gust.Option[T] {
	return p.deInner().NthBack(n)
}

func (p *deIterImpl[T]) TryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return p.deInner().TryRfold(init, fold)
}

func (p *deIterImpl[T]) Rfold(init any, fold func(any, T) any) any {
	return p.deInner().Rfold(init, fold)
}

func (p *deIterImpl[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	return p.deInner().Rfind(predicate)
}

func (p *deIterImpl[T]) DeFuse() DeIterator[T] {
	return p.deReset(newDeFuseIterator(p.deInner()))
}

func (p *deIterImpl[T]) IntoDePeekable() DePeekableIterator[T] {
	q := &dePeekableIterImpl[T]{
		deIterImpl: *p,
	}
	p.clear()
	return q
}

func (p *deIterImpl[T]) DeSkip(n uint) DeIterator[T] {
	return p.deReset(newDeSkipIterator(p.deInner(), n))
}

func (p *deIterImpl[T]) DeTake(n uint) DeIterator[T] {
	return p.deReset(newDeTakeIterator(p.deInner(), n))
}

func (p *deIterImpl[T]) DeChain(other DeIterator[T]) DeIterator[T] {
	return p.deReset(newDeChainIterator(p.deInner(), other.deInner()))
}

func (p *deIterImpl[T]) DeFilter(f func(T) bool) DeIterator[T] {
	return p.deReset(newDeFilterIterator(p.deInner(), f))
}

func (p *deIterImpl[T]) DeFilterMap(f func(T) gust.Option[T]) DeIterator[T] {
	return p.deReset(newDeFilterMapIterator(p.deInner(), f))
}

func (p *deIterImpl[T]) IntoDeFilterMap(f func(T) gust.Option[any]) DeIterator[any] {
	return p.deRenewAny(newDeFilterMapIterator(p.deInner(), f))
}

func (p *deIterImpl[T]) DeInspect(f func(T)) DeIterator[T] {
	return p.deReset(newDeInspectIterator(p.deInner(), f))
}

func (p *deIterImpl[T]) DeMap(f func(T) T) DeIterator[T] {
	return p.deReset(newDeMapIterator(p.deInner(), f))
}

func (p *deIterImpl[T]) IntoDeMap(f func(T) any) DeIterator[any] {
	return p.deRenewAny(newDeMapIterator(p.deInner(), f))
}

// -----------------------------------------------------------------

var _ PeekableIterator[any] = (*peekableIterImpl[any])(nil)

type peekableIterImpl[T any] struct {
	iterImpl[T]
}

func (p *peekableIterImpl[T]) Peek() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

func (p *peekableIterImpl[T]) PeekPtr() gust.Option[*T] {
	// TODO implement me
	panic("implement me")
}

func (p *peekableIterImpl[T]) NextIf(f func(T) bool) gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

// Intersperse creates a new iterator which places a copy of `separator` between adjacent
// items of the original iterator.
func (p *peekableIterImpl[T]) Intersperse(separator T) PeekableIterator[T] {
	panic("not implemented")
	// inner:= newIntersperseIterator[T](&PeekableIterator[T]{Iterator: p.Iterator}, separator)
	// p.Iterator.clear()
	// &PeekableIterator[T]{Iterator: p.Iterator}
}

// IntersperseWith creates a new iterator which places an item generated by `separator`
// between adjacent items of the original iterator.
func (p *peekableIterImpl[T]) IntersperseWith(separator func() T) Iterator[T] {
	panic("not implemented")
	// return p.reset(newIntersperseWithIterator[T](p.innerIterator, separator))
}

var _ DePeekableIterator[any] = (*dePeekableIterImpl[any])(nil)

type dePeekableIterImpl[T any] struct {
	deIterImpl[T]
}

func (d dePeekableIterImpl[T]) NextBack() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) NthBack(n uint) gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) TryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) Rfold(init any, fold func(any, T) any) any {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) Peek() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) PeekPtr() gust.Option[*T] {
	// TODO implement me
	panic("implement me")
}

func (d dePeekableIterImpl[T]) NextIf(f func(T) bool) gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
