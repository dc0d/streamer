package streamer

type skipWhileStream struct {
	input  Iterator
	skipFn func(TItem) bool

	skipped      bool
	leftoverItem TItem
}

func newSkipWhileStream(input Iterator, skipFn func(TItem) bool) (res *skipWhileStream) {
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

func (sw *skipWhileStream) Value() TItem {
	if sw.leftoverItem != nil {
		value := sw.leftoverItem
		sw.leftoverItem = nil
		return value
	}
	return sw.input.Value()
}
