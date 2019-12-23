package streamer

type chunkEveryStream struct {
	input     Iterator
	chunkSize int

	lastChunk []interface{}
}

func newChunkEveryStream(input Iterator, chunkSize int) (res *chunkEveryStream) {
	res = &chunkEveryStream{
		input:     input,
		chunkSize: chunkSize,
	}
	return
}

func (ce *chunkEveryStream) Next() bool {
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

func (ce *chunkEveryStream) Value() interface{} {
	item := ce.lastChunk
	ce.lastChunk = nil
	return item
}
