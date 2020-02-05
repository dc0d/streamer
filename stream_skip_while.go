package streamer

type skipWhileStream struct {
	input  Iterator
	skipFn func(interface{}) bool

	skipped bool
}

func newSkipWhileStream(input Iterator, skipFn func(interface{}) bool) (res *skipWhileStream) {
	res = &skipWhileStream{
		input:  input,
		skipFn: skipFn,
	}
	return
}

func (sw *skipWhileStream) Next() (interface{}, bool) {
	if !sw.skipped {
		var (
			item interface{}
			ok   bool
		)

		item, ok = sw.input.Next()
		for ok && sw.skipFn(item) {
			item, ok = sw.input.Next()
		}

		sw.skipped = true

		if item != nil {
			return item, true
		}
	}

	return sw.input.Next()
}
