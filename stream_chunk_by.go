package streamer

type chunkByStream struct {
	input   Iterator
	chunkFn func(x TItem) TSeparator

	leftoverChunkItem TItem
	lastChunk         []TItem
	lastFlag          TItem
}

func newChunkByStream(input Iterator, chunkFn func(x TItem) TSeparator) (res *chunkByStream) {
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
		flag           TItem
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

func (cs *chunkByStream) Value() TItem {
	item := cs.lastChunk
	cs.lastChunk = nil
	return item
}
