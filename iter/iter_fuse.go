package iter

import (
	"github.com/andeya/gust"
)

var (
	_ innerIterator[any] = (*fuseIterator[any])(nil)
	_ iRealNext[any]     = (*fuseIterator[any])(nil)
	_ iRealNth[any]      = (*fuseIterator[any])(nil)
	_ iRealLast[any]     = (*fuseIterator[any])(nil)
	_ iRealCount         = (*fuseIterator[any])(nil)
	_ iRealSizeHint      = (*fuseIterator[any])(nil)
	_ iRealTryFold[any]  = (*fuseIterator[any])(nil)
	_ iRealFold[any]     = (*fuseIterator[any])(nil)
	_ iRealFind[any]     = (*fuseIterator[any])(nil)
)

func newFuseIterator[T any](iter innerIterator[T]) innerIterator[T] {
	p := &fuseIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type fuseIterator[T any] struct {
	deIterBackground[T]
	iter   innerIterator[T]
	isNone bool
}

func (f *fuseIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (f *fuseIterator[T]) realNext() gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Next().InspectNone(func() {
		f.isNone = true
	})
}

func (f *fuseIterator[T]) realNth(n uint) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Nth(n).InspectNone(func() {
		f.isNone = true
	})
}

func (f *fuseIterator[T]) realLast() gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Last()
}

func (f *fuseIterator[T]) realCount() uint {
	if f.isNone {
		return 0
	}
	f.isNone = true
	return f.iter.Count()
}

func (f *fuseIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if f.isNone {
		return 0, gust.Some[uint](0)
	}
	return f.iter.SizeHint()
}

func (f *fuseIterator[T]) realTryFold(acc any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if f.isNone {
		return gust.AnyBreak("fuse iterator is empty")
	}
	r := f.iter.TryFold(acc, fold)
	f.isNone = true
	return r
}

func (f *fuseIterator[T]) realFold(acc any, fold func(any, T) any) any {
	if f.isNone {
		return acc
	}
	r := f.iter.Fold(acc, fold)
	f.isNone = true
	return r
}

func (f *fuseIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Find(predicate).InspectNone(func() {
		f.isNone = true
	})
}

var (
	_ innerDeIterator[any] = (*deFuseIterator[any])(nil)
	_ iRealNext[any]       = (*deFuseIterator[any])(nil)
	_ iRealNth[any]        = (*deFuseIterator[any])(nil)
	_ iRealLast[any]       = (*deFuseIterator[any])(nil)
	_ iRealCount           = (*deFuseIterator[any])(nil)
	_ iRealSizeHint        = (*deFuseIterator[any])(nil)
	_ iRealTryFold[any]    = (*deFuseIterator[any])(nil)
	_ iRealFold[any]       = (*deFuseIterator[any])(nil)
	_ iRealFind[any]       = (*deFuseIterator[any])(nil)
	_ iRealRemaining       = (*deFuseIterator[any])(nil)
	_ iRealNextBack[any]   = (*deFuseIterator[any])(nil)
	_ iRealNthBack[any]    = (*deFuseIterator[any])(nil)
	_ iRealTryRfold[any]   = (*deFuseIterator[any])(nil)
	_ iRealRfold[any]      = (*deFuseIterator[any])(nil)
	_ iRealRfind[any]      = (*deFuseIterator[any])(nil)
)

func newDeFuseIterator[T any](iter innerDeIterator[T]) innerDeIterator[T] {
	p := &deFuseIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// deFuseIterator double ended fuse iterator
type deFuseIterator[T any] struct {
	fuseIterator[T]
}

func (d *deFuseIterator[T]) realRemaining() uint {
	return d.iter.(innerDeIterator[T]).Remaining()
}

func (d *deFuseIterator[T]) realNextBack() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(innerDeIterator[T])
	return deIter.NextBack().InspectNone(func() {
		d.isNone = true
	})
}

func (d *deFuseIterator[T]) realNthBack(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(innerDeIterator[T])
	return deIter.NthBack(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d *deFuseIterator[T]) realTryRfold(acc any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if d.isNone {
		return gust.AnyBreak("fuse iterator is empty")
	}
	deIter := d.iter.(innerDeIterator[T])
	r := deIter.TryRfold(acc, fold)
	d.isNone = true
	return r
}

func (d *deFuseIterator[T]) realRfold(acc any, fold func(any, T) any) any {
	if d.isNone {
		return acc
	}
	deIter := d.iter.(innerDeIterator[T])
	r := deIter.Rfold(acc, fold)
	d.isNone = true
	return r
}

func (d *deFuseIterator[T]) realRfind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(innerDeIterator[T])
	return deIter.Rfind(predicate).InspectNone(func() {
		d.isNone = true
	})
}
