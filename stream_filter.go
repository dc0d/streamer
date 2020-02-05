package streamer

type filterStream struct {
	input    Iterator
	filterFn func(interface{}) bool
}

func newFilterStream(input Iterator, filterFn func(interface{}) bool) (res *filterStream) {
	res = &filterStream{
		input:    input,
		filterFn: filterFn,
	}
	return
}

func (fs *filterStream) Next() (interface{}, bool) {
	var (
		item interface{}
		ok   = true
	)

	for ok {
		if item, ok = fs.input.Next(); !ok {
			break
		}
		if !fs.filterFn(item) {
			continue
		}
		return item, true
	}

	return nil, false
}
