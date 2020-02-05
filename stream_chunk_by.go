package streamer

type chunkByStream struct {
	input   Iterator
	chunkFn func(x interface{}) interface{}

	leftoverChunkItem interface{}
	lastFlag          interface{}
}

func newChunkByStream(input Iterator, chunkFn func(x interface{}) interface{}) (res *chunkByStream) {
	res = &chunkByStream{
		input:   input,
		chunkFn: chunkFn,
	}
	return
}

func (cs *chunkByStream) Next() (interface{}, bool) {
	var (
		flag           interface{}
		flagCalculated bool
		chunk          []interface{}
	)

	if cs.lastFlag != nil {
		flag = cs.lastFlag
		flagCalculated = true
		cs.lastFlag = nil
	}

	if cs.leftoverChunkItem != nil {
		chunk = append(chunk, cs.leftoverChunkItem)
		cs.leftoverChunkItem = nil
	}

	for item, ok := cs.input.Next(); ok; item, ok = cs.input.Next() {
		current := item
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
		chunk = append(chunk, current)
	}

	if len(chunk) == 0 {
		return nil, false
	}

	return chunk, true
}
