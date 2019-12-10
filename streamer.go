package streamer

type SliceIterator struct {
	input   []interface{}
	current int
}

func NewSliceIterator(input ...interface{}) *SliceIterator {
	res := &SliceIterator{
		input:   input,
		current: -1,
	}

	return res
}

func (it *SliceIterator) Next() bool {
	if it.current+1 < len(it.input) {
		it.current++
		return true
	}
	return false
}

func (it *SliceIterator) Value() interface{} { return it.input[it.current] }

//

type Stream struct {
	input    Iterator
	iterator Iterator
}

func NewStream(input Iterator) (res *Stream) {
	res = &Stream{input: input}
	return
}

func (st *Stream) Next() bool { return st.iterator.Next() }

func (st *Stream) Value() interface{} { return st.iterator.Value() }

func (st *Stream) Map(mapFn func(x interface{}) interface{}) {
	st.iterator = newMapperStream(st.input, mapFn)
}

func (st *Stream) ChunkBy(chunkFn func(x interface{}) interface{}) {
	st.iterator = newChunkByStream(st.input, chunkFn)
}

func (st *Stream) ChunkEvery(chunkSize int) {
	st.iterator = newChunkEveryStrem(st.input, chunkSize)
}

func (st *Stream) Skip(skipCount int) {
	st.iterator = newSkipCountStream(st.input, skipCount)
}

func (st *Stream) SkipWhile(skipFn func(interface{}) bool) {
	st.iterator = newSkipWhileStream(st.input, skipFn)
}

//

type skipWhileStream struct {
	input  Iterator
	skipFn func(interface{}) bool

	skipped      bool
	leftoverItem interface{}
}

func newSkipWhileStream(input Iterator, skipFn func(interface{}) bool) (res *skipWhileStream) {
	res = &skipWhileStream{
		input:  input,
		skipFn: skipFn,
	}
	return
}

func (sw *skipWhileStream) Next() bool {
	if !sw.skipped {
		for sw.input.Next() {
			value := sw.input.Value()
			if !sw.skipFn(value) {
				sw.leftoverItem = value
				break
			}
		}
		sw.skipped = true
	}
	if sw.leftoverItem != nil {
		return true
	}
	return sw.input.Next()
}

func (sw *skipWhileStream) Value() interface{} {
	if sw.leftoverItem != nil {
		value := sw.leftoverItem
		sw.leftoverItem = nil
		return value
	}
	return sw.input.Value()
}

//

type skipCountStream struct {
	input     Iterator
	skipCount int

	skipped bool
}

func newSkipCountStream(input Iterator, skipCount int) (res *skipCountStream) {
	res = &skipCountStream{
		input:     input,
		skipCount: skipCount,
	}
	return
}

func (sc *skipCountStream) Next() bool {
	if !sc.skipped {
		for i := 0; i < sc.skipCount && sc.input.Next(); i++ {
			sc.input.Value()
		}
		sc.skipped = true
	}
	return sc.input.Next()
}

func (sc *skipCountStream) Value() interface{} { return sc.input.Value() }

//

type chunkEveryStrem struct {
	input     Iterator
	chunkSize int

	lastChunk []interface{}
}

func newChunkEveryStrem(input Iterator, chunkSize int) (res *chunkEveryStrem) {
	res = &chunkEveryStrem{
		input:     input,
		chunkSize: chunkSize,
	}
	return
}

func (ce *chunkEveryStrem) Next() bool {
	ce.lastChunk = nil
	for ce.input.Next() {
		ce.lastChunk = append(ce.lastChunk, ce.input.Value())
		if len(ce.lastChunk) == ce.chunkSize {
			break
		}
	}
	return len(ce.lastChunk) > 0
}

func (ce *chunkEveryStrem) Value() interface{} { return ce.lastChunk }

//

type chunkByStream struct {
	input   Iterator
	chunkFn func(x interface{}) interface{}

	leftoverChunkItem interface{}
	lastChunk         []interface{}
	lastFlag          interface{}
}

func newChunkByStream(input Iterator, chunkFn func(x interface{}) interface{}) (res *chunkByStream) {
	res = &chunkByStream{
		input:   input,
		chunkFn: chunkFn,
	}
	return
}

func (cs *chunkByStream) Next() bool {
	var (
		flag           interface{}
		flagCalculated bool
	)

	if cs.lastFlag != nil {
		flag = cs.lastFlag
		flagCalculated = true
		cs.lastFlag = nil
	}

	cs.lastChunk = nil
	if cs.leftoverChunkItem != nil {
		cs.lastChunk = append(cs.lastChunk, cs.leftoverChunkItem)
		cs.leftoverChunkItem = nil
	}

	for cs.input.Next() {
		current := cs.input.Value()
		cond := cs.chunkFn(current)

		if !flagCalculated {
			flag = cond
			flagCalculated = true
		}

		if flag != cond {
			cs.leftoverChunkItem = current
			cs.lastFlag = cond
			break
		}

		flag = cond
		cs.lastChunk = append(cs.lastChunk, current)
	}

	return len(cs.lastChunk) > 0
}

func (cs *chunkByStream) Value() interface{} { return cs.lastChunk }

//

type mapperStream struct {
	input Iterator
	mapFn func(x interface{}) interface{}
}

func newMapperStream(input Iterator, mapFn func(x interface{}) interface{}) (res *mapperStream) {
	res = &mapperStream{
		input: input,
		mapFn: mapFn,
	}
	return
}

func (ms *mapperStream) Next() bool         { return ms.input.Next() }
func (ms *mapperStream) Value() interface{} { return ms.mapFn(ms.input.Value()) }

//

type Iterator interface {
	Next() bool
	Value() interface{}
}
