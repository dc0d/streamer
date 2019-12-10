package streamer

type SliceIterator struct {
	input   []interface{}
	current int
	item    interface{}
}

func NewSliceIterator(input ...interface{}) *SliceIterator {
	res := &SliceIterator{
		input:   input,
		current: -1,
	}

	return res
}

func (it *SliceIterator) Next() bool {
	if it.item != nil {
		return true
	}

	if it.current+1 < len(it.input) {
		it.current++
		it.item = it.input[it.current]
		return true
	}
	return false
}

func (it *SliceIterator) Value() interface{} {
	item := it.item
	it.item = nil
	return item
}

//

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
	iterator := newChunkEveryStrem(st.input, chunkSize)
	return NewStream(iterator)
}

func (st *Stream) Skip(skipCount int) *Stream {
	iterator := newSkipCountStream(st.input, skipCount)
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

//

type takeWhileStream struct {
	input  Iterator
	takeFn func(interface{}) bool

	item  interface{}
	taken bool
}

func newTakeWhileStream(input Iterator, takeFn func(interface{}) bool) (res *takeWhileStream) {
	res = &takeWhileStream{
		input:  input,
		takeFn: takeFn,
	}
	return
}

func (tw *takeWhileStream) Next() bool {
	if tw.taken {
		return false
	}
	if tw.item != nil {
		return true
	}
	next := tw.input.Next()
	if !next {
		return false
	}
	item := tw.input.Value()
	if tw.takeFn(item) {
		tw.item = item
		return true
	}
	tw.taken = true
	return false
}

func (tw *takeWhileStream) Value() interface{} {
	value := tw.item
	tw.item = nil
	return value
}

//

type takeStream struct {
	input     Iterator
	takeCount int

	item interface{}
}

func newTakeStream(input Iterator, takeCount int) (res *takeStream) {
	res = &takeStream{
		input:     input,
		takeCount: takeCount,
	}
	return
}

func (tc *takeStream) Next() bool {
	if tc.item != nil {
		return true
	}
	if tc.takeCount == 0 {
		return false
	}
	next := tc.input.Next()
	if !next {
		return false
	}
	tc.item = tc.input.Value()
	tc.takeCount--
	return true
}

func (tc *takeStream) Value() interface{} {
	value := tc.item
	tc.item = nil
	return value
}

//

type filterStream struct {
	input    Iterator
	filterFn func(interface{}) bool

	item interface{}
}

func newFilterStream(input Iterator, filterFn func(interface{}) bool) (res *filterStream) {
	res = &filterStream{
		input:    input,
		filterFn: filterFn,
	}
	return
}

func (fs *filterStream) Next() bool {
	if fs.item != nil {
		return true
	}
	for fs.input.Next() {
		item := fs.input.Value()
		if fs.filterFn(item) {
			fs.item = item
			return true
		}
	}
	return false
}

func (fs *filterStream) Value() interface{} {
	value := fs.item
	fs.item = nil
	return value
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
}

func newSkipCountStream(input Iterator, skipCount int) (res *skipCountStream) {
	res = &skipCountStream{
		input:     input,
		skipCount: skipCount,
	}
	return
}

func (sc *skipCountStream) Next() bool {
	if sc.skipCount != 0 {
		for i := 0; i < sc.skipCount && sc.input.Next(); i++ {
			sc.input.Value()
		}
		sc.skipCount = 0
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
	if ce.lastChunk != nil {
		return true
	}

	for ce.input.Next() {
		ce.lastChunk = append(ce.lastChunk, ce.input.Value())
		if len(ce.lastChunk) == ce.chunkSize {
			break
		}
	}
	return len(ce.lastChunk) > 0
}

func (ce *chunkEveryStrem) Value() interface{} {
	item := ce.lastChunk
	ce.lastChunk = nil
	return item
}

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
	if cs.lastChunk != nil {
		return true
	}

	var (
		flag           interface{}
		flagCalculated bool
	)

	if cs.lastFlag != nil {
		flag = cs.lastFlag
		flagCalculated = true
		cs.lastFlag = nil
	}

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

func (cs *chunkByStream) Value() interface{} {
	item := cs.lastChunk
	cs.lastChunk = nil
	return item
}

//

type mapperStream struct {
	input Iterator
	mapFn func(x interface{}) interface{}

	item interface{}
}

func newMapperStream(input Iterator, mapFn func(x interface{}) interface{}) (res *mapperStream) {
	res = &mapperStream{
		input: input,
		mapFn: mapFn,
	}
	return
}

func (ms *mapperStream) Next() bool {
	if ms.item != nil {
		return true
	}
	next := ms.input.Next()
	if next {
		ms.item = ms.mapFn(ms.input.Value())
	}
	return next
}

func (ms *mapperStream) Value() interface{} {
	item := ms.item
	ms.item = nil
	return item
}

//

type Iterator interface {
	Next() bool
	Value() interface{}
}
