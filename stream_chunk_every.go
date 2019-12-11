package streamer

type chunkEveryStream struct {
	input     Iterator
	chunkSize int

	lastChunk []TItem
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

func (ce *chunkEveryStream) Value() TItem {
	item := ce.lastChunk
	ce.lastChunk = nil
	return item
}
