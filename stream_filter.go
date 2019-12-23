package streamer

type filterStream struct {
	input    Iterator
	filterFn func(interface{}) bool

	currentItem interface{}
}

func newFilterStream(input Iterator, filterFn func(interface{}) bool) (res *filterStream) {
	res = &filterStream{
		input:    input,
		filterFn: filterFn,
	}
	return
}

func (fs *filterStream) Next() bool {
	if fs.currentItem != nil {
		return true
	}
	for fs.input.Next() {
		item := fs.input.Value()
		if fs.filterFn(item) {
			fs.currentItem = item
			return true
		}
	}
	return false
}

func (fs *filterStream) Value() interface{} {
	value := fs.currentItem
	fs.currentItem = nil
	return value
}
