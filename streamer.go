package streamer

type (
	TItem      = interface{}
	TIn        = interface{}
	TOut       = interface{}
	TSeparator = interface{}
)

//

type Stream struct {
	input Iterator
}

func NewStream(input Iterator) (res *Stream) {
	res = &Stream{input: input}
	return
}

func (st *Stream) Next() bool { return st.input.Next() }

func (st *Stream) Value() TItem { return st.input.Value() }

func (st *Stream) Map(mapFn func(x TIn) TOut) *Stream {
	iterator := newMapperStream(st.input, mapFn)
	return NewStream(iterator)
}

func (st *Stream) ChunkBy(chunkFn func(x TItem) TSeparator) *Stream {
	iterator := newChunkByStream(st.input, chunkFn)
	return NewStream(iterator)
}

func (st *Stream) ChunkEvery(chunkSize int) *Stream {
	iterator := newChunkEveryStream(st.input, chunkSize)
	return NewStream(iterator)
}

func (st *Stream) Skip(skipCount int) *Stream {
	iterator := newSkipStream(st.input, skipCount)
	return NewStream(iterator)
}

func (st *Stream) SkipWhile(skipFn func(TItem) bool) *Stream {
	iterator := newSkipWhileStream(st.input, skipFn)
	return NewStream(iterator)
}

func (st *Stream) Filter(filterFn func(TItem) bool) *Stream {
	iterator := newFilterStream(st.input, filterFn)
	return NewStream(iterator)
}

func (st *Stream) Take(takeCount int) *Stream {
	iterator := newTakeStream(st.input, takeCount)
	return NewStream(iterator)
}

func (st *Stream) TakeWhile(takeFn func(TItem) bool) *Stream {
	iterator := newTakeWhileStream(st.input, takeFn)
	return NewStream(iterator)
}

//

type Iterator interface {
	Next() bool
	Value() TItem
}
