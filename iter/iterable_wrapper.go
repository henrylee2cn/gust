package iter

import (
	"github.com/andeya/gust"
)

func wrapIterable[T any](data gust.Iterable[T]) iterableWrapper[T] {
	return iterableWrapper[T]{data: data}
}

type iterableWrapper[T any] struct {
	data gust.Iterable[T]
}

func (iter *iterableWrapper[T]) realNext() gust.Option[T] {
	return iter.data.Next()
}

func (iter *iterableWrapper[T]) realCount() uint {
	if c, ok := iter.data.(gust.IterableCount); ok {
		return c.Count()
	}
	var a uint
	for iter.data.Next().IsSome() {
		a++
	}
	return a
}

func (iter *iterableWrapper[T]) realSizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.data.(gust.IterableSizeHint); ok {
		return cover.SizeHint()
	}
	return 0, gust.None[uint]()
}

func (iter *iterableWrapper[T]) realNextBack() gust.Option[T] {
	if i, _ := iter.data.(iNextBack[T]); i != nil {
		return i.NextBack()
	}
	panic("not implemented")
}

func (iter *iterableWrapper[T]) realRemaining() uint {
	if i, _ := iter.data.(iRemaining[T]); i != nil {
		return i.Remaining()
	}
	panic("not implemented")
}

type baseIterator[T any] struct {
	iterBackground[T]
	iterableWrapper[T]
}

func fromIterable[T any](data gust.Iterable[T]) innerIterator[T] {
	iter := &baseIterator[T]{iterableWrapper: wrapIterable[T](data)}
	iter.setFacade(iter)
	return iter
}

type baseDeIterator[T any] struct {
	deIterBackground[T]
	iterableWrapper[T]
}

func fromDeIterable[T any](data gust.DeIterable[T]) innerDeIterator[T] {
	iter := &baseDeIterator[T]{iterableWrapper: wrapIterable[T](data)}
	iter.setFacade(iter)
	return iter
}
