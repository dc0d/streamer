package streamer

type chunkEveryStream struct {
	input     Iterator
	chunkSize int
}

func newChunkEveryStream(input Iterator, chunkSize int) (res *chunkEveryStream) {
	res = &chunkEveryStream{
		input:     input,
		chunkSize: chunkSize,
	}
	return
}

func (ce *chunkEveryStream) Next() (interface{}, bool) {
	var (
		chunk []interface{}
	)

	for item, ok := ce.input.Next(); ok; item, ok = ce.input.Next() {
		chunk = append(chunk, item)
		if len(chunk) == ce.chunkSize {
			break
		}
	}

	if len(chunk) == 0 {
		return nil, false
	}

	return chunk, true
}
