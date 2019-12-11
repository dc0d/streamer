package streamer

type mapperStream struct {
	input Iterator
	mapFn func(x TIn) TOut

	currentItem TItem
}

func newMapperStream(input Iterator, mapFn func(x TIn) TOut) (res *mapperStream) {
	res = &mapperStream{
		input: input,
		mapFn: mapFn,
	}
	return
}

func (ms *mapperStream) Next() bool {
	if ms.currentItem != nil {
		return true
	}
	next := ms.input.Next()
	if next {
		ms.currentItem = ms.mapFn(ms.input.Value())
	}
	return next
}

func (ms *mapperStream) Value() TItem {
	item := ms.currentItem
	ms.currentItem = nil
	return item
}
