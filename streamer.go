package streamer

// Iterator interface needs to be implemented by every iterator type.
// Next returns true as long as there are any value that is not consumed yet. This behaviour is consistent with the way channels work in Go.
// Value consumes the next available value. For each call to Next there should be also one call to Value.
type Iterator interface {
	Next() bool
	Value() interface{}
}

//

// Stream consumes values from an Iterator and makes it possible to apply different actions to that stream.
type Stream struct {
	input Iterator
}

func NewStream(input Iterator) (res *Stream) {
	res = &Stream{input: input}
	return
}

func (st *Stream) Next() bool { return st.input.Next() }

func (st *Stream) Value() interface{} { return st.input.Value() }

func (st *Stream) Map(mapFn func(x interface{}) interface{}) *Stream {
	iterator := newMapperStream(st.input, mapFn)
	return NewStream(iterator)
}

func (st *Stream) ChunkBy(chunkFn func(x interface{}) interface{}) *Stream {
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

func (st *Stream) SkipWhile(skipFn func(interface{}) bool) *Stream {
	iterator := newSkipWhileStream(st.input, skipFn)
	return NewStream(iterator)
}

func (st *Stream) Filter(filterFn func(interface{}) bool) *Stream {
	iterator := newFilterStream(st.input, filterFn)
	return NewStream(iterator)
}

func (st *Stream) Take(takeCount int) *Stream {
	iterator := newTakeStream(st.input, takeCount)
	return NewStream(iterator)
}

func (st *Stream) TakeWhile(takeFn func(interface{}) bool) *Stream {
	iterator := newTakeWhileStream(st.input, takeFn)
	return NewStream(iterator)
}
