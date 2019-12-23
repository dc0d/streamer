package streamer

type SliceIterator struct {
	input       []interface{}
	current     int
	currentItem interface{}
}

func NewSliceIterator(input []interface{}) *SliceIterator {
	res := &SliceIterator{
		input:   input,
		current: -1,
	}

	return res
}

func (it *SliceIterator) Next() bool {
	if it.currentItem != nil {
		return true
	}

	if it.current+1 < len(it.input) {
		it.current++
		it.currentItem = it.input[it.current]
		return true
	}
	return false
}

func (it *SliceIterator) Value() interface{} {
	item := it.currentItem
	it.currentItem = nil
	return item
}
